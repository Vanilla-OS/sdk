package tests

import (
	"testing"

	"github.com/vanilla-os/sdk/pkg/v1/hardware"
)

func TestGetPeripheralList(t *testing.T) {
	err := hardware.LoadPCIDeviceMap()
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	peripherals, err := hardware.GetPeripheralList()
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	for _, peripheral := range peripherals {
		t.Logf("ID: %s", peripheral.ID)
		t.Logf("Name: %s", peripheral.Name)
		t.Log("Type:", peripheral.Type)
	}

	if len(peripherals) == 0 {
		t.Log("No peripherals found")
	}
}
