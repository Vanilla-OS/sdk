package types

// FileInfo represents information about a file
type FileInfo struct {
	// Path to the file
	Path string

	// ParentPath is the path to the parent directory of the file
	ParentPath string

	// IsDirectory is true if the file is a directory
	IsDirectory bool

	// Size of the file in bytes
	Size int64

	// Permissions of the file as a Permission struct
	Permissions Permission

	// Extension of the file (e.g. "txt")
	Extension string
}

// FileDiff represents the difference between two files.
type FileDiffInfo struct {
	// AddedLines are the lines added to the file compared to the other file
	AddedLines []string

	// RemovedLines are the lines removed from the file compared to the
	// other file
	RemovedLines []string
}
