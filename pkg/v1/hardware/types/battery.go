package types

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

// BatteryStats represents battery statistics
type BatteryStats struct {
	Capacity       int           `json:"capacity"`
	CapacityDesign int           `json:"capacityDesign"`
	Percentage     int           `json:"percentage"`
	Status         BatteryStatus `json:"status"`
	Voltage        int           `json:"voltage"`
}

// BatteryStatus represents the status of a battery
type BatteryStatus string

const (
	// BatteryStatusUnknown indicates an unknown battery status
	BatteryStatusUnknown BatteryStatus = "unknown"

	// BatteryStatusCharging indicates the battery is currently charging
	BatteryStatusCharging BatteryStatus = "charging"

	// BatteryStatusNotCharging indicates the battery is not charging
	BatteryStatusNotCharging BatteryStatus = "not charging"

	// BatteryStatusDischarging indicates the battery is currently discharging
	BatteryStatusDischarging BatteryStatus = "discharging"

	// BatteryStatusFull indicates the battery is fully charged
	BatteryStatusFull BatteryStatus = "full"
)
