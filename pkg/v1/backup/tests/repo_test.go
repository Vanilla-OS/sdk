package tests

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vanilla-os/sdk/pkg/v1/backup"
)

func TestRepositorySnapshotLifecycle(t *testing.T) {
	repoRoot := t.TempDir()
	repo, err := backup.OpenRepository(repoRoot)
	if err != nil {
		t.Fatalf("open repo: %v", err)
	}

	src := t.TempDir()
	if err := os.MkdirAll(filepath.Join(src, "d"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(src, "d", "f"), []byte("x"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	snap, err := repo.CreateSnapshot(src, backup.CreateSnapshotOptions{Deduplicate: true, DedupWorkers: 1})
	if err != nil {
		t.Fatalf("create snapshot: %v", err)
	}

	dst := t.TempDir()
	out := filepath.Join(dst, "restore")
	if err := repo.RestoreSnapshot(snap.Manifest.ID, out, backup.DefaultCopyOptions()); err != nil {
		t.Fatalf("restore snapshot: %v", err)
	}

	b, err := os.ReadFile(filepath.Join(out, "d", "f"))
	if err != nil {
		t.Fatalf("read restored: %v", err)
	}
	if string(b) != "x" {
		t.Fatalf("unexpected restored content")
	}

	removed, err := repo.PruneKeepLast(0)
	if err != nil {
		t.Fatalf("prune: %v", err)
	}
	if len(removed) != 1 {
		t.Fatalf("expected 1 removed")
	}
}
