package app

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"

	"github.com/vanilla-os/sdk/pkg/v1/app/types"
	"github.com/vanilla-os/sdk/pkg/v1/logs"
)

// NewApp creates a new Vanilla OS application, which can be used to
// interact with the system. The application is created with the
// default configuration if no options are provided.
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
		RDNN:    options.RDNN,
		Name:    options.Name,
		Version: options.Version,
	}
	app.Sign = generateAppSign(app)

	logger, err := logs.NewLogger(app)
	if err != nil {
		return nil, err
	}
	app.Logger = logger

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
