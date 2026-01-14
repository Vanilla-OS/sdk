package tests

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

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
