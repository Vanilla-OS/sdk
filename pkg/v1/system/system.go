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
	// There are many ways to check if the system is running in a container,
	// we have to check multiple methods to be sure.

	// Check if /run/.containerenv file exists
	if _, err := os.Stat("/run/.containerenv"); err == nil {
		return types.Container, nil
	}

	// Check if /.dockerenv file exists
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return types.Container, nil
	}

	// Check if hypervisor information is present in /proc/cpuinfo
	cpuInfo, err := os.ReadFile("/proc/cpuinfo")
	if err == nil && strings.Contains(string(cpuInfo), "hypervisor") {
		return types.VM, nil
	}

	// No clear indication of VM or Container, assuming BareMetal
	return types.BareMetal, nil
}

// RunningInVM returns true if the system is running in a virtual machine,
// otherwise it returns false.
//
// Example:
//
//	if system.RunningInVM() {
//		fmt.Println("Running in a virtual machine")
//	} else {
//		fmt.Println("Not running in a virtual machine")
//	}
func RunningInVM() bool {
	info, err := getMachineType()
	if err != nil {
		return false
	}
	return info == types.VM
}

// RunningInContainer returns true if the system is running in a container,
// otherwise it returns false.
//
// Example:
//
//	if system.RunningInContainer() {
//		fmt.Println("Running in a container")
//	} else {
//		fmt.Println("Not running in a container")
//	}
func RunningInContainer() bool {
	info, err := getMachineType()
	if err != nil {
		return false
	}
	return info == types.Container
}

// RunningInBareMetal returns true if the system is running on bare metal,
// otherwise it returns false.
//
// Example:
//
//	if system.RunningInBareMetal() {
//		fmt.Println("Running on bare metal")
//	} else {
//		fmt.Println("Not running on bare metal")
//	}
func RunningInBareMetal() bool {
	info, err := getMachineType()
	if err != nil {
		return false
	}
	return info == types.BareMetal
}

// RunningInVMOrContainer returns true if the system is running in a virtual
// machine or a container, otherwise it returns false.
//
// Example:
//
//	if system.RunningInVMOrContainer() {
//		fmt.Println("Running in a virtual machine or a container")
//	} else {
//		fmt.Println("Not running in a virtual machine or a container")
//	}
func RunningInVMOrContainer() bool {
	info, err := getMachineType()
	if err != nil {
		return false
	}
	return info == types.VM || info == types.Container
}
