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
	"path/filepath"
	"syscall"
)

// GetFreeSpaceBytes returns the available free space for the filesystem hosting path.
//
// Example:
//
//	free, err := fs.GetFreeSpaceBytes("/")
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//	fmt.Printf("Free bytes: %d\n", free)
func GetFreeSpaceBytes(path string) (uint64, error) {
	var st syscall.Statfs_t
	err := syscall.Statfs(path, &st)
	if err != nil {
		return 0, err
	}
	return st.Bavail * uint64(st.Bsize), nil
}

// GetTotalSpaceBytes returns the total size for the filesystem hosting path.
//
// Example:
//
//	total, err := fs.GetTotalSpaceBytes("/")
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//	fmt.Printf("Total bytes: %d\n", total)
func GetTotalSpaceBytes(path string) (uint64, error) {
	var st syscall.Statfs_t
	err := syscall.Statfs(path, &st)
	if err != nil {
		return 0, err
	}
	return st.Blocks * uint64(st.Bsize), nil
}

// IsWritableDir checks if dir is writable by attempting to create a temporary file.
//
// Example:
//
//	writable, err := fs.IsWritableDir("/tmp")
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//	fmt.Printf("Writable: %v\n", writable)
func IsWritableDir(dir string) (bool, error) {
	testPath := filepath.Join(dir, ".vos-sdk-writable-check")
	f, err := os.OpenFile(testPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		if os.IsPermission(err) {
			return false, nil
		}
		return false, err
	}
	_ = f.Close()
	_ = os.Remove(testPath)
	return true, nil
}
