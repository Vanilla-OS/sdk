package types

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

// Peripheral represents a system peripheral.
type Peripheral struct {
	ID   string     `json:"id"`
	Name string     `json:"name"`
	Type DeviceType `json:"type"`
}
