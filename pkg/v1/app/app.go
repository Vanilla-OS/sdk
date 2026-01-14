package app

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/fs"
	"os"

	"github.com/vanilla-os/sdk/pkg/v1/app/types"
	"github.com/vanilla-os/sdk/pkg/v1/cli"
	"github.com/vanilla-os/sdk/pkg/v1/i18n"
	"github.com/vanilla-os/sdk/pkg/v1/logs"
	"github.com/vorlif/spreak"
)

// App represents a Vanilla OS application
type App struct {
	// Sign is a unique signature for the application
	Sign types.Sign

	// RDNN is the reverse domain name notation of the application
	RDNN string

	// Name is the name of the application
	Name string

	// Version is the version of the application
	Version string

	// Log is the logger for the application
	Log *logs.Logger

	// LC (Localizer) is the localizer for the application
	LC spreak.Localizer

	// LocalesFS is the file system containing the locales for the application
	LocalesFS fs.FS

	// CLI is the command line interface for the application
	CLI *cli.Command
}

// NewApp creates a new Vanilla OS application, which can be used to
// interact with the system. The application is created with the
// default configuration if no options are provided.
//
// Notes:
//
// If the project does not provide a FS for the locales, the
// localizer will not be created. If the localizer fails to be created
// in any way, the application will continue to work without it but
// translation keys will be returned as they are (in English).
//
// The logger, instead, is mandatory for each application. So if the
// logger fails to be created, the application will return an error and
// will not work.
//
// Example:
//
//	app, err := app.NewApp({
//		RDNN: "com.vanilla-os.batsignal",
//		Name: "BatSignal",
//		Version: "1.0.0",
//		LocalesFS: localesFS,
//		DefaultLocale: "en",
//		CLIOptions: &cli.CLIOptions{
//			Use: "batsignal",
//			Short: "A simple CLI to call Batman",
//			Long: "A simple CLI to call Batman using the BatSignal",
//		},
//	})
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//	fmt.Printf("App Sign: %s\n", app.Sign)
func NewApp(options types.AppOptions) (*App, error) {
	app := App{
		RDNN:      options.RDNN,
		Name:      options.Name,
		Version:   options.Version,
		LocalesFS: options.LocalesFS,
	}
	app.Sign = generateAppSign(&app)

	// here we prepare a logger for the application
	logger, err := logs.NewLogger(string(app.Sign))
	if err != nil {
		return &app, err // logger is mandatory for each application
	}
	app.Log = &logger

	// here we prepare a localizer for the application
	locale := os.Getenv("LANGUAGE")
	if locale == "" {
		locale = os.Getenv("LC_ALL")
	}
	if locale == "" {
		locale = os.Getenv("LC_MESSAGES")
	}
	if locale == "" {
		locale = os.Getenv("LANG")
	}
	if locale == "" {
		locale = options.DefaultLocale
	}

	var localizer *spreak.Localizer
	if options.LocalesFS != nil {
		var err error
		localizer, err = i18n.NewLocalizer(options.LocalesFS, app.RDNN, locale)
		if err != nil {
			return &app, err
		}
	}

	if localizer != nil {
		app.LC = *localizer
	}

	app.checkMissingStrings()

	return &app, nil
}

// WithCLI assigns a command created from a struct (declarative model) to the application CLI.
//
// Example:
//
//	root := &RootCmd{
//		Poll: PollCmd{},
//		Man:  ManCmd{},
//	}
//	app.WithCLI(root)
//	app.CLI.Execute()
func (app *App) WithCLI(root any) error {
	cmd, err := cli.NewCommandFromStruct(root)
	if err != nil {
		return err
	}
	app.CLI = cmd
	return nil
}

// generateAppSign generates a unique signature for the application
// based on the RDNN, name and version. The signature is used to
// identify the application.
func generateAppSign(app *App) types.Sign {
	sign := fmt.Sprintf("%s-%s-%s", app.RDNN, app.Name, app.Version)

	h := sha1.New()
	h.Write([]byte(sign))
	return types.Sign(base64.URLEncoding.EncodeToString(h.Sum(nil)))
}
