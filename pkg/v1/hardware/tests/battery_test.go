package tests

import (
	"testing"

	"github.com/vanilla-os/sdk/pkg/v1/hardware"
)

func TestGetBatteryStats(t *testing.T) {
	batteryStats, err := hardware.GetBatteryStats()
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	if batteryStats == nil {
		t.Skip("No battery included")
		return
	}

	t.Logf("Capacity: %d", batteryStats.Capacity)
	t.Logf("Percentage: %d", batteryStats.Percentage)
	t.Logf("Status: %s", batteryStats.Status)
	t.Logf("Voltage: %d", batteryStats.Voltage)
}
