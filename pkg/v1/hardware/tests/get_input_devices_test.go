package tests

import (
	"testing"

	"github.com/vanilla-os/sdk/pkg/v1/hardware"
)

func TestGetInputDevices(t *testing.T) {
	err := hardware.LoadPCIDeviceMap()
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	devices, err := hardware.GetInputDevices()
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	for _, device := range devices {
		t.Logf("Name: %s", device.Name)
		t.Logf("Product: %s", device.Product)
	}
}
