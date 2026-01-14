package tests

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

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

	for i, device := range devices {
		t.Logf("Name: %s", device.Name)
		t.Logf("Product: %s", device.Product)

		if i == 5 {
			return
		}
	}
}

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

	for i, peripheral := range peripherals {
		_, vendorName, err := hardware.GetPCIDeviceByIDs(peripheral.Vendor, peripheral.Device)
		if err != nil {
			t.Logf("Error: Unknown vendor: %s", peripheral.Vendor)
			continue
		}
		t.Logf("VendorName: %s Vendor: %s", vendorName, peripheral.Vendor)

		if i == 5 {
			return
		}
	}
}

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

	for i, device := range devices {
		t.Logf("ID: %s", device.ID)
		t.Logf("Name: %s", device.Name)
		t.Logf("Class: %s", device.Class)
		t.Logf("Vendor: %s", device.Vendor)
		t.Logf("Device: %s", device.Device)
		t.Logf("SubsystemDevice: %s", device.SubsystemDevice)
		t.Logf("SubsystemVendor: %s", device.SubsystemVendor)
		t.Logf("Modalias: %s", device.Modalias)

		if i == 5 {
			return
		}
	}
}

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

	for i, peripheral := range peripherals {
		t.Logf("ID: %s", peripheral.ID)
		t.Logf("Name: %s", peripheral.Name)
		t.Log("Type:", peripheral.Type)

		if i == 5 {
			return
		}
	}

	if len(peripherals) == 0 {
		t.Log("No peripherals found")
	}
}
