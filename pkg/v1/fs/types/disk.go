package types

// DiskInfo represents information about a disk
type DiskInfo struct {
	// Full path to the disk (e.g., "/dev/sda")
	Path string

	// Size of the disk in bytes
	Size int64

	// Partitions on the disk
	Partitions []PartitionInfo

	// HumanSize is the size of the disk in human-readable format
	HumanSize string
}

// PartitionInfo represents information about a disk partition
type PartitionInfo struct {
	// Full path to the partition (e.g., "/dev/sda1")
	Path string

	// Size of the partition in bytes
	Size int64

	// HumanSize is the size of the partition in human-readable format
	HumanSize string
}
