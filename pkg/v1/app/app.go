package app

import (
	"crypto/sha1"
	"embed"
	"encoding/base64"
	"fmt"

	"github.com/vanilla-os/sdk/pkg/v1/app/types"
	"github.com/vanilla-os/sdk/pkg/v1/cli"
	cliTypes "github.com/vanilla-os/sdk/pkg/v1/cli/types"
	"github.com/vanilla-os/sdk/pkg/v1/i18n"
	"github.com/vanilla-os/sdk/pkg/v1/logs"
	logsTypes "github.com/vanilla-os/sdk/pkg/v1/logs/types"
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
	Log logsTypes.Logger

	// LC (Localizer) is the localizer for the application
	LC spreak.Localizer

	// LocalesFS is the file system containing the locales for the application
	LocalesFS embed.FS

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
	app.Log = logger

	// here we prepare a localizer for the application
	localizer, err := i18n.NewLocalizer(options.LocalesFS, app.RDNN, options.DefaultLocale)
	if err == nil {
		app.LC = *localizer
	} // something went wrong, perhaps the FS is not provided

	return &app, nil
}

// WithCLI adds a command line interface to the application
//
// Example:
//
//	app.WithCLI(&cli.CLIOptions{
//		Use: "batsignal",
//		Short: "A simple CLI to call Batman",
//		Long: "A simple CLI to call Batman using the BatSignal",
//	})
func (app *App) WithCLI(options *cliTypes.CLIOptions) {
	app.CLI = cli.NewCLI(options)
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
