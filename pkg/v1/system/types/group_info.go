package types

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

// GroupInfo represents information about a group
type GroupInfo struct {
	// GID is the group's ID
	GID string `json:"gid"`

	// Name is the group's name
	Name string `json:"name"`
}
