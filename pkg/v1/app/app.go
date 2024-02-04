package app

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"

	"github.com/vanilla-os/sdk/pkg/v1/app/types"
	"github.com/vanilla-os/sdk/pkg/v1/i18n"
	"github.com/vanilla-os/sdk/pkg/v1/logs"
)

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
//	})
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//	fmt.Printf("App Sign: %s\n", app.Sign)
func NewApp(options types.AppOptions) (*types.App, error) {
	app := &types.App{
		RDNN:      options.RDNN,
		Name:      options.Name,
		Version:   options.Version,
		LocalesFS: options.LocalesFS,
	}
	app.Sign = generateAppSign(app)

	// here we prepare a logger for the application
	logger, err := logs.NewLogger(app)
	if err != nil {
		return nil, err // logger is mandatory for each application
	}
	app.Logger = logger

	// here we prepare a localizer for the application
	localizer, err := i18n.NewLocalizer(app, options.DefaultLocale)
	if err == nil {
		app.LC = *localizer
	} // something went wrong, perhaps the FS is not provided

	return app, nil
}

// generateAppSign generates a unique signature for the application
// based on the RDNN, name and version. The signature is used to
// identify the application.
func generateAppSign(app *types.App) string {
	sign := fmt.Sprintf("%s-%s-%s", app.RDNN, app.Name, app.Version)

	h := sha1.New()
	h.Write([]byte(sign))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}
