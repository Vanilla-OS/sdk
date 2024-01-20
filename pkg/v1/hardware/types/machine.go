package types

// ChassisType represents the standardized chassis type.
type ChassisType string

const (
	OtherChassis          ChassisType = "Other"
	UnknownChassis        ChassisType = "Unknown"
	DesktopChassis        ChassisType = "Desktop"
	LowProfileDesktop     ChassisType = "Low Profile Desktop"
	PizzaBoxChassis       ChassisType = "Pizza Box"
	MiniTowerChassis      ChassisType = "Mini Tower"
	TowerChassis          ChassisType = "Tower"
	PortableChassis       ChassisType = "Portable"
	LaptopChassis         ChassisType = "Laptop"
	NotebookChassis       ChassisType = "Notebook"
	HandHeldChassis       ChassisType = "Hand Held"
	DockingStation        ChassisType = "Docking Station"
	AllInOneChassis       ChassisType = "All in One"
	SubNotebookChassis    ChassisType = "Sub Notebook"
	SpaceSavingChassis    ChassisType = "Space-saving"
	LunchBoxChassis       ChassisType = "Lunch Box"
	MainServerChassis     ChassisType = "Main Server Chassis"
	ExpansionChassis      ChassisType = "Expansion Chassis"
	SubChassis            ChassisType = "SubChassis"
	BusExpansionChassis   ChassisType = "Bus Expansion Chassis"
	PeripheralChassis     ChassisType = "Peripheral Chassis"
	RAIDChassis           ChassisType = "RAID Chassis"
	RackMountChassis      ChassisType = "Rack Mount Chassis"
	SealedCasePC          ChassisType = "Sealed-case PC"
	MultiSystemChassis    ChassisType = "Multi-system chassis"
	CompactPCIChassis     ChassisType = "Compact PCI"
	AdvancedTCAChassis    ChassisType = "Advanced TCA"
	BladeChassis          ChassisType = "Blade"
	BladeEnclosureChassis ChassisType = "Blade Enclosure"
	TabletChassis         ChassisType = "Tablet"
	ConvertibleChassis    ChassisType = "Convertible"
	DetachableChassis     ChassisType = "Detachable"
	IoTGatewayChassis     ChassisType = "IoT Gateway"
	EmbeddedPCChassis     ChassisType = "Embedded PC"
	MiniPCChassis         ChassisType = "Mini PC"
	StickPCChassis        ChassisType = "Stick PC"
)

// MachineInfo contains various information about the current machine
type MachineInfo struct {
	ProductName  string      `json:"productName"`
	Manufacturer string      `json:"manufacturer"`
	Version      string      `json:"version"`
	Chassis      ChassisInfo `json:"chassis"`
	Bios         BiosInfo    `json:"bios"`
	Board        BoardInfo   `json:"board"`
}

// ChassisInfo contains various information about the machine chassis
type ChassisInfo struct {
	ID           int         `json:"id"`
	Type         ChassisType `json:"type"`
	Manufacturer string      `json:"manufacturer"`
	Version      string      `json:"version"`
}

// BiosInfo contains various information about the machine bios
type BiosInfo struct {
	Vendor  string `json:"vendor"`
	Version string `json:"version"`
	Release string `json:"release"`
}

// BoardInfo contains various information about the machine board
type BoardInfo struct {
	ProductName  string `json:"product"`
	Manufacturer string `json:"manufacturer"`
	Version      string `json:"version"`
}
