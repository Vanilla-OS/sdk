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

	"github.com/vanilla-os/sdk/pkg/v1/app"
	"github.com/vanilla-os/sdk/pkg/v1/app/types"
)

func TestNewApp(t *testing.T) {
	app, err := app.NewApp(types.AppOptions{
		RDNN:    "com.vanilla-os.batsignal",
		Name:    "BatSignal",
		Version: "1.0.0",
	})
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	if app.Sign == "" {
		t.Errorf("Expected app.Sign to be non-empty, got empty")
	}
	t.Log("App created")
	t.Logf("\tSign: %s\n", app.Sign)
	t.Logf("\tRDNN: %s\n", app.RDNN)
	t.Logf("\tName: %s\n", app.Name)
	t.Logf("\tVersion: %s\n", app.Version)
	t.Logf("\tTesting bundled logger:")
	app.Log.File.Info().Msg("Robin reached the file logger")
	app.Log.Term.Info().Msg("Robin reached the console logger")
}
