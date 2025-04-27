package tests

import (
	"os/exec"
	"testing"

	"github.com/vanilla-os/sdk/pkg/v1/media"
	media_types "github.com/vanilla-os/sdk/pkg/v1/media/types"
)

// checkWpctl checks if wpctl is available in the PATH.
func checkWpctl(t *testing.T) bool {
	_, err := exec.LookPath("wpctl")
	if err != nil {
		t.Skip("Skipping media tests: 'wpctl' command not found in PATH.")
		return false
	}
	return true
}

func TestGetAudioDevices(t *testing.T) {
	if !checkWpctl(t) {
		return
	}

	devices, err := media.GetAudioDevices()
	if err != nil {
		// Don't fail if wpctl exists but returns an error (e.g., pipewire not running)
		t.Skipf("GetAudioDevices failed (is pipewire running?): %v", err)
		return
	}

	t.Logf("Found %d audio devices via wpctl.", len(devices))
	foundDefaultOutput := false
	foundDefaultInput := false
	for _, device := range devices {
		t.Logf("  ID: %s, Name: %s, Type: %s, Default: %t", device.ID, device.Name, device.Type, device.IsDefault)
		if device.IsDefault {
			if device.Type == media_types.AudioOutput {
				foundDefaultOutput = true
			}
			if device.Type == media_types.AudioInput {
				foundDefaultInput = true
			}
		}
		// Basic validation
		if device.ID == "" {
			t.Errorf("Found device with empty ID: %+v", device)
		}
		if device.Name == "" {
			t.Errorf("Found device with empty Name: %+v", device)
		}
	}

	if !foundDefaultOutput && len(devices) > 0 {
		t.Log("Warning: No default audio output (Sink) marked with '*' found.")
	}
	if !foundDefaultInput && len(devices) > 0 {
		t.Log("Warning: No default audio input (Source) marked with '*' found.")
	}
}

func TestGetSetMasterVolume(t *testing.T) {
	if !checkWpctl(t) {
		return
	}

	initialVolumeInfo, err := media.GetMasterVolume()
	if err != nil {
		t.Skipf("GetMasterVolume failed (is pipewire running / default sink available?): %v", err)
		return
	}
	t.Logf("Initial Master Volume: %d%%, Muted: %t", initialVolumeInfo.LevelPercent, initialVolumeInfo.IsMuted)

	// Test setting volume
	testVolumePercent := 65
	if initialVolumeInfo.LevelPercent == testVolumePercent {
		testVolumePercent = 55
		if initialVolumeInfo.LevelPercent == testVolumePercent {
			testVolumePercent = 75
		}
	}

	t.Logf("Attempting to set volume to %d%%", testVolumePercent)
	err = media.SetMasterVolume(testVolumePercent)
	if err != nil {
		t.Errorf("SetMasterVolume to %d%% failed: %v", testVolumePercent, err)
		media.SetMasterVolume(initialVolumeInfo.LevelPercent)
		return
	}

	// Verify volume was set
	currentVolumeInfo, err := media.GetMasterVolume()
	if err != nil {
		t.Errorf("GetMasterVolume after set failed: %v", err)
	} else {
		if currentVolumeInfo.LevelPercent != testVolumePercent {
			t.Errorf("Volume set verification failed: expected %d%%, got %d%%", testVolumePercent, currentVolumeInfo.LevelPercent)
		} else {
			t.Logf("Volume successfully set and verified to %d%%", currentVolumeInfo.LevelPercent)
		}
	}

	// Restore original volume
	t.Logf("Restoring initial volume to %d%%", initialVolumeInfo.LevelPercent)
	err = media.SetMasterVolume(initialVolumeInfo.LevelPercent)
	if err != nil {
		t.Logf("Warning: Failed to restore initial volume to %d%%: %v", initialVolumeInfo.LevelPercent, err)
	}
}

func TestSetMuteStatus(t *testing.T) {
	if !checkWpctl(t) {
		return
	}

	initialVolumeInfo, err := media.GetMasterVolume()
	if err != nil {
		t.Skipf("GetMasterVolume failed (is pipewire running / default sink available?): %v", err)
		return
	}
	initialMuteStatus := initialVolumeInfo.IsMuted
	t.Logf("Initial Mute Status: %t", initialMuteStatus)

	// Toggle mute status
	targetMuteStatus := !initialMuteStatus
	t.Logf("Attempting to set mute status to %t", targetMuteStatus)
	err = media.SetMuteStatus(targetMuteStatus)
	if err != nil {
		t.Errorf("SetMuteStatus(%t) failed: %v", targetMuteStatus, err)
		// Attempt to restore original status anyway
		media.SetMuteStatus(initialMuteStatus) // Best effort restore
		return
	}

	// Verify mute status
	// time.Sleep(100 * time.Millisecond) // Optional delay
	currentVolumeInfo, err := media.GetMasterVolume()
	if err != nil {
		t.Errorf("GetMasterVolume after set mute failed: %v", err)
	} else {
		if currentVolumeInfo.IsMuted != targetMuteStatus {
			t.Errorf("Mute status verification failed: expected %t, got %t", targetMuteStatus, currentVolumeInfo.IsMuted)
		} else {
			t.Logf("Mute status successfully set and verified to %t", currentVolumeInfo.IsMuted)
		}
	}

	// Restore initial mute status
	t.Logf("Restoring initial mute status to %t", initialMuteStatus)
	err = media.SetMuteStatus(initialMuteStatus)
	if err != nil {
		t.Logf("Warning: Failed to restore initial mute status to %t: %v", initialMuteStatus, err)
	}
}
