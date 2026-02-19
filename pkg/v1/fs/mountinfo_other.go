//go:build !linux

package fs

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

import "fmt"

// MountInfo is not supported on non-Linux platforms.
type MountInfo struct{}

func GetMountInfo() ([]MountInfo, error) {
	return nil, fmt.Errorf("mountinfo is only supported on Linux")
}

func GetMountpoint(source string) (string, error) {
	return "", fmt.Errorf("mountinfo is only supported on Linux")
}
