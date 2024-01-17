package tests

import (
	"testing"

	"github.com/vanilla-os/sdk/pkg/v1/hardware"
)

func TestGetPCIDevices(t *testing.T) {
	err := hardware.LoadPCIDeviceMap()
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	devices, err := hardware.GetPCIDevices()
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	for _, device := range devices {
		t.Logf("ID: %s", device.ID)
		t.Logf("Name: %s", device.Name)
		t.Logf("Class: %s", device.Class)
		t.Logf("Vendor: %s", device.Vendor)
		t.Logf("Device: %s", device.Device)
		t.Logf("SubsystemDevice: %s", device.SubsystemDevice)
		t.Logf("SubsystemVendor: %s", device.SubsystemVendor)
		t.Logf("Modalias: %s", device.Modalias)
	}
}
