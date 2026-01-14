package fs

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// IsMounted checks if the given source path is mounted in the given
// destination path. It does so by reading the /proc/mounts file.
//
// Example:
//
//	mounted, err := fs.IsMounted("tmpfs", "/tmp")
//	if err != nil {
//		fmt.Printf("Error checking if /tmp is mounted: %v", err)
//		return
//	}
//	fmt.Printf("/tmp is mounted: %v", mounted)
func IsMounted(source string, destination string) (bool, error) {
	mounts, err := os.Open("/proc/mounts")
	if err != nil {
		return false, fmt.Errorf("error opening /proc/mounts: %w", err)
	}
	defer mounts.Close()

	scanner := bufio.NewScanner(mounts)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, source) && strings.Contains(line, destination) {
			return true, nil
		}
	}

	return false, nil
}

// Mount mounts the given source path in the given destination path. It also
// creates the destination path if it does not exist. An error is returned if
// the source path does not exist.
//
// Notes (for developers):
//
// Avoid mapping the fsType into a custom type, as it would require further
// maintenance once a new filesystem type is added to Vanilla OS.
//
// Example:
//
//	err := fs.Mount("/dev/sda1", "/mnt", "ext4", "", syscall.MS_RDONLY)
//	if err != nil {
//		fmt.Printf("Error mounting /dev/sda1: %v", err)
//		return
//	}
func Mount(source, destination, fsType, data string, mode uintptr) error {
	info, err := os.Stat(source)
	if err != nil {
		return err
	}

	if info.IsDir() {
		_ = os.MkdirAll(destination, 0o755)
	} else {
		file, _ := os.Create(destination)

		defer func() { _ = file.Close() }()
	}

	return syscall.Mount(source, destination, fsType, mode, data)
}

// MountBind mounts bind the given source path in the given destination path.
//
// Notes:
//
// This is just a wrapper of the Mount function, for convenience.
//
// Example:
//
//	err := fs.MountBind("/bruce", "/batman")
//	if err != nil {
//		fmt.Printf("Error bind mounting /batman: %v", err)
//		return
//	}
func MountBind(src, dest string) error {
	return Mount(
		src,
		dest,
		"bind",
		"",
		syscall.MS_BIND|
			syscall.MS_REC|
			syscall.MS_RDONLY|
			syscall.MS_NOSUID|
			syscall.MS_NOEXEC|
			syscall.MS_NODEV|
			syscall.MS_PRIVATE,
	)
}

// MountOverlay mounts the given lower, upper and work directories in the
// given destination path as an overlay filesystem.
//
// Notes:
//
// This is just a wrapper of the Mount function, for convenience.
//
// Example:
//
//	err := fs.MountOverlay("/batman/lower", "/batman/upper", "/batman/work")
//	if err != nil {
//		fmt.Printf("Error overlay mounting /batman: %v", err)
//		return
//	}
func MountOverlay(lowerDir, upperDir, workDir string) error {
	return Mount(
		"overlay",
		lowerDir,
		"overlay",
		fmt.Sprintf(
			"lowerdir=%s,upperdir=%s,workdir=%s,userxattr",
			lowerDir, upperDir, workDir,
		),
		0,
	)
}

// MountFuseOverlay mounts the given lower, upper and work directories in the
// given destination path as an overlay filesystem using fuse-overlayfs.
//
// Notes:
//
// This implementation uses the fuse-overlayfs command-line tool, if that
// is not available in the system, this function will return an error.
//
// Example:
//
//	err := fs.MountFuseOverlay("/batman", "/batman/lower", "/batman/upper", "/batman/work")
//	if err != nil {
//		fmt.Printf("Error fuse-overlayfs mounting /batman: %v", err)
//		return
//	}
func MountFuseOverlay(targetDir, lowerDir, upperDir, workDir string) (err error) {
	// TODO: Replace the command-line tool with a Go library, if available.
	c := exec.Command(
		"fuse-overlayfs",
		targetDir,
		"-o",
		fmt.Sprintf(
			"lowerdir=%s,upperdir=%s,workdir=%s",
			lowerDir, upperDir, workDir,
		),
	)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

// Unmount unmounts the given path. An error is returned if the path is not
// mounted.
//
// Example:
//
//	err := fs.Unmount("/mnt")
//	if err != nil {
//		fmt.Printf("Error unmounting /mnt: %v", err)
//		return
//	}
func Unmount(target string) error {
	return syscall.Unmount(target, 0)
}

// UnmountFuseOverlay unmounts the given path using the fuse-overlayfs command-
// line tool. An error is returned if the path is not mounted.
//
// Notes:
//
// This implementation uses the fuse-overlayfs command-line tool, if that
// is not available in the system, this function will return an error.
//
// Example:
//
//	err := fs.UnmountFuseOverlay("/batman")
//	if err != nil {
//		fmt.Printf("Error fuse-overlayfs unmounting /batman: %v", err)
//		return
//	}
func UnmountFuseOverlay(targetDir string) (err error) {
	// TODO: Replace the command-line tool with a Go library, if available.
	c := exec.Command("fusermount", "-u", targetDir)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
