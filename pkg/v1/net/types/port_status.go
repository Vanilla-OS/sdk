package types

// PortStatus represents the status of a network port
type PortStatus string

const (
	PortOpen    PortStatus = "open"
	PortClosed  PortStatus = "closed"
	PortUnknown PortStatus = "unknown"
)
