//go:build !linux || !cgo

package fs

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
