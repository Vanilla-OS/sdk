package types

// BatteryStats represents battery statistics
type BatteryStats struct {
	Percentage int           `json:"percentage"`
	Status     BatteryStatus `json:"status"`
	Voltage    int           `json:"voltage"`
}

// BatteryStatus represents the status of a battery
type BatteryStatus string

const (
	// BatteryStatusUnknown indicates an unknown battery status
	BatteryStatusUnknown BatteryStatus = "unknown"

	// BatteryStatusCharging indicates the battery is currently charging
	BatteryStatusCharging BatteryStatus = "charging"

	// BatteryStatusDischarging indicates the battery is currently discharging
	BatteryStatusDischarging BatteryStatus = "discharging"

	// BatteryStatusFull indicates the battery is fully charged
	BatteryStatusFull BatteryStatus = "full"
)
