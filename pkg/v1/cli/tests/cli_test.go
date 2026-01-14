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

	"github.com/vanilla-os/sdk/pkg/v1/cli"
)

func TestDeclarativeCLI(t *testing.T) {
	type TestCmd struct {
		cli.Base
	}

	cmd, err := cli.NewCommandFromStruct(&TestCmd{})
	if err != nil {
		t.Fatalf("Failed to create declarative command: %v", err)
	}

	if cmd == nil {
		t.Fatal("Command is nil")
	}

	if cmd.GetRoot() == nil {
		t.Error("Command root is nil")
	}

	// We can't easily test execution output here without capturing stdout/mocking,
	// but we can ensure structural setup is correct.

	// Check man page generation
	man, err := cli.GenerateManPage(&TestCmd{}, func(key string) string {
		return key
	})
	if err != nil {
		t.Errorf("Failed to generate man page: %v", err)
	}
	if man == "" {
		t.Error("Man page is empty")
	}
}
