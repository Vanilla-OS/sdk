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

func TestGetMachineInfo(t *testing.T) {
	info, err := hardware.GetMachineInfo()
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	t.Logf("ProductName: %s", info.ProductName)
	t.Logf("Manufacturer: %s", info.Manufacturer)
	t.Logf("Version: %s", info.Version)
	t.Log("Chassis:")
	t.Logf("\tID: %d", info.Chassis.ID)
	t.Logf("\tType: %s", info.Chassis.Type)
	t.Logf("\tManufacturer: %s", info.Chassis.Manufacturer)
	t.Logf("\tVersion: %s", info.Chassis.Version)
	t.Log("Bios:")
	t.Logf("\tVendor: %s", info.Bios.Vendor)
	t.Logf("\tVersion: %s", info.Bios.Version)
	t.Logf("\tRelease: %s", info.Bios.Release)
	t.Log("Board:")
	t.Logf("\tProductName: %s", info.Board.ProductName)
	t.Logf("\tManufacturer: %s", info.Board.Manufacturer)
	t.Logf("\tVersion: %s", info.Board.Version)
}
