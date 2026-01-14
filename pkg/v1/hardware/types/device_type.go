package types

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

// DeviceType represents a system device type
type DeviceType string

const (
	// InputDeviceType represents an input device
	InputDeviceType DeviceType = "input"

	// PCIDeviceType represents a PCI device
	PCIDeviceType DeviceType = "pci"
)
