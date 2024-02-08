package tests

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/sdk/pkg/v1/cli"
	"github.com/vanilla-os/sdk/pkg/v1/cli/types"
)

func TestNewCLI(t *testing.T) {
	root := cli.NewCLI(&types.CLIOptions{
		Use:   "batsignal",
		Short: "A signal for Batman",
		Long:  "A signal for Batman to save Gotham",
	})
	if root == nil {
		t.Errorf("Error: root is nil")
	}

	if root.Use != "batsignal" {
		t.Errorf("Error: expected %s, got %s", "batsignal", root.Use)
	} else {
		t.Logf("Success: expected %s, got %s", "batsignal", root.Use)
	}

	// Here we try to add a command to the root command
	testFn := func(cmd *cobra.Command, args []string) error {
		return nil
	}
	testCommand := cli.NewCommand(
		"test",
		"Test the signal for Batman",
		"Test the signal",
		testFn,
	)
	root.AddCommand(testCommand)
	if len(root.Children()) != 1 {
		t.Errorf("Error: expected 1, got %d", len(root.Children()))
	} else {
		t.Logf("Success: expected 1, got %d", len(root.Children()))
	}

	err := root.Execute()
	if err != nil {
		t.Errorf("Error: %v", err)
	} else {
		t.Logf("Root command executed successfully")
	}
}
