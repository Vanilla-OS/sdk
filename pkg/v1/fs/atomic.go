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

	"golang.org/x/sys/unix"
)

// AtomicSwap atomically swaps two files or directories
//
// Example:
//
//	err := fs.AtomicSwap("/tmp/file1", "/tmp/file2")
//	if err != nil {
//		fmt.Printf("Error swapping files: %v", err)
//		return
//	}
func AtomicSwap(sourcePath, destinationPath string) error {
	orig, err := os.Open(sourcePath)
	if err != nil {
		return err
	}

	newfile, err := os.Open(destinationPath)
	if err != nil {
		return err
	}

	err = unix.Renameat2(int(orig.Fd()), sourcePath, int(newfile.Fd()), destinationPath, unix.RENAME_EXCHANGE)
	if err != nil {
		return err
	}

	return nil
}
