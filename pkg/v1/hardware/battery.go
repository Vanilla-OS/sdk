package hardware

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/vanilla-os/sdk/pkg/v1/hardware/types"
)

// GetBatteryStats retrieves battery statistics using sysfs. If the battery
// capacity information is not available, it returns nil assuming it's not a
// portable device.
//
// Example:
//
//	batteryStats, err := hardware.GetBatteryStats()
//	if err != nil {
//	  fmt.Printf("Error: %v\n", err)
//	  return
//	}
//
//	fmt.Printf("Percentage: %d\n", batteryStats.Percentage)
//	fmt.Printf("Status: %s\n", batteryStats.Status)
//	fmt.Printf("Voltage: %d\n", batteryStats.Voltage)
func GetBatteryStats() (*types.BatteryStats, error) {
	const sysfsBatteryPath = "/sys/class/power_supply/BAT0"

	percentageContent, err := readSysFile(sysfsBatteryPath, "capacity")
	if err != nil {
		// If battery capacity information is not available, assume it's not
		// a portable device
		return nil, nil
	}

	percentage, err := strconv.Atoi(strings.TrimSpace(percentageContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse battery percentage: %v", err)
	}

	statusContent, err := readSysFile(sysfsBatteryPath, "status")
	if err != nil {
		return nil, fmt.Errorf("failed to read battery status: %v", err)
	}

	var status types.BatteryStatus
	switch strings.ToLower(strings.TrimSpace(statusContent)) {
	case "charging":
		status = types.BatteryStatusCharging
	case "discharging":
		status = types.BatteryStatusDischarging
	case "full":
		status = types.BatteryStatusFull
	default:
		status = types.BatteryStatusUnknown
	}

	voltageContent, err := readSysFile(sysfsBatteryPath, "voltage_now")
	if err != nil {
		return nil, fmt.Errorf("failed to read battery voltage: %v", err)
	}

	voltage, err := strconv.Atoi(strings.TrimSpace(voltageContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse battery voltage: %v", err)
	}

	batteryStats := &types.BatteryStats{
		Percentage: percentage,
		Status:     status,
		Voltage:    voltage,
	}

	return batteryStats, nil
}
