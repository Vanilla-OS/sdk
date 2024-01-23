package tests

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/vanilla-os/sdk/pkg/v1/chroot"
)

// Note: this function tests both EnterChroot, ExitChroot and RunChroot, so
// there is no need to test them separately.
func TestRunChroot(t *testing.T) {
	// Skip the test if not running as root
	if os.Geteuid() != 0 {
		t.Skip("Skipping test: root privileges required for chroot operations")
	}

	// Download Alpine tarball to use as the new root filesystem
	alpineURL := "https://dl-cdn.alpinelinux.org/alpine/v3.19/releases/x86_64/alpine-minirootfs-3.19.0-x86_64.tar.gz"
	tmpDir := t.TempDir()
	alpineTarball := filepath.Join(tmpDir, "alpine-minirootfs.tar.gz")

	err := downloadFile(alpineURL, alpineTarball)
	if err != nil {
		t.Fatalf("Failed to download Alpine tarball: %v", err)
	}
	alpineRoot := filepath.Join(tmpDir, "alpine-root")
	err = os.Mkdir(alpineRoot, 0755)
	if err != nil {
		t.Fatalf("Failed to create Alpine root directory: %v", err)
	}
	err = extractTarball(alpineTarball, alpineRoot)
	if err != nil {
		t.Fatalf("Failed to extract Alpine tarball: %v", err)
	}

	// Run a function in the chroot root
	fErr, err := chroot.RunChroot(alpineRoot, func() error {
		// Check if we are in the new root by reading /etc/os-release
		osReleaseFile := filepath.Join("/etc", "os-release")
		osRelease, err := os.ReadFile(osReleaseFile)
		if err != nil {
			return err
		}
		t.Logf("Distribution (from inside): %s", osRelease)

		return nil
	})
	if err != nil {
		t.Fatalf("Failed to run function in chroot root: %v", err)
	}
	if fErr != nil {
		t.Fatalf("Function in chroot root failed: %v", fErr)
	}

	// Check if we are still in the old root by reading /etc/os-release
	osReleaseFile := filepath.Join("/etc", "os-release")
	osRelease, err := os.ReadFile(osReleaseFile)
	if err != nil {
		t.Fatalf("Failed to read /etc/os-release: %v", err)
	}
	t.Logf("Distribution (from outside): %s", osRelease)
}

func downloadFile(url, destination string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	file, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	return err
}

func extractTarball(tarball, destination string) error {
	cmd := exec.Command("tar", "xzf", tarball, "-C", destination)
	return cmd.Run()
}
