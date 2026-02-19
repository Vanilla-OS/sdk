//go:build !linux

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
	"time"
)

func extractUnixMeta(info os.FileInfo) (uid, gid int, at, mt time.Time) {
	mt = info.ModTime()
	at = mt
	return 0, 0, at, mt
}
