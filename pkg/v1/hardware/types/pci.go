package types

// PCIDeviceMap represents a map of PCIDeviceMapVendor structs
type PCIDeviceMap []PCIDeviceMapVendor

// PCIDeviceMapDevice represents a PCI (Peripheral Component Interconnect)
// device
type PCIDeviceMapDevice struct {
	ID   string
	Name string
}

// PCIDeviceMapVendor represents a PCI (Peripheral Component Interconnect)
// vendor
type PCIDeviceMapVendor struct {
	ID      string
	Name    string
	Devices []PCIDeviceMapDevice
}

// PCIDevice represents a PCI (Peripheral Component Interconnect) device
type PCIDevice struct {
	// ID is the unique identifier for the PCI device
	ID string `json:"id"`

	// Class represents the class code of the PCI device
	Class string `json:"class"`

	// Name is the name or description of the PCI device
	Name string `json:"name"`

	// VendorName is the name of the PCI device vendor
	VendorName string `json:"vendor_name"`

	// Vendor is the PCI vendor identifier
	Vendor string `json:"vendor"`

	// Device is the PCI device identifier
	Device string `json:"device"`

	// SubsystemDevice is the identifier for the PCI device's subsystem
	SubsystemDevice string `json:"subsystem_device"`

	// SubsystemVendor is the identifier for the PCI device's subsystem vendor
	SubsystemVendor string `json:"subsystem_vendor"`

	// Modalias is a string that represents the PCI device in a modalias format
	Modalias string `json:"modalias"`
}
