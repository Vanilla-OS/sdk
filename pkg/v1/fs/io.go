package fs

import (
	"os"
)

// CopyFile copies the contents of one file to another
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
func MoveFile(sourcePath, destinationPath string) error {
	err := os.Rename(sourcePath, destinationPath)
	if err != nil {
		return err
	}
	return nil
}

// DeleteFile deletes a file
func DeleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}

// CreateDirectory creates a new directory
func CreateDirectory(directoryPath string) error {
	err := os.Mkdir(directoryPath, 0755)
	if err != nil {
		return err
	}
	return nil
}

// DeleteDirectory deletes a directory and its contents
func DeleteDirectory(directoryPath string) error {
	err := os.RemoveAll(directoryPath)
	if err != nil {
		return err
	}
	return nil
}
