package types

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

// VolumeInfo represents volume information, typically for the default sink/source.
type VolumeInfo struct {
	// Volume level (0-100, potentially > 100, max 150 defined by this SDK)
	LevelPercent int `json:"level_percent"`

	// Mute status
	IsMuted bool `json:"is_muted"`
}
