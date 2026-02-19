package backup

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/mirkobrombin/dabadee/pkg/dabadee"
	"github.com/mirkobrombin/dabadee/pkg/hash"
	"github.com/mirkobrombin/dabadee/pkg/processor"
	"github.com/mirkobrombin/dabadee/pkg/storage"
	"github.com/vanilla-os/sdk/pkg/v1/fs"
)

type Repository struct {
	Root         string
	SnapshotsDir string
	ObjectsDir   string
}

type CreateSnapshotOptions struct {
	ID                string
	Deduplicate       bool
	DedupWorkers      int
	DedupWithMetadata bool
	CopyOptions       fs.CopyTreeOptions
}

type Snapshot struct {
	Manifest SnapshotManifest
	Path     string
	TreePath string
}

// OpenRepository opens (or initializes) a snapshot repository under root.
//
// Example:
//
//	repo, err := backup.OpenRepository("/var/lib/myrepo")
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//	_ = repo
func OpenRepository(root string) (*Repository, error) {
	snapshots := filepath.Join(root, "snapshots")
	objects := filepath.Join(root, "objects")

	if err := os.MkdirAll(snapshots, 0o755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Join(snapshots, ".tmp"), 0o755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(objects, 0o755); err != nil {
		return nil, err
	}

	return &Repository{Root: root, SnapshotsDir: snapshots, ObjectsDir: objects}, nil
}

// CreateSnapshot copies sourcePath into a new snapshot directory.
//
// Example:
//
//	snap, err := repo.CreateSnapshot("/home/user", backup.CreateSnapshotOptions{Deduplicate: true})
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//	fmt.Printf("Snapshot: %s\n", snap.Manifest.ID)
func (r *Repository) CreateSnapshot(sourcePath string, opts CreateSnapshotOptions) (*Snapshot, error) {
	l, err := r.lockExclusive()
	if err != nil {
		return nil, err
	}
	defer func() { _ = l.Close() }()

	id := opts.ID
	if id == "" {
		id = newSnapshotID()
	}

	finalPath := filepath.Join(r.SnapshotsDir, id)
	tmpBase := filepath.Join(r.SnapshotsDir, ".tmp")
	tmpPath := filepath.Join(tmpBase, id)
	treePath := filepath.Join(tmpPath, "tree")

	if err := os.MkdirAll(tmpBase, 0o755); err != nil {
		return nil, err
	}
	_ = os.RemoveAll(tmpPath)
	if err := os.MkdirAll(treePath, 0o755); err != nil {
		return nil, err
	}

	copyOpts := opts.CopyOptions
	if copyOpts.Workers == 0 {
		copyOpts = DefaultCopyOptions()
	}
	if err := fs.CopyTree(sourcePath, treePath, copyOpts); err != nil {
		_ = os.RemoveAll(tmpPath)
		return nil, err
	}

	if opts.Deduplicate {
		workers := opts.DedupWorkers
		if workers <= 0 {
			workers = 2
		}

		s, err := storage.NewStorage(storage.StorageOptions{Root: r.ObjectsDir, WithMetadata: opts.DedupWithMetadata})
		if err != nil {
			_ = os.RemoveAll(tmpPath)
			return nil, err
		}
		h := hash.NewSHA256Generator()
		p := processor.NewDedupProcessor(treePath, "", s, h, workers)
		d := dabadee.NewDaBaDee(p, false)
		if err := d.Run(); err != nil {
			_ = os.RemoveAll(tmpPath)
			return nil, err
		}
	}

	m := SnapshotManifest{
		ID:          id,
		CreatedAt:   time.Now().UTC(),
		SourcePath:  sourcePath,
		Deduplicate: opts.Deduplicate,
	}
	if err := writeManifest(filepath.Join(tmpPath, "manifest.json"), m); err != nil {
		_ = os.RemoveAll(tmpPath)
		return nil, err
	}

	if _, err := os.Stat(finalPath); err == nil {
		_ = os.RemoveAll(tmpPath)
		return nil, fmt.Errorf("snapshot already exists: %s", id)
	}
	if err := os.Rename(tmpPath, finalPath); err != nil {
		_ = os.RemoveAll(tmpPath)
		return nil, err
	}

	return &Snapshot{Manifest: m, Path: finalPath, TreePath: filepath.Join(finalPath, "tree")}, nil
}

// RestoreSnapshot restores snapshotID into destination.
//
// Example:
//
//	err := repo.RestoreSnapshot("20260218T120000Z-acde1234", "/restore", backup.DefaultCopyOptions())
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
func (r *Repository) RestoreSnapshot(snapshotID, destination string, copyOpts fs.CopyTreeOptions) error {
	l, err := r.lockShared()
	if err != nil {
		return err
	}
	defer func() { _ = l.Close() }()

	tree := filepath.Join(r.SnapshotsDir, snapshotID, "tree")
	if _, err := os.Stat(tree); err != nil {
		return err
	}
	if copyOpts.Workers == 0 {
		copyOpts = DefaultCopyOptions()
	}
	return fs.CopyTree(tree, destination, copyOpts)
}

// ListSnapshots lists known snapshot manifests.
//
// Example:
//
//	snaps, err := repo.ListSnapshots()
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//	for _, s := range snaps {
//		fmt.Printf("%s %s\n", s.ID, s.CreatedAt)
//	}
func (r *Repository) ListSnapshots() ([]SnapshotManifest, error) {
	l, err := r.lockShared()
	if err != nil {
		return nil, err
	}
	defer func() { _ = l.Close() }()

	entries, err := os.ReadDir(r.SnapshotsDir)
	if err != nil {
		return nil, err
	}

	snaps := make([]SnapshotManifest, 0)
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		m, err := readManifest(filepath.Join(r.SnapshotsDir, e.Name(), "manifest.json"))
		if err != nil {
			continue
		}
		snaps = append(snaps, m)
	}

	sort.Slice(snaps, func(i, j int) bool { return snaps[i].CreatedAt.After(snaps[j].CreatedAt) })
	return snaps, nil
}

