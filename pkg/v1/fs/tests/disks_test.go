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
		t.Skipf("No disks found")
		return
	}

	for _, disk := range disks {
		t.Logf("Disk: %s", disk.Path)
		t.Logf("Size: %d", disk.Size)
		t.Logf("HumanSize: %s", disk.HumanSize)
		t.Logf("Filesystem: %s", disk.Filesystem)
		t.Logf("Mountpoint: %s", disk.Mountpoint)
		t.Logf("Label: %s", disk.Label)
		t.Logf("UUID: %s", disk.UUID)
		if len(disk.Partitions) > 0 {
			t.Logf("Partitions:")
			for _, partition := range disk.Partitions {
				t.Logf("--------------------------------")
				t.Logf("Partition: %s", partition.Path)
				t.Logf("Size: %d", partition.Size)
				t.Logf("HumanSize: %s", partition.HumanSize)
				t.Logf("Filesystem: %s", partition.Filesystem)
				t.Logf("Mountpoint: %s", partition.Mountpoint)
				t.Logf("Label: %s", partition.Label)
				t.Logf("UUID: %s", partition.UUID)
				t.Logf("PARTUUID: %s", partition.PARTUUID)
			}
		}
		t.Logf("--------------------------------")
	}
}
