//go:build !linux

package luks

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

import "errors"

var ErrUnsupported = errors.New("luks is only supported on Linux")

func Probe(devicePath string) (bool, error) { return false, ErrUnsupported }
func FormatLUKS1(devicePath string, passphrase []byte, params Params) error {
	return ErrUnsupported
}
func Open(devicePath, mapperName string, passphrase []byte) error { return ErrUnsupported }
func Close(mapperName string) error                               { return ErrUnsupported }
func AddKey(devicePath string, oldPassphrase, newPassphrase []byte) error {
	return ErrUnsupported
}
func DelKey(devicePath string, passphrase []byte) error { return ErrUnsupported }
