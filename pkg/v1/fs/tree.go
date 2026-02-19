package fs

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// CopyTreeOptions controls CopyTree behavior.
type CopyTreeOptions struct {
	Workers             int
	PreserveOwnership   bool
	PreserveTimestamps  bool
	PreservePermissions bool
	AllowSpecial        bool
	OnProgress          func(CopyTreeProgress)
}

// CopyTreeProgress is emitted via CopyTreeOptions.OnProgress when set.
type CopyTreeProgress struct {
	SourcePath      string
	DestinationPath string
	BytesCopied     int64
}

type dirMeta struct {
	path string
	mode os.FileMode
	uid  int
	gid  int
	at   time.Time
	mt   time.Time
}

// CopyTree copies a directory tree from source to destination.
//
// Notes:
// - It is streaming (does not load whole files into memory).
// - It preserves symlinks (does not follow them).
// - It preserves permissions/ownership/timestamps when requested.
//
// Example:
//
//	err := fs.CopyTree("/source", "/destination", fs.CopyTreeOptions{Workers: 4})
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
func CopyTree(source, destination string, opts CopyTreeOptions) error {
	if opts.Workers <= 0 {
		opts.Workers = runtime.GOMAXPROCS(0)
		if opts.Workers < 1 {
			opts.Workers = 1
		}
	}
	if !opts.PreservePermissions {
		opts.PreservePermissions = true
	}
	if !opts.PreserveOwnership {
		opts.PreserveOwnership = true
	}
	if !opts.PreserveTimestamps {
		opts.PreserveTimestamps = true
	}

	if err := os.MkdirAll(destination, 0o755); err != nil {
		return err
	}

	jobs := make(chan fileJob, opts.Workers*2)
	var wg sync.WaitGroup

	var firstErr error
	var errMu sync.Mutex
	setErr := func(err error) {
		errMu.Lock()
		defer errMu.Unlock()
		if firstErr == nil {
			firstErr = err
		}
	}

	progress := func(p CopyTreeProgress) {
		if opts.OnProgress != nil {
			opts.OnProgress(p)
		}
	}

	for i := 0; i < opts.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				if err := copyRegularFile(j.src, j.dst, j.mode, j.uid, j.gid, j.at, j.mt, opts, progress); err != nil {
					setErr(err)
				}
			}
		}()
	}

	dirs := make([]dirMeta, 0)

	walkErr := filepath.WalkDir(source, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		dst := filepath.Join(destination, rel)

		info, err := d.Info()
		if err != nil {
			return err
		}

		mode := info.Mode()
		uid, gid, at, mt := extractUnixMeta(info)

		switch {
		case mode.IsDir():
			if err := os.MkdirAll(dst, 0o755); err != nil {
				return err
			}
			dirs = append(dirs, dirMeta{path: dst, mode: mode, uid: uid, gid: gid, at: at, mt: mt})
			return nil

		case mode&os.ModeSymlink != 0:
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			if err := os.Symlink(link, dst); err != nil {
				if errors.Is(err, os.ErrExist) {
					_ = os.Remove(dst)
					return os.Symlink(link, dst)
				}
				return err
			}
			if opts.PreserveOwnership {
				_ = os.Lchown(dst, uid, gid)
			}
			return nil

		case mode.IsRegular():
			jobs <- fileJob{src: path, dst: dst, mode: mode, uid: uid, gid: gid, at: at, mt: mt}
			return nil

		default:
			if opts.AllowSpecial {
				// Special files support can be added later via mknod.
				return nil
			}
			return fmt.Errorf("unsupported file type: %s", path)
		}
	})

	close(jobs)
	wg.Wait()

	if walkErr != nil {
		return walkErr
	}
	if firstErr != nil {
		return firstErr
	}

	// Apply dir metadata bottom-up to keep directory mtimes stable.
	for i := len(dirs) - 1; i >= 0; i-- {
		dm := dirs[i]
		if opts.PreservePermissions {
			_ = os.Chmod(dm.path, dm.mode)
		}
		if opts.PreserveOwnership {
			_ = os.Chown(dm.path, dm.uid, dm.gid)
		}
		if opts.PreserveTimestamps {
			_ = os.Chtimes(dm.path, dm.at, dm.mt)
		}
	}

	return nil
}

type fileJob struct {
	src  string
	dst  string
	mode os.FileMode
	uid  int
	gid  int
	at   time.Time
	mt   time.Time
}

func copyRegularFile(src, dst string, mode os.FileMode, uid, gid int, at, mt time.Time, opts CopyTreeOptions, progress func(CopyTreeProgress)) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}

	buf := make([]byte, 1024*1024)
	n, copyErr := io.CopyBuffer(out, in, buf)
	closeErr := out.Close()
	if copyErr != nil {
		return copyErr
	}
	if closeErr != nil {
		return closeErr
	}

	if opts.PreservePermissions {
		_ = os.Chmod(dst, mode)
	}
	if opts.PreserveOwnership {
		_ = os.Chown(dst, uid, gid)
	}
	if opts.PreserveTimestamps {
		_ = os.Chtimes(dst, at, mt)
	}

	progress(CopyTreeProgress{SourcePath: src, DestinationPath: dst, BytesCopied: n})
	return nil
}
