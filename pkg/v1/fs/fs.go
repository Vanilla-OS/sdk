package fs

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/vanilla-os/sdk/pkg/v1/fs/types"
)

// GetFileList returns a list of files in the specified directory.
// If recursive is true, the function will recursively search for files, if
// fullPaths is true, the full path of the file will be returned instead of
// the relative path.
//
// Example:
//
//	fileList, err := fs.GetFileList("/batmans/cave", true, false)
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//
//	for _, file := range fileList {
//		fmt.Printf("Path: %s\n", file.Path)
//		fmt.Printf("Size: %d\n", file.Size)
//		fmt.Printf("Permissions: %s\n", file.Permissions.String())
//		fmt.Printf("Extension: %s\n", file.Extension)
//	}
func GetFileList(directory string, recursive, fullPaths bool) ([]types.FileInfo, error) {
	var fileList []types.FileInfo

	walkFunc := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !recursive && path != directory {
			return fs.SkipDir
		}
		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		filePath := path
		if !fullPaths {
			filePath, _ = filepath.Rel(directory, path)
		}

		permissions := convertPermissions(info.Mode())

		fileList = append(fileList, types.FileInfo{
			Path:        filePath,
			ParentPath:  filepath.Dir(filePath),
			IsDirectory: false,
			Size:        info.Size(),
			Permissions: permissions,
			Extension:   GetFileExtension(path),
		})
		return nil
	}

	if err := filepath.WalkDir(directory, walkFunc); err != nil {
		return nil, err
	}

	return fileList, nil
}

// convertPermissions converts file mode to Permission type
func convertPermissions(mode os.FileMode) types.Permission {
	return types.Permission{
		OwnerRead:     mode&0400 != 0,
		OwnerWrite:    mode&0200 != 0,
		OwnerExecute:  mode&0100 != 0,
		GroupRead:     mode&0040 != 0,
		GroupWrite:    mode&0020 != 0,
		GroupExecute:  mode&0010 != 0,
		OthersRead:    mode&0004 != 0,
		OthersWrite:   mode&0002 != 0,
		OthersExecute: mode&0001 != 0,
	}
}

// GetFile returns the file info of the specified file.
// If fullPath is true, the full path of the file will be returned instead of
// the relative path.
//
// Example:
//
//	file, err := fs.GetFile("/batmans/cave/batmobile.txt", false)
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//
//	fmt.Printf("Path: %s\n", file.Path)
//	fmt.Printf("Size: %d\n", file.Size)
//	fmt.Printf("Permissions: %s\n", file.Permissions.String())
func GetFile(filePath string, fullPath bool) (types.FileInfo, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return types.FileInfo{}, err
	}

	permissions := convertPermissions(info.Mode())
	file := types.FileInfo{
		Path:        filePath,
		ParentPath:  filepath.Dir(filePath),
		IsDirectory: info.IsDir(),
		Size:        info.Size(),
		Permissions: permissions,
		Extension:   GetFileExtension(filePath),
	}

	if !fullPath {
		file.Path, _ = filepath.Rel(".", filePath)
	}

	return file, nil
}

// GetFileExtension returns the extension of the given file
func GetFileExtension(filePath string) string {
	ext := filepath.Ext(filePath)
	if ext != "" {
		return strings.TrimPrefix(ext, ".")
	}
	return ""
}

// IsFile checks whether the given path is a regular file
func IsFile(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}

// IsDirectory checks whether the given path is a directory
func IsDirectory(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if err != nil {
		return false
	}
	return info.Mode().IsDir()
}

// GetFileSize returns the size of the specified file in bytes
func GetFileSize(filePath string) int64 {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0
	}
	return info.Size()
}

// WriteFileContent writes content to the specified file
func WriteFileContent(filePath, content string) error {
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}

// FileExists checks if a file exists.
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// DirectoryExists checks if a directory exists.
func DirectoryExists(directoryPath string) bool {
	fileInfo, err := os.Stat(directoryPath)
	return err == nil && fileInfo.IsDir()
}

// ListDirectories returns a list of directories in the specified directory.
func ListDirectories(directoryPath string) ([]string, error) {
	var directories []string
	files, err := os.ReadDir(directoryPath)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			directories = append(directories, file.Name())
		}
	}
	return directories, nil
}
