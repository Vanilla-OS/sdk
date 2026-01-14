package tests

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

import (
	"testing"

	"github.com/vanilla-os/sdk/pkg/v1/system"
)

func TestGetSystemInfo(t *testing.T) {
	systemInfo, err := system.GetSystemInfo()
	if err != nil {
		t.Errorf("Error getting system info: %v", err)
		return
	}

	if systemInfo.OS == "" {
		t.Errorf("OS is empty")
	}

	if systemInfo.Version == "" {
		t.Errorf("Version is empty")
	}

	if systemInfo.Codename == "" {
		t.Errorf("Codename is empty")
	}

	if systemInfo.Arch == "" {
		t.Errorf("Arch is empty")
	}

	if systemInfo.MachineType == "" {
		t.Errorf("MachineType is empty")
	}

	t.Logf("OS: %s", systemInfo.OS)
	t.Logf("Version: %s", systemInfo.Version)
	t.Logf("Codename: %s", systemInfo.Codename)
	t.Logf("Arch: %s", systemInfo.Arch)
	t.Logf("MachineType: %s", systemInfo.MachineType)
}
