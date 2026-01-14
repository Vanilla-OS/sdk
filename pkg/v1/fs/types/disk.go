package types

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

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
