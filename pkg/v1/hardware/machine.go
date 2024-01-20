package hardware

import (
	"strconv"
	"strings"

	"github.com/vanilla-os/sdk/pkg/v1/hardware/types"
)

func GetMachineInfo() (types.MachineInfo, error) {
	var info types.MachineInfo
	var chassisInfo types.ChassisInfo

	const sysfsDmiPath = "/sys/class/dmi/id"

	productName, err := readSysFile(sysfsDmiPath, "product_name")
	if err != nil {
		info.ProductName = ""
	} else {
		info.ProductName = strings.TrimSpace(productName)
	}

	manufacturer, err := readSysFile(sysfsDmiPath, "sys_vendor")
	if err != nil {
		info.Manufacturer = "Unknown"
	} else {
		info.Manufacturer = strings.TrimSpace(manufacturer)
	}

	version, err := readSysFile(sysfsDmiPath, "product_version")
	if err != nil {
		info.Version = ""
	} else {
		info.Version = strings.TrimSpace(version)
	}

	// Chassis information
	chassisType, err := readSysFile(sysfsDmiPath, "chassis_type")
	if err != nil {
		chassisInfo.ID = -1
		chassisInfo.Type = types.UnknownChassis
	} else {
		chassisInfo.ID, err = strconv.Atoi(strings.TrimSpace(chassisType))
		if err != nil {
			return info, err
		}
		chassisInfo.Type = MapChassisType(chassisInfo.ID)
	}

	chassisVendor, err := readSysFile(sysfsDmiPath, "chassis_vendor")
	if err != nil {
		chassisInfo.Manufacturer = "Unknown"
	} else {
		chassisInfo.Manufacturer = strings.TrimSpace(chassisVendor)
	}

	chassisVersion, err := readSysFile(sysfsDmiPath, "chassis_version")
	if err != nil {
		chassisInfo.Version = "Unknown"
	} else {
		chassisInfo.Version = strings.TrimSpace(chassisVersion)
	}

	info.Chassis = chassisInfo

	// BIOS information
	biosVendor, err := readSysFile(sysfsDmiPath, "bios_vendor")
	if err != nil {
		info.Bios.Vendor = "Unknown"
	} else {
		info.Bios.Vendor = strings.TrimSpace(biosVendor)
	}

	biosVersion, err := readSysFile(sysfsDmiPath, "bios_version")
	if err != nil {
		info.Bios.Version = "Unknown"
	} else {
		info.Bios.Version = strings.TrimSpace(biosVersion)
	}

	biosRelease, err := readSysFile(sysfsDmiPath, "bios_release")
	if err != nil {
		info.Bios.Release = "Unknown"
	} else {
		info.Bios.Release = strings.TrimSpace(biosRelease)
	}

	// Board information
	boardName, err := readSysFile(sysfsDmiPath, "board_name")
	if err != nil {
		info.Board.ProductName = "Unknown"
	} else {
		info.Board.ProductName = strings.TrimSpace(boardName)
	}

	boardVendor, err := readSysFile(sysfsDmiPath, "board_vendor")
	if err != nil {
		info.Board.Manufacturer = "Unknown"
	} else {
		info.Board.Manufacturer = strings.TrimSpace(boardVendor)
	}

	boardVersion, err := readSysFile(sysfsDmiPath, "board_version")
	if err != nil {
		info.Board.Version = "Unknown"
	} else {
		info.Board.Version = strings.TrimSpace(boardVersion)
	}

	return info, nil
}

// MapChassisType maps the chassis type to a standardized representation.
// Refer to https://www.dmtf.org/sites/default/files/standards/documents/DSP0134_3.7.0.pdf
// section 7.4.1 (System Enclosure or Chassis Types) and
// https://superuser.com/a/1107191 for more information.
func MapChassisType(chassisTypeID int) types.ChassisType {
	switch chassisTypeID {
	case 1:
		return types.OtherChassis
	case 2:
		return types.UnknownChassis
	case 3:
		return types.DesktopChassis
	case 4:
		return types.LowProfileDesktop
	case 5:
		return types.PizzaBoxChassis
	case 6:
		return types.MiniTowerChassis
	case 7:
		return types.TowerChassis
	case 8:
		return types.PortableChassis
	case 9:
		return types.LaptopChassis
	case 10:
		return types.NotebookChassis
	case 11:
		return types.HandHeldChassis
	case 12:
		return types.DockingStation
	case 13:
		return types.AllInOneChassis
	case 14:
		return types.SubNotebookChassis
	case 15:
		return types.SpaceSavingChassis
	case 16:
		return types.LunchBoxChassis
	case 17:
		return types.MainServerChassis
	case 18:
		return types.ExpansionChassis
	case 19:
		return types.SubChassis
	case 20:
		return types.BusExpansionChassis
	case 21:
		return types.PeripheralChassis
	case 22:
		return types.RAIDChassis
	case 23:
		return types.RackMountChassis
	case 24:
		return types.SealedCasePC
	case 25:
		return types.MultiSystemChassis
	case 26:
		return types.CompactPCIChassis
	case 27:
		return types.AdvancedTCAChassis
	case 28:
		return types.BladeChassis
	case 29:
		return types.BladeEnclosureChassis
	case 30:
		return types.TabletChassis
	case 31:
		return types.ConvertibleChassis
	case 32:
		return types.DetachableChassis
	case 33:
		return types.IoTGatewayChassis
	case 34:
		return types.EmbeddedPCChassis
	case 35:
		return types.MiniPCChassis
	case 36:
		return types.StickPCChassis
	default:
		return types.UnknownChassis
	}
}
