package media

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/vanilla-os/sdk/pkg/v1/media/types"
)

// runCommand executes a command and returns its standard output, or an error.
func runCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("command '%s %s' failed: %w (stderr: %s)", name, strings.Join(args, " "), err, stderr.String())
	}
	return stdout.String(), nil
}

// GetAudioDevices retrieves a list of audio devices by parsing 'wpctl status'.
// Note: This parsing logic is specific to the expected output format of wpctl
// and might break if the format changes in future versions.
//
// Example:
//
//	import mediaTypes "github.com/vanilla-os/sdk/pkg/v1/media/types"
//
//	devices, err := media.GetAudioDevices()
//	if err != nil {
//		log.Fatalf("Error retrieving audio devices: %v", err)
//	}
//
//	for _, device := range devices {
//		fmt.Printf("Device ID: %s, Name: %s, Type: %s, Default: %t\n", device.ID, device.Name, device.Type, device.IsDefault)
//	}
func GetAudioDevices() ([]types.MediaDevice, error) {
	output, err := runCommand("wpctl", "status")
	if err != nil {
		return nil, fmt.Errorf("failed to run 'wpctl status': %w", err)
	}

	devices := make([]types.MediaDevice, 0)
	lines := strings.Split(output, "\n")

	var currentType types.DeviceType = types.Unknown
	// Regex to capture device ID, name, and default status indicator (*) from
	// the wpctl status output.
	// It parses lines under "Sinks:" or "Sources:" sections.
	// Example lines:
	//   │  * 55. Alder Lake PCH-P High Definition Audio Controller Stereo [vol: 1.08]
	//   │    58. Alder Lake PCH-P High Definition Audio Controller Stereo [vol: 1.00]
	// The regex expects: Tree symbol, optional '*', numeric ID, dot, name, optional volume block.
	// It captures the default indicator ('*'), the numeric ID, and the device name.
	deviceRegex := regexp.MustCompile(`^\s*([└├│])(?:─|\s)*\s*(\*)?\s*(\d+)\.\s+(.*?)(\s*\[vol:.*?\])?\s*$`)
	inAudioSection := false

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Detect the start of the main "Audio" section
		if strings.HasPrefix(trimmedLine, "Audio") && !strings.Contains(trimmedLine, "/") {
			inAudioSection = true
			currentType = types.Unknown
			continue
		}
		if (strings.HasPrefix(trimmedLine, "Video") || strings.HasPrefix(trimmedLine, "Settings")) && !strings.Contains(trimmedLine, "/") {
			inAudioSection = false
			continue
		}

		// We only want to parse lines inside the main "Audio" section
		if inAudioSection {
			// Detect subsection headers (Sinks, Sources) and skip them
			if strings.Contains(trimmedLine, "Sinks:") {
				currentType = types.AudioOutput
				continue
			} else if strings.Contains(trimmedLine, "Sources:") {
				currentType = types.AudioInput
				continue
			} else if strings.Contains(trimmedLine, "Devices:") || strings.Contains(trimmedLine, "Filters:") || strings.Contains(trimmedLine, "Streams:") {
				currentType = types.Unknown
				continue
			}

			if currentType == types.AudioInput || currentType == types.AudioOutput {
				matches := deviceRegex.FindStringSubmatch(line)
				if len(matches) >= 5 {
					dev := types.MediaDevice{
						ID:          matches[3],
						Name:        strings.TrimSpace(matches[4]),
						Type:        currentType,
						IsDefault:   matches[2] == "*",
						Description: strings.TrimSpace(matches[4]),
					}
					devices = append(devices, dev)
				}
			}
		}
	}

	if len(devices) == 0 {
		// If wpctl ran but no devices were parsed, the format might have changed.
		// Or maybe there really are no devices. Return empty list without error for now.
		fmt.Println("Warning: No audio devices parsed from wpctl status.")
	}

	return devices, nil
}

// GetMasterVolume retrieves the current master volume level and mute status
// for the default audio sink.
//
// Example:
//
//	import mediaTypes "github.com/vanilla-os/sdk/pkg/v1/media/types"
//
//	volumeInfo, err := media.GetMasterVolume()
//	if err != nil {
//		log.Fatalf("Error retrieving master volume: %v", err)
//	}
//	fmt.Printf("Volume Level: %d%%, Muted: %t\n", volumeInfo.LevelPercent, volumeInfo.IsMuted)
func GetMasterVolume() (*types.VolumeInfo, error) {
	// @DEFAULT_AUDIO_SINK@ is a special specifier for the default audio sink
	output, err := runCommand("wpctl", "get-volume", "@DEFAULT_AUDIO_SINK@")
	if err != nil {
		// Check if the error indicates the sink wasn't found
		if strings.Contains(err.Error(), "No node found matching specifier") {
			return nil, fmt.Errorf("default audio sink not found")
		}
		return nil, fmt.Errorf("failed to run 'wpctl get-volume': %w", err)
	}

	volumeStr := strings.TrimSpace(output)
	isMuted := strings.Contains(volumeStr, "[MUTED]")

	// Extract volume value
	re := regexp.MustCompile(`Volume:\s*([0-9.]+)`)
	matches := re.FindStringSubmatch(volumeStr)
	if len(matches) != 2 {
		return nil, fmt.Errorf("could not parse volume from wpctl output: '%s'", volumeStr)
	}

	volumeLevel, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse volume value '%s': %w", matches[1], err)
	}

	// Convert 0.0-1.0+ scale to a valid percentage
	volumePercent := int(volumeLevel * 100)

	return &types.VolumeInfo{
		LevelPercent: volumePercent,
		IsMuted:      isMuted,
	}, nil
}

// SetMasterVolume sets the master volume level for the default audio sink.
//
// Example:
//
//	err := media.SetMasterVolume(75)
//	if err != nil {
//		log.Fatalf("Error setting master volume: %v", err)
//	}
//	fmt.Println("Master volume set to 75%")
func SetMasterVolume(volumePercent int) error {
	if volumePercent < 0 {
		volumePercent = 0
	}
	// PipeWire often supports > 100%, common limit is 150%
	if volumePercent > 150 {
		volumePercent = 150
	}

	volumeArg := fmt.Sprintf("%d%%", volumePercent)
	_, err := runCommand("wpctl", "set-volume", "@DEFAULT_AUDIO_SINK@", volumeArg)
	if err != nil {
		return fmt.Errorf("failed to run 'wpctl set-volume': %w", err)
	}
	return nil
}

// SetMuteStatus sets the mute status for the default audio sink.
// It accepts a boolean value: true to mute, false to unmute.
//
// Example:
//
//	err := media.SetMuteStatus(true)
//	if err != nil {
//		log.Fatalf("Error muting audio: %v", err)
//	}
func SetMuteStatus(muted bool) error {
	muteArg := "0"
	if muted {
		muteArg = "1"
	}

	_, err := runCommand("wpctl", "set-mute", "@DEFAULT_AUDIO_SINK@", muteArg)
	if err != nil {
		return fmt.Errorf("failed to run 'wpctl set-mute': %w", err)
	}
	return nil
}
