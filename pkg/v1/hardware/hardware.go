package hardware

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/vanilla-os/sdk/pkg/v1/hardware/types"
)

//go:embed resources/pci.ids
var pciIDsData string

var pciDeviceMap types.PCIDeviceMap

// GetPeripheralList returns a list of all system peripherals.
//
// Example:
//
//	peripherals, err := hardware.GetPeripheralList()
//	if err != nil {
//	  fmt.Printf("Error: %v\n", err)
//	  return
//	}
//
//	for _, peripheral := range peripherals {
//	  fmt.Printf("ID: %s\n", peripheral.ID)
//	  fmt.Printf("Name: %s\n", peripheral.Name)
//	  fmt.Printf("Type: %s\n", peripheral.Type)
//	}
func GetPeripheralList() ([]types.Peripheral, error) {
	peripherals := make([]types.Peripheral, 0)

	inputDevices, err := GetInputDevices()
	if err != nil {
		return nil, err
	}

	for _, inputDevice := range inputDevices {
		peripherals = append(peripherals, types.Peripheral{
			ID:   inputDevice.Product,
			Name: inputDevice.Name,
			Type: types.InputDeviceType,
		})
	}

	pciDevices, err := GetPCIDevices()
	if err != nil {
		return nil, err
	}

	for _, pciDevice := range pciDevices {
		peripherals = append(peripherals, types.Peripheral{
			ID:   pciDevice.ID,
			Name: pciDevice.Name,
			Type: types.PCIDeviceType,
		})
	}

	return peripherals, nil
}

// GetInputDevices returns a list of input devices with specific details.
//
// Example:
//
//	inputDevices, err := hardware.GetInputDevices()
//	if err != nil {
//	  fmt.Printf("Error: %v\n", err)
//	  return
//	}
//
//	for _, inputDevice := range inputDevices {
//	  fmt.Printf("Name: %s\n", inputDevice.Name)
//	  fmt.Printf("Product: %s\n", inputDevice.Product)
//	}
func GetInputDevices() ([]types.InputDevice, error) {
	inputDevices := make([]types.InputDevice, 0)

	inputDevicePaths, err := getInputDevicePaths()
	if err != nil {
		return nil, err
	}

	for _, path := range inputDevicePaths {
		deviceInfo, err := readInputDeviceInfo(path)
		if err != nil {
			return nil, err
		}

		inputDevices = append(inputDevices, deviceInfo)
	}

	return inputDevices, nil
}

// getInputDevicePaths returns a list of paths for input devices.
func getInputDevicePaths() ([]string, error) {
	var inputDevicePaths []string

	files, err := os.ReadDir("/sys/class/input")
	if err != nil {
		return nil, fmt.Errorf("error reading /sys/class/input: %v", err)
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "event") {
			inputDevicePaths = append(inputDevicePaths, fmt.Sprintf("/sys/class/input/%s/device", file.Name()))
		}
	}

	return inputDevicePaths, nil
}

// readInputDeviceInfo reads specific information about an input device.
func readInputDeviceInfo(devicePath string) (types.InputDevice, error) {
	deviceInfo := types.InputDevice{}

	ueventPath := fmt.Sprintf("%s/uevent", devicePath)
	file, err := os.Open(ueventPath)
	if err != nil {
		return deviceInfo, fmt.Errorf("error opening %s: %v", ueventPath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "PRODUCT=") {
			deviceInfo.Product = strings.TrimPrefix(line, "PRODUCT=")
		} else if strings.HasPrefix(line, "NAME=") {
			deviceInfo.Name = strings.TrimPrefix(line, "NAME=")
		}
	}

	return deviceInfo, nil
}

