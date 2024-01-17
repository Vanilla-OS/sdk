package types

// DeviceType represents a system device type
type DeviceType string

const (
	// InputDeviceType represents an input device
	InputDeviceType DeviceType = "input"

	// PCIDeviceType represents a PCI device
	PCIDeviceType DeviceType = "pci"
)
