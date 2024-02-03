package tests

import (
	"os"
	"os/exec"
	"testing"

	"github.com/vanilla-os/sdk/pkg/v1/fs"
)

func TestIsMounted(t *testing.T) {
	mounted, err := fs.IsMounted("tmpfs", "/tmp")
	if err != nil {
		t.Errorf("Error checking if /tmp is mounted: %v", err)
		return
	}
	t.Logf("/tmp is mounted: %v", mounted)
}

func TestMountFuseOverlay(t *testing.T) {
	// skip if fuse-overlayfs is not available
	if _, err := exec.LookPath("fuse-overlayfs"); err != nil {
		t.Skip("fuse-overlayfs is not available")
	}

	// preparing a temporary directory
	tmpDir := t.TempDir()
	mountDir := tmpDir + "/overlay"
	lowerDir := tmpDir + "/lower"
	upperDir := tmpDir + "/upper"
	workDir := tmpDir + "/work"

	// creating the directories
	if err := os.MkdirAll(mountDir, 0o755); err != nil {
		t.Fatalf("Error creating mount directory: %v", err)
	}
	if err := os.MkdirAll(lowerDir, 0o755); err != nil {
		t.Fatalf("Error creating lower directory: %v", err)
	}
	if err := os.MkdirAll(upperDir, 0o755); err != nil {
		t.Fatalf("Error creating upper directory: %v", err)
	}
	if err := os.MkdirAll(workDir, 0o755); err != nil {
		t.Fatalf("Error creating work directory: %v", err)
	}

	// mounting the overlay filesystem
	if err := fs.MountFuseOverlay(mountDir, lowerDir, upperDir, workDir); err != nil {
		t.Fatalf("Error mounting overlay filesystem: %v", err)
	}

	// checking if the overlay filesystem is mounted
	mounted, err := fs.IsMounted("fuse-overlayfs", mountDir)
	if err != nil {
		t.Errorf("Error checking if overlay filesystem is mounted: %v", err)
		return
	}
	t.Logf("Overlay filesystem is mounted: %v", mounted)

	// unmounting the overlay filesystem
	if err := fs.UnmountFuseOverlay(mountDir); err != nil {
		t.Fatalf("Error unmounting overlay filesystem: %v", err)
	}
}
