package types

// VolumeInfo represents volume information, typically for the default sink/source.
type VolumeInfo struct {
	// Volume level (0-100, potentially > 100, max 150 defined by this SDK)
	LevelPercent int `json:"level_percent"`

	// Mute status
	IsMuted bool `json:"is_muted"`
}
