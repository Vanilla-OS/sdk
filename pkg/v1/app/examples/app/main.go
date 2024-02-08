package main

import (
	"fmt"
	"os"

	"github.com/vanilla-os/sdk/pkg/v1/app"
	"github.com/vanilla-os/sdk/pkg/v1/app/types"
)

func main() {
	// Here we create a new Vanilla OS application
	myApp, err := app.NewApp(types.AppOptions{
		RDNN:    "com.vanillaos.batsignal",
		Name:    "BatSignal",
		Version: "1.0.0",
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Then log a welcome message
	myApp.Log.Term.Info().Msgf("Welcome to %s (%s)!", myApp.Name, myApp.Version)
	myApp.Log.Term.Info().Msg("You just called Batman!")

	// And finally, log a message to the file logger
	myApp.Log.File.Info().Msg("Batman reached the file logger")
}
