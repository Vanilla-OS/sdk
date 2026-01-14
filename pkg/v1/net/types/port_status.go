package types

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

// PortStatus represents the status of a network port
type PortStatus string

const (
	PortOpen    PortStatus = "open"
	PortClosed  PortStatus = "closed"
	PortUnknown PortStatus = "unknown"
)
