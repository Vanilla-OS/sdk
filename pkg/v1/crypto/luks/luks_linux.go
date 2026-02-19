//go:build linux && cgo

package luks

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

import (
	"errors"
	"path/filepath"

	cryptsetup "github.com/kcolford/go-cryptsetup"
)

// ErrNotLUKS is returned when the target device doesn't look like a LUKS volume.
var ErrNotLUKS = errors.New("not a LUKS device")

// Probe checks if devicePath contains a LUKS1 header.
//
// Example:
//
//	ok, err := luks.Probe("/dev/sda2")
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//	fmt.Printf("Is LUKS: %v\n", ok)
func Probe(devicePath string) (bool, error) {
	d, err := cryptsetup.NewDevice(devicePath)
	if err != nil {
		return false, err
	}
	defer d.Close()

	if err := d.Load(cryptsetup.LuksParams{}); err != nil {
		return false, nil
	}
	return true, nil
}

// FormatLUKS1 formats the given block device as LUKS1.
//
// Example:
//
//	err := luks.FormatLUKS1("/dev/sdb1", []byte("secret"), luks.Params{})
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
func FormatLUKS1(devicePath string, passphrase []byte, params Params) error {
	d, err := cryptsetup.NewDevice(devicePath)
	if err != nil {
		return err
	}
	defer d.Close()

	p := cryptsetup.LuksParams{}
	p.Cipher = params.Cipher
	p.Mode = params.Mode
	p.Hash = params.Hash
	p.VolumeKeySize = params.VolumeKeySize
	p.DataAlignment = params.DataAlignment
	p.DataDevice = params.DataDevice

	return d.Format(passphrase, p)
}

// Open activates a LUKS device mapping with the given mapperName.
//
// Example:
//
//	err := luks.Open("/dev/sdb1", "vos-backup", []byte("secret"))
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
func Open(devicePath, mapperName string, passphrase []byte) error {
	d, err := cryptsetup.NewDevice(devicePath)
	if err != nil {
		return err
	}
	defer d.Close()

	if err := d.Load(cryptsetup.LuksParams{}); err != nil {
		return ErrNotLUKS
	}
	return d.Activate(mapperName, passphrase)
}

// Close deactivates a LUKS device mapping.
//
// Example:
//
//	err := luks.Close("vos-backup")
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
func Close(mapperName string) error {
	d, err := cryptsetup.NewDevice(filepath.Join("/dev/mapper", mapperName))
	if err != nil {
		return err
	}
	defer d.Close()
	return d.Deactivate(mapperName)
}

// AddKey adds a new passphrase to a LUKS device.
//
// Example:
//
//	err := luks.AddKey("/dev/sdb1", []byte("old"), []byte("new"))
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
func AddKey(devicePath string, oldPassphrase, newPassphrase []byte) error {
	d, err := cryptsetup.NewDevice(devicePath)
	if err != nil {
		return err
	}
	defer d.Close()

	if err := d.Load(cryptsetup.LuksParams{}); err != nil {
		return ErrNotLUKS
	}
	return d.AddKey(oldPassphrase, newPassphrase)
}

// DelKey removes a passphrase from a LUKS device.
//
// Example:
//
//	err := luks.DelKey("/dev/sdb1", []byte("secret"))
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
func DelKey(devicePath string, passphrase []byte) error {
	d, err := cryptsetup.NewDevice(devicePath)
	if err != nil {
		return err
	}
	defer d.Close()

	if err := d.Load(cryptsetup.LuksParams{}); err != nil {
		return ErrNotLUKS
	}
	return d.DelKey(passphrase)
}
