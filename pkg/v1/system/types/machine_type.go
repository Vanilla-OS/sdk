package types

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
