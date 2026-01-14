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
)

// Notes:
// 	Why do we provide those wrappers?
//	We provide wrappers for built-in functions for two primary reasons:
//		1. To provide a standard interface for common operations, making
//		   it easier for developers to understand and investigate the code.
//		2. (not yet) The SDK, if used to create a new VOS App instance, will
//		   be able to log all the application's operations in the system log,
//		   making it easier to debug and monitor the application by both the
//		   developer and the system administrator.

// CopyFile copies the contents of one file to another
//
// Example:
//
//	err := fs.CopyFile("/tmp/batman", "/tmp/robin")
//	if err != nil {
//		fmt.Printf("Error copying file: %v", err)
//		return
//	}
func CopyFile(sourcePath, destinationPath string) error {
	source, err := os.ReadFile(sourcePath)
	if err != nil {
		return err
	}
	err = os.WriteFile(destinationPath, source, 0644)
	if err != nil {
		return err
	}
	return nil
}

// MoveFile moves a file from one location to another
//
// Example:
//
//	err := fs.MoveFile("/tmp/batman", "/tmp/robin")
//	if err != nil {
//		fmt.Printf("Error moving file: %v", err)
//		return
//	}
func MoveFile(sourcePath, destinationPath string) error {
	err := os.Rename(sourcePath, destinationPath)
	if err != nil {
		return err
	}
	return nil
}

// DeleteFile deletes a file
//
// Example:
//
//	err := fs.DeleteFile("/tmp/batman")
//	if err != nil {
//		fmt.Printf("Error deleting file: %v", err)
//		return
//	}
func DeleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}

// CreateDirectory creates a new directory
//
// Example:
//
//	err := fs.CreateDirectory("/tmp/batman")
//	if err != nil {
//		fmt.Printf("Error creating directory: %v", err)
//		return
//	}
func CreateDirectory(directoryPath string) error {
	err := os.Mkdir(directoryPath, 0755)
	if err != nil {
		return err
	}
	return nil
}

// DeleteDirectory deletes a directory and its contents
//
// Example:
//
//	err := fs.DeleteDirectory("/tmp/batman")
//	if err != nil {
//		fmt.Printf("Error deleting directory: %v", err)
//		return
//	}
func DeleteDirectory(directoryPath string) error {
	err := os.RemoveAll(directoryPath)
	if err != nil {
		return err
	}
	return nil
}
