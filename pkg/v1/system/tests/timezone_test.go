package tests

import (
	"testing"

	"github.com/vanilla-os/sdk/pkg/v1/system"
)

func TestGetSupportedTimezones(t *testing.T) {
	timezones, err := system.GetSupportedTimezones()
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	for i, timezone := range timezones {
		t.Logf("%s\n", timezone.Name)
		t.Logf("%s\n", timezone.Location)

		if i == 5 {
			break
		}
	}
}
