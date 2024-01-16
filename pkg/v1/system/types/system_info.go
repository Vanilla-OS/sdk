package types

// SystemInfo is a struct that contains information about the system
type SystemInfo struct {
	// OS is the name of the operating system
	OS string

	// Version is the version of the operating system
	Version string

	// Codename is how the operating system is referred to internally
	Codename string

	// Arch is the architecture of the operating system
	Arch string

	// MachineType is the type of machine that the operating system is running
	// on. This can be baremetal, vm or container.
	MachineType MachineType
}
