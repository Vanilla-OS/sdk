package main

import (
	"fmt"
	"os"
	"time"

	"github.com/vanilla-os/sdk/pkg/v1/app"
	"github.com/vanilla-os/sdk/pkg/v1/app/types"
	"github.com/vanilla-os/sdk/pkg/v1/cli"
)

var myApp *app.App

type RootCmd struct {
	cli.Base
	Poll PollCmd `cmd:"poll" help:"Ask the user preferred hero"`
	Man  ManCmd  `cmd:"man" help:"Generate man page"`
}

type ManCmd struct {
	cli.Base
}

func (c *ManCmd) Run() error {
	man, err := cli.GenerateManPage(&RootCmd{})
	if err != nil {
		return err
	}
	fmt.Print(man)
	return nil
}

type PollCmd struct {
	cli.Base
}

func (c *PollCmd) Run() error {
	fmt.Printf("Welcome to %s (%s)!\n", myApp.Name, myApp.Version)
	hero, err := myApp.CLI.SelectOption(
		"What is your preferred hero?",
		[]string{"Batman", "Ironman", "Spiderman", "Robin", "None"},
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	switch hero {
	case "Batman":
		myApp.Log.Term.Info().Msg("I am Batman!")
	case "Ironman":
		myApp.Log.Term.Info().Msg("Yeah Ironman is cool!")
	case "Spiderman":
		myApp.Log.Term.Info().Msg("Spiderman is ok.")
	case "Robin":
		myApp.Log.Term.Info().Msg("Nobody likes Robin.")
	case "None":
		// Let's ask if they want to pick an unlisted hero
		confirm, err := myApp.CLI.ConfirmAction(
			"Do you want to pick an unlisted hero?",
			"Y", "n",
			true,
		)
		if err != nil {
			fmt.Println(err)
			return err
		}

		if confirm {
			hero, err := myApp.CLI.PromptText(
				"Enter your preferred hero",
				"Batman",
			)
			if err != nil {
				fmt.Println(err)
				return err
			}

			// Let's simulate a spinner while we process the input
			spinner := myApp.CLI.StartSpinner("Looking for your hero...")
			time.Sleep(3 * time.Second)
			spinner.Stop()
			myApp.Log.Term.Info().Msgf("You picked %s", hero)

			// Showing a progressbar now
			bar := myApp.CLI.StartProgressBar("Preparing your hero...", 100)
			for i := 0; i <= 100; i++ {
				bar.Increment(1)
				time.Sleep(50 * time.Millisecond)
			}

			myApp.Log.Term.Info().Msgf("Here is %s!", hero)
		}
	}

	return nil
}

func main() {
	// Here we create a new Vanilla OS application
	var err error
	myApp, err = app.NewApp(types.AppOptions{
		RDNN:    "com.vanillaos.batpoll",
		Name:    "BatPoll",
		Version: "1.0.0",
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Here we add a CLI to the application
	err = myApp.WithCLI(&RootCmd{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// And finally, run the CLI
	myApp.CLI.Execute()
}
