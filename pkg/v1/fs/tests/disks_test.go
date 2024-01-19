package tests

import (
	"testing"

	"github.com/vanilla-os/sdk/pkg/v1/fs"
)

func TestGetDiskList(t *testing.T) {
	disks, err := fs.GetDiskList()
	if err != nil {
		t.Errorf("Error getting disk list: %v", err)
		return
	}

	if len(disks) == 0 {
		t.Errorf("No disks found")
		return
	}

	for _, disk := range disks {
		t.Logf("Disk: %s", disk.Path)
		t.Logf("Size: %d", disk.Size)
		t.Logf("HumanSize: %s", disk.HumanSize)
		for _, partition := range disk.Partitions {
			t.Logf("Partition: %s", partition.Path)
			t.Logf("Size: %d", partition.Size)
			t.Logf("HumanSize: %s", partition.HumanSize)
		}
	}
}