func (r *Repository) listSnapshotsUnlocked() ([]SnapshotManifest, error) {
	entries, err := os.ReadDir(r.SnapshotsDir)
	if err != nil {
		return nil, err
	}

	snaps := make([]SnapshotManifest, 0)
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		m, err := readManifest(filepath.Join(r.SnapshotsDir, e.Name(), "manifest.json"))
		if err != nil {
			continue
		}
		snaps = append(snaps, m)
	}

	sort.Slice(snaps, func(i, j int) bool { return snaps[i].CreatedAt.After(snaps[j].CreatedAt) })
	return snaps, nil
}

// PruneKeepLast removes snapshots keeping only the most recent keepLast.
//
// Example:
//
//	removed, err := repo.PruneKeepLast(7)
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//	fmt.Printf("Removed: %v\n", removed)
func (r *Repository) PruneKeepLast(keepLast int) ([]string, error) {
	l, err := r.lockExclusive()
	if err != nil {
		return nil, err
	}
	defer func() { _ = l.Close() }()

	if keepLast < 0 {
		return nil, fmt.Errorf("keepLast must be >= 0")
	}

	snaps, err := r.listSnapshotsUnlocked()
	if err != nil {
		return nil, err
	}

	if len(snaps) <= keepLast {
		return nil, nil
	}

	removed := make([]string, 0)
	for _, s := range snaps[keepLast:] {
		_ = os.RemoveAll(filepath.Join(r.SnapshotsDir, s.ID))
		removed = append(removed, s.ID)
	}

	// Best-effort object cleanup if DaBaDee storage exists.
	if _, err := os.Stat(filepath.Join(r.ObjectsDir, ".dabadee")); err == nil {
		if st, err := storage.NewStorage(storage.StorageOptions{Root: r.ObjectsDir}); err == nil {
			_ = st.RemoveOrphans()
		}
	}

	return removed, nil
}

type SnapshotManifest struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	SourcePath  string    `json:"source_path"`
	Deduplicate bool      `json:"deduplicate"`
}

func writeManifest(path string, m SnapshotManifest) error {
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

func readManifest(path string) (SnapshotManifest, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return SnapshotManifest{}, err
	}
	var m SnapshotManifest
	if err := json.Unmarshal(b, &m); err != nil {
		return SnapshotManifest{}, err
	}
	return m, nil
}

func newSnapshotID() string {
	stamp := time.Now().UTC().Format("20060102T150405Z")
	suffix := make([]byte, 4)
	_, _ = rand.Read(suffix)
	return fmt.Sprintf("%s-%s", stamp, hex.EncodeToString(suffix))
}

// DefaultCopyOptions returns sane defaults for snapshot copy/restore.
//
// Example:
//
//	opts := backup.DefaultCopyOptions()
//	_ = opts
func DefaultCopyOptions() fs.CopyTreeOptions {
	return fs.CopyTreeOptions{Workers: 2}
}
