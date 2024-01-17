package tests

import (
	"testing"

	"github.com/vanilla-os/sdk/pkg/v1/hardware"
)

func TestGetPCIDeviceByIDs(t *testing.T) {
	err := hardware.LoadPCIDeviceMap()
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	peripherals, err := hardware.GetPCIDevices()
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	for _, peripheral := range peripherals {
		_, vendorName, err := hardware.GetPCIDeviceByIDs(peripheral.Vendor, peripheral.Device)
		if err != nil {
			t.Logf("Error: Unknown vendor: %s", peripheral.Vendor)
			continue
		}
		t.Logf("VendorName: %s Vendor: %s", vendorName, peripheral.Vendor)
	}
}
