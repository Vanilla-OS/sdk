//go:build linux && cgo

package fs

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
