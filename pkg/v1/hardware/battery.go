package hardware

import (
	"fmt"
	"os"
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
	slot, err := getBatterySlot()
	if err != nil {
		// If battery slot is not found, assume it's not a portable device
		return nil, nil
	}

	sysfsBatteryPath := "/sys/class/power_supply/" + slot

	capacityContent, err := readSysFile(sysfsBatteryPath, "charge_full")
	if err != nil {
		return nil, fmt.Errorf("failed to read battery capacity: %v", err)
	}

	capacity, err := strconv.Atoi(strings.TrimSpace(capacityContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse battery capacity: %v", err)
	}

	capacityDesignContent, err := readSysFile(sysfsBatteryPath, "charge_full_design")
	if err != nil {
		return nil, fmt.Errorf("failed to read battery design capacity: %v", err)
	}

	capacityDesign, err := strconv.Atoi(strings.TrimSpace(capacityDesignContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse battery design capacity: %v", err)
	}

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
	case "not charging":
		status = types.BatteryStatusNotCharging
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
		Capacity:       capacity,
		CapacityDesign: capacityDesign,
		Percentage:     percentage,
		Status:         status,
		Voltage:        voltage,
	}

	return batteryStats, nil
}

// GetBatteryHealth calculates the battery health based on the battery
// statistics. It returns the battery health percentage.
//
// Example:
//
//	batteryStats, err := hardware.GetBatteryStats()
//	if err != nil {
//	  fmt.Printf("Error: %v\n", err)
//	  return
//	}
//
//	batteryHealth, err := hardware.GetBatteryHealth(batteryStats)
//	if err != nil {
//	  fmt.Printf("Error: %v\n", err)
//	  return
//	}
//
//	fmt.Printf("Health: %f\n", batteryHealth)
func GetBatteryHealth(batteryStats *types.BatteryStats) (float64, error) {
	health := (float64(batteryStats.Capacity) / float64(batteryStats.CapacityDesign)) * 100
	if health > 100 {
		health = 100
	}

	return health, nil
}

// getBatterySlot returns the available battery slot.
func getBatterySlot() (string, error) {
	const sysfsPowerSupplyPath = "/sys/class/power_supply"

	files, err := os.ReadDir(sysfsPowerSupplyPath)
	if err != nil {
		return "", fmt.Errorf("failed to read power supply directory: %v", err)
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "BAT") {
			return file.Name(), nil
		}
	}

	return "", fmt.Errorf("no battery slot found")
}
