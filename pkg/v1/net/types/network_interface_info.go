package types

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

// NetworkInterfaceInfo represents information about a network interface
type NetworkInterfaceInfo struct {
	// Name is the name of the network interface
	Name string `json:"name"`

	// HardwareAddr is the hardware address of the network interface
	HardwareAddr string `json:"hardware_address"`

	// IPAddresses is a list of IP addresses associated with the
	// network interface
	IPAddresses []string `json:"ip_addresses"`

	// Status is the status of the network interface
	Status NetworkInterfaceStatus `json:"status"`

	// Running indicates whether the network interface is running
	Running bool `json:"running"`

	// SupportsBroadcast indicates whether the network interface supports
	// broadcast
	SupportsBroadcast bool `json:"supports_broadcast"`

	// SupportsMulticast indicates whether the network interface supports
	// multicast
	SupportsMulticast bool `json:"supports_multicast"`

	// IsLoopback indicates whether the network interface is a loopback
	// interface
	IsLoopback bool `json:"is_loopback"`

	// IsP2P indicates whether the network interface is a point-to-point
	// interface
	IsP2P bool `json:"is_point_to_point"`
}

// NetworkInterfaceStatus represents the status of a network interface
type NetworkInterfaceStatus string

const (
	// NetworkInterfaceStatusUp indicates that the network interface is up
	NetworkInterfaceStatusUp NetworkInterfaceStatus = "up"

	// NetworkInterfaceStatusDown indicates that the network interface is down
	NetworkInterfaceStatusDown NetworkInterfaceStatus = "down"

	// NetworkInterfaceStatusUnknown indicates that the status of the network
	// interface is unknown
	NetworkInterfaceStatusUnknown NetworkInterfaceStatus = "unknown"
)
