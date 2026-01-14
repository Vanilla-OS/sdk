//go:build !linux || !cgo

package fs

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

// GetFilesystemInfo returns information about the filesystem of a file or
// partition by reading from /etc/mtab.
func GetFilesystemInfo(path string) map[string]string {
	return map[string]string{
		"LABEL":    "",
		"TYPE":     "",
		"UUID":     "",
		"PARTUUID": "",
	}
}
