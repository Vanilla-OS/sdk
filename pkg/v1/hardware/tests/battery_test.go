package tests

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

import (
	"fmt"
	"strings"
	"testing"

	"github.com/vanilla-os/sdk/pkg/v1/hardware"
)

func TestGetBatteryStats(t *testing.T) {
	batteryStats, err := hardware.GetBatteryStats()
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			t.Skipf("Skipping test due to missing battery files: %v", err)
			return
		}
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

func TestGetBatteryHealth(t *testing.T) {
	batteryStats, err := hardware.GetBatteryStats()
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			t.Skipf("Skipping test due to missing battery files: %v", err)
			return
		}
		t.Errorf("Error: %v", err)
		return
	}

	if batteryStats == nil {
		t.Skip("No battery included")
		return
	}

	batteryHealth, err := hardware.GetBatteryHealth(batteryStats)
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	t.Logf("Health: %s", fmt.Sprintf("%f", batteryHealth))
}
