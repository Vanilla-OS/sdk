package types

// NetworkInterfaceInfo represents information about a network interface
type NetworkInterfaceInfo struct {
	// Name is the name of the network interface
	Name string `json:"name"`

	// HardwareAddr is the hardware address of the network interface
	HardwareAddr string `json:"hardware_address"`

	// IPAddresses is a list of IP addresses associated with the
	// network interface
	IPAddresses []string `json:"ip_addresses"`
}
