package system

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/shirou/gopsutil/host"
	"github.com/vanilla-os/sdk/pkg/v1/system/types"
)

// GetSystemInfo returns information about the system, such as OS, version,
// codename, architecture and machine type. If the machine type cannot be
// determined, it will be set to BareMetal. If any error occurs, it will
// be returned.
//
// Example:
//
//	info, err := system.GetSystemInfo()
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//	fmt.Printf("OS: %s\n", info.OS)
func GetSystemInfo() (*types.SystemInfo, error) {
	hostInfo, err := host.Info()
	if err != nil {
		return nil, err
	}

	osReleaseInfo, err := readOSRelease()
	if err != nil {
		return nil, err
	}

	machineType, err := getMachineType()
	if err != nil {
		return nil, err
	}

	info := &types.SystemInfo{
		OS:          osReleaseInfo.Name,
		Version:     osReleaseInfo.Version,
		Codename:    osReleaseInfo.Codename,
		Arch:        hostInfo.KernelArch,
		MachineType: machineType,
	}
	return info, nil
}

// readOSRelease reads the /etc/os-release file and returns the OS name,
// version and codename. In the future releases on Vanilla OS, we may
// consider storing this information in a different file, perhaps using
// a better format. If any error occurs, it will be returned.
func readOSRelease() (*types.OSReleaseInfo, error) {
	osReleaseInfo := &types.OSReleaseInfo{}

	file, err := os.Open("/etc/os-release")
	if err != nil {
		return nil, fmt.Errorf("failed to open /etc/os-release: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.Trim(strings.TrimSpace(parts[1]), "\"")

			switch key {
			case "NAME":
				osReleaseInfo.Name = value
			case "VERSION_ID":
				osReleaseInfo.Version = value
			case "VERSION_CODENAME":
				osReleaseInfo.Codename = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading /etc/os-release: %v", err)
	}

	if osReleaseInfo.Name == "" || osReleaseInfo.Version == "" {
		return nil, fmt.Errorf("missing or invalid information in /etc/os-release")
	}

	return osReleaseInfo, nil
}

// getMachineType returns the machine type, which can be BareMetal, VM or
// Container. If any error occurs, it will be returned.
func getMachineType() (types.MachineType, error) {
	// Check if hypervisor information is present in /proc/cpuinfo
	cpuInfo, err := os.ReadFile("/proc/cpuinfo")
	if err == nil && strings.Contains(string(cpuInfo), "hypervisor") {
		return types.VM, nil
	}

	// Check if /run/.containerenv file exists
	if _, err := os.Stat("/run/.containerenv"); err == nil {
		return types.Container, nil
	}

	// No clear indication of VM or Container, assuming BareMetal
	return types.BareMetal, nil
}
