package types

// BaseInfo contains the common fields for both disks and partitions
type BaseInfo struct {
	Path       string
	Size       int64
	HumanSize  string
	Filesystem string
	Mountpoint string
	Label      string
	UUID       string
	PARTUUID   string
}

// PartitionInfo represents information about a disk partition
type PartitionInfo struct {
	BaseInfo // Embedded struct, inheriting fields from BaseInfo
}

// DiskInfo represents information about a disk
type DiskInfo struct {
	BaseInfo   // Embedded struct, inheriting fields from BaseInfo
	Partitions []PartitionInfo
}
