package types

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

// DeviceType represents the type of media device (e.g., audio input/output).
type DeviceType string

const (
	AudioInput  DeviceType = "audio_input"
	AudioOutput DeviceType = "audio_output"
	Unknown     DeviceType = "unknown"
)

// MediaDevice represents a system media device, typically parsed from wpctl.
type MediaDevice struct {
	// Numeric ID from wpctl status
	ID string `json:"id"`

	// User-friendly name from wpctl status
	Name string `json:"name"`

	// Often same as Name when parsed from wpctl status
	Description string `json:"description"`

	// Type (Sink/Output or Source/Input)
	Type DeviceType `json:"type"`

	// Is it the default device?
	IsDefault bool `json:"is_default"`
}