// GetPCIDevices returns a list of PCI devices from /sys/bus/pci/devices.
//
// Example:
//
//	pciDevices, err := hardware.GetPCIDevices()
//	if err != nil {
//	  fmt.Printf("Error: %v\n", err)
//	  return
//	}
//
//	for _, pciDevice := range pciDevices {
//	  fmt.Printf("ID: %s\n", pciDevice.ID)
//	}
func GetPCIDevices() ([]types.PCIDevice, error) {
	pciDevices := make([]types.PCIDevice, 0)

	files, err := os.ReadDir("/sys/bus/pci/devices")
	if err != nil {
		return nil, fmt.Errorf("error reading /sys/bus/pci/devices: %v", err)
	}

	for _, file := range files {
		pciDevice := types.PCIDevice{
			ID: file.Name(),
		}

		devicePath := fmt.Sprintf("/sys/bus/pci/devices/%s", file.Name())

		if vendor, err := readSysFile(devicePath, "vendor"); err == nil {
			pciDevice.Vendor = vendor
		}

		if device, err := readSysFile(devicePath, "device"); err == nil {
			pciDevice.Device = device
		}

		if class, err := readSysFile(devicePath, "class"); err == nil {
			pciDevice.Class = class
		}

		if subsystemDevice, err := readSysFile(devicePath, "subsystem_device"); err == nil {
			pciDevice.SubsystemDevice = subsystemDevice
		}

		if subsystemVendor, err := readSysFile(devicePath, "subsystem_vendor"); err == nil {
			pciDevice.SubsystemVendor = subsystemVendor
		}

		if modalias, err := readSysFile(devicePath, "modalias"); err == nil {
			pciDevice.Modalias = modalias
		}

		if device, vendorName, err := GetPCIDeviceByIDs(pciDevice.Vendor, pciDevice.Device); err == nil {
			pciDevice.Name = device.Name
			pciDevice.VendorName = vendorName
		}

		pciDevices = append(pciDevices, pciDevice)
	}

	return pciDevices, nil
}

// GetPCIDeviceByIDs returns a PCIDeviceMapDevice by vendor ID and device ID,
// useful for getting the name of a PCI device or the vendor name. It returns
// an error if the device is not found.
//
// Example:
//
//	device, vendorName, err := hardware.GetPCIDeviceByIDs("8086", "10f8")
//	if err != nil {
//	  fmt.Printf("Error: %v\n", err)
//	  return
//	}
//
//	fmt.Printf("Device name: %s\n", device.Name)
//	fmt.Printf("Vendor name: %s\n", vendorName)
func GetPCIDeviceByIDs(vendorID, deviceID string) (types.PCIDeviceMapDevice, string, error) {
	// Cleanup the vendor ID and device ID
	vendorID = strings.TrimPrefix(vendorID, "0x")
	deviceID = strings.TrimPrefix(deviceID, "0x")

	for _, vendor := range pciDeviceMap {
		if vendor.ID == vendorID {
			for _, device := range vendor.Devices {
				if device.ID == deviceID {
					return device, vendor.Name, nil
				}
			}
		}
	}

	return types.PCIDeviceMapDevice{}, "", fmt.Errorf("device not found for vendor ID: %s and device ID: %s", vendorID, deviceID)
}

// readSysFile reads the content of a file in the /sys directory.
func readSysFile(devicePath, fileName string) (string, error) {
	filePath := fmt.Sprintf("%s/%s", devicePath, fileName)
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening %s: %v", filePath, err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("error reading %s: %v", filePath, err)
	}

	return strings.TrimSpace(string(content)), nil
}

// LoadPCIDeviceMap loads the PCI IDs from /usr/share/misc/pci.ids. This function
// should be called at the beginning of the program or before using any other
// function from this package.
//
// # File syntax Syntax:
// # vendor  vendor_name
// #	device  device_name				<-- single tab
// #		subvendor subdevice  subsystem_name	<-- two tabs
// #	...
//
// Example:
//
//	err := hardware.Loadtypes.PCIDeviceMap()
//	if err != nil {
//	  fmt.Printf("Error: %v\n", err)
//	  return
//	}
//
// Note: Subdevice is not supported yet.
func LoadPCIDeviceMap() error {
	pciDeviceMap = make(types.PCIDeviceMap, 0)

	pciIDsData, err := os.ReadFile("/usr/share/misc/pci.ids")
	if err != nil {
		fmt.Println("Warning: /usr/share/misc/pci.ids not found, using embedded data assuming we are in a development environment.")
		pciIDsData = []byte(pciIDsData)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(pciIDsData)))
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#") {
			continue
		}

		if line == "" {
			continue
		}

		if !strings.HasPrefix(line, "\t") {
			// Vendor
			parts := strings.Split(line, "  ")
			currentVendor := types.PCIDeviceMapVendor{
				ID:   parts[0],
				Name: parts[1],
			}
			pciDeviceMap = append(pciDeviceMap, currentVendor)
		} else if strings.HasPrefix(line, "\t\t") {
			// Subdevice not supported yet
			continue
		} else {
			// Device
			parts := strings.Split(line, "  ")
			if len(parts) != 2 {
				continue
			}

			id := strings.TrimSpace(parts[0])
			name := strings.TrimSpace(parts[1])
			currentDevice := types.PCIDeviceMapDevice{
				ID:   id,
				Name: name,
			}
			pciDeviceMap[len(pciDeviceMap)-1].Devices = append(pciDeviceMap[len(pciDeviceMap)-1].Devices, currentDevice)
		}
	}

	return nil
}
