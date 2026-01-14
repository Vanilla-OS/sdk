package types

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

// Timezone represents a supported timezone
type Timezone struct {
	// Name is the name of the timezone
	Name string `json:"name"`

	// Location is the location of the timezone file
	Location string `json:"location"`
}
