//go:build linux

package fs

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func resolveDiskSymlink(dir, value string) (string, error) {
	p := filepath.Join(dir, value)
	_, err := os.Lstat(p)
	if err != nil {
		return "", err
	}
	resolved, err := filepath.EvalSymlinks(p)
	if err != nil {
		return "", err
	}
	return resolved, nil
}

// GetDeviceByUUID resolves /dev/disk/by-uuid/<uuid>.
//
// Example:
//
//	dev, err := fs.GetDeviceByUUID("5a1b2c3d-4e5f-6789-abcd-ef0123456789")
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//	fmt.Printf("Device: %s\n", dev)
func GetDeviceByUUID(uuid string) (string, error) {
	return resolveDiskSymlink("/dev/disk/by-uuid", uuid)
}

// GetDeviceByPARTUUID resolves /dev/disk/by-partuuid/<partuuid>.
//
// Example:
//
//	dev, err := fs.GetDeviceByPARTUUID("aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee")
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//	fmt.Printf("Device: %s\n", dev)
func GetDeviceByPARTUUID(partuuid string) (string, error) {
	return resolveDiskSymlink("/dev/disk/by-partuuid", partuuid)
}

// GetDeviceByLabel resolves /dev/disk/by-label/<label>.
//
// Example:
//
//	dev, err := fs.GetDeviceByLabel("VANILLA")
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//	fmt.Printf("Device: %s\n", dev)
func GetDeviceByLabel(label string) (string, error) {
	return resolveDiskSymlink("/dev/disk/by-label", label)
}

// IsRemovableDevice checks if a block device is marked as removable via sysfs.
//
// Example:
//
//	rem, err := fs.IsRemovableDevice("/dev/sda")
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//	fmt.Printf("Removable: %v\n", rem)
func IsRemovableDevice(devicePath string) (bool, error) {
	name := filepath.Base(devicePath)
	b, err := os.ReadFile(filepath.Join("/sys/class/block", name, "removable"))
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(b)) == "1", nil
}

// GetDeviceSysPath returns the sysfs path for a /dev/* device.
//
// Example:
//
//	sp, err := fs.GetDeviceSysPath("/dev/sda")
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//	fmt.Printf("Sysfs path: %s\n", sp)
func GetDeviceSysPath(devicePath string) (string, error) {
	name := filepath.Base(devicePath)
	p := filepath.Join("/sys/class/block", name)
	if _, err := os.Stat(p); err != nil {
		return "", fmt.Errorf("device not found in sysfs: %w", err)
	}
	return p, nil
}
