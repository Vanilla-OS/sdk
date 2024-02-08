package tests

import (
	"testing"

	"github.com/vanilla-os/sdk/pkg/v1/app"
	appTypes "github.com/vanilla-os/sdk/pkg/v1/app/types"
	"github.com/vanilla-os/sdk/pkg/v1/logs"
)

func TestNewLogger(t *testing.T) {
	app, err := app.NewApp(appTypes.AppOptions{
		RDNN:    "com.vanilla-os.batsignal",
		Name:    "BatSignal",
		Version: "1.0.0",
	})
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	logger, err := logs.NewLogger(string(app.Sign))
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	logger.File.Info().Msg("Batman reached the file logger")
	logger.Term.Info().Msg("Batman reached the console logger")

	logger.File.Info().Str("where", "file").Msg("Batman is saving Gotham")
}
