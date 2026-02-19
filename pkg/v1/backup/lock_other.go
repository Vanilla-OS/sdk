//go:build !linux

package backup

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

type repoLock struct{}

func (l *repoLock) Close() error { return nil }

func (r *Repository) lockShared() (*repoLock, error) { return &repoLock{}, nil }

func (r *Repository) lockExclusive() (*repoLock, error) { return &repoLock{}, nil }
