//go:build linux

package fs

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

import (
	"os"
	"syscall"
	"time"
)

func extractUnixMeta(info os.FileInfo) (uid, gid int, at, mt time.Time) {
	mt = info.ModTime()
	at = mt

	st, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, 0, at, mt
	}

	uid = int(st.Uid)
	gid = int(st.Gid)

	at = time.Unix(int64(st.Atim.Sec), int64(st.Atim.Nsec))
	mt = time.Unix(int64(st.Mtim.Sec), int64(st.Mtim.Nsec))
	return uid, gid, at, mt
}
