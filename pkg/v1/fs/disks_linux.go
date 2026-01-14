//go:build linux && cgo

package fs

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

import (
	"path/filepath"

	"github.com/jochenvg/go-udev"
)

// GetFilesystemInfo returns information about the filesystem of a file or
// partition by reading from /etc/mtab.
func GetFilesystemInfo(path string) map[string]string {
	info := make(map[string]string)

	u := udev.Udev{}
	d := u.NewDeviceFromSyspath(filepath.Join("/sys/class/block", filepath.Base(path)))
	info["LABEL"] = d.PropertyValue("ID_FS_LABEL")
	info["TYPE"] = d.PropertyValue("ID_FS_TYPE")
	info["UUID"] = d.PropertyValue("ID_FS_UUID")
	info["PARTUUID"] = d.PropertyValue("ID_PART_ENTRY_UUID")
	return info
}
