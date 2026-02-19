//go:build linux

package backup

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
	"syscall"
)

type repoLock struct {
	f *os.File
}

func (l *repoLock) Close() error {
	if l == nil || l.f == nil {
		return nil
	}
	_ = syscall.Flock(int(l.f.Fd()), syscall.LOCK_UN)
	return l.f.Close()
}

func (r *Repository) lockShared() (*repoLock, error) {
	return r.lock(syscall.LOCK_SH)
}

func (r *Repository) lockExclusive() (*repoLock, error) {
	return r.lock(syscall.LOCK_EX)
}

func (r *Repository) lock(flag int) (*repoLock, error) {
	lockPath := filepath.Join(r.Root, ".repo.lock")
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, err
	}
	if err := syscall.Flock(int(f.Fd()), flag); err != nil {
		_ = f.Close()
		return nil, err
	}
	return &repoLock{f: f}, nil
}
