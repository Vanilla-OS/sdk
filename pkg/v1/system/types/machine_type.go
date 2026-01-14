package types

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

// MachineType is the representation of the type of machine
type MachineType string

const (
	// VM is a virtual machine
	VM MachineType = "vm"

	// BareMetal is a bare metal machine, like a workstation or server
	BareMetal MachineType = "baremetal"

	// Container is a pretty much self-explanatory
	Container MachineType = "container"
)
