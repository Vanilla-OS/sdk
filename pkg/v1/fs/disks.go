package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/vanilla-os/sdk/pkg/v1/fs/types"
)

// GetDiskList returns a list of disks on the system
//
// Example:
//
//	disks, err := fs.GetDiskList()
//	if err != nil {
//		fmt.Printf("Error getting disk list: %v", err)
//		return
//	}
//	for _, disk := range disks {
//		fmt.Printf("Disk: %s\n", disk.Path)
//		fmt.Printf("Size: %d\n", disk.Size)
//		fmt.Printf("HumanSize: %s\n", disk.HumanSize)
//		for _, partition := range disk.Partitions {
//			fmt.Printf("Partition: %s\n", partition.Path)
//			fmt.Printf("Size: %d\n", partition.Size)
//			fmt.Printf("HumanSize: %s\n", partition.HumanSize)
//		}
//	}
func GetDiskList() ([]types.DiskInfo, error) {
	var disks []types.DiskInfo

	files, err := os.ReadDir("/sys/class/block")
	if err != nil {
		return nil, err
	}

	diskMap := make(map[string]types.DiskInfo)

	for _, file := range files {
		name := file.Name()

		// Skip non-disk entries
		if strings.HasPrefix(name, "loop") ||
			strings.HasPrefix(name, "ram") ||
			strings.HasPrefix(name, "zram") ||
			strings.HasPrefix(name, "dm") {
			continue
		}

		diskPath := filepath.Join("/dev", name)
		diskInfo, err := GetDiskInfo(diskPath)
		if err != nil {
			return nil, err
		}

		partitions, err := GetPartitionList(diskPath)
		if err != nil {
			return nil, err
		}
		diskInfo.Partitions = partitions

		diskMap[diskPath] = diskInfo
	}

	for _, disk := range diskMap {
		disks = append(disks, disk)
	}

	return disks, nil
}

// GetPartitionList returns a list of disk partitions on the specified disk
//
// Example:
//
//	partitions, err := fs.GetPartitionList("/dev/sda")
//	if err != nil {
//		fmt.Printf("Error getting partition list: %v", err)
//		return
//	}
//	for _, partition := range partitions {
//		fmt.Printf("Partition: %s\n", partition.Path)
//		fmt.Printf("Size: %d\n", partition.Size)
//		fmt.Printf("HumanSize: %s\n", partition.HumanSize)
//	}
func GetPartitionList(diskPath string) ([]types.PartitionInfo, error) {
	var partitions []types.PartitionInfo

	files, err := os.ReadDir(filepath.Join("/sys/class/block", filepath.Base(diskPath)))
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		name := file.Name()

		// Skip non-partition entries
		if !strings.HasPrefix(name, filepath.Base(diskPath)) {
			continue
		}

		partitionPath := filepath.Join("/dev", name)
		partitionInfo, err := GetPartitionInfo(partitionPath)
		if err != nil {
			return nil, err
		}

		partitions = append(partitions, partitionInfo)
	}

	return partitions, nil
}

// GetDiskInfo returns information about a specific disk
//
// Example:
//
//	diskInfo, err := fs.GetDiskInfo("/dev/sda")
//	if err != nil {
//		fmt.Printf("Error getting disk info: %v", err)
//		return
//	}
//	fmt.Printf("Path: %s\n", diskInfo.Path)
//	fmt.Printf("Size: %d\n", diskInfo.Size)
//	fmt.Printf("HumanSize: %s\n", diskInfo.HumanSize)
//	for _, partition := range diskInfo.Partitions {
//		fmt.Printf("Partition: %s\n", partition.Path)
//		fmt.Printf("Size: %d\n", partition.Size)
//		fmt.Printf("HumanSize: %s\n", partition.HumanSize)
//	}
func GetDiskInfo(diskPath string) (types.DiskInfo, error) {
	disk := types.DiskInfo{
		Path: diskPath,
	}

	sizePath := filepath.Join("/sys/class/block", filepath.Base(diskPath), "size")
	size, err := os.ReadFile(sizePath)
	if err != nil {
		return types.DiskInfo{}, err
	}

	sectorSize := 512
	disk.Size = int64(sectorSize) * int64(parseUint64(strings.TrimSpace(string(size))))
	disk.HumanSize = GetHumanSize(disk.Size)

	return disk, nil
}

// GetPartitionInfo returns information about a specific disk partition
//
// Example:
//
//	partitionInfo, err := fs.GetPartitionInfo("/dev/sda1")
//	if err != nil {
//		fmt.Printf("Error getting partition info: %v", err)
//		return
//	}
//	fmt.Printf("Path: %s\n", partitionInfo.Path)
//	fmt.Printf("Size: %d\n", partitionInfo.Size)
//	fmt.Printf("HumanSize: %s\n", partitionInfo.HumanSize)
func GetPartitionInfo(partitionPath string) (types.PartitionInfo, error) {
	partition := types.PartitionInfo{
		Path: partitionPath,
	}

	sizePath := filepath.Join("/sys/class/block", filepath.Base(partitionPath), "size")
	size, err := os.ReadFile(sizePath)
	if err != nil {
		return types.PartitionInfo{}, err
	}

	sectorSize := 512
	partition.Size = int64(sectorSize) * int64(parseUint64(strings.TrimSpace(string(size))))
	partition.HumanSize = GetHumanSize(partition.Size)

	return partition, nil
}

// parseUint64 parses a string into a uint64 or returns 0 if parsing fails
func parseUint64(s string) uint64 {
	value, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}
	return value
}

// GetHumanSize converts the size from bytes to a human-readable format.
// For example, 1024 bytes would be converted to "1 kB".
//
// Example:
//
//	fmt.Println(GetHumanSize(1024)) // 1 kB
func GetHumanSize(size int64) string {
	const (
		kB = 1024.0
		mB = kB * 1024.0
		gB = mB * 1024.0
		tB = gB * 1024.0
	)

	sizeFloat := float64(size)

	switch {
	case size < int64(kB):
		return fmt.Sprintf("%d B", size)
	case size < int64(mB):
		return fmt.Sprintf("%.2f kB", sizeFloat/kB)
	case size < int64(gB):
		return fmt.Sprintf("%.2f MB", sizeFloat/mB)
	case size < int64(tB):
		return fmt.Sprintf("%.2f GB", sizeFloat/gB)
	default:
		return fmt.Sprintf("%.2f TB", sizeFloat/tB)
	}
}
