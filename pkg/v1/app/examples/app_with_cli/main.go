package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/sdk/pkg/v1/app"
	"github.com/vanilla-os/sdk/pkg/v1/app/types"
	"github.com/vanilla-os/sdk/pkg/v1/cli"
	cliTypes "github.com/vanilla-os/sdk/pkg/v1/cli/types"
)

var myApp *app.App

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

	// We need to add a CLI to our application
	myApp.WithCLI(&cliTypes.CLIOptions{
		Use:   "batpoll",
		Short: "CLI to ask the user preferred hero",
		Long:  "A simple CLI to ask the user about their preferred hero",
	})

	// Let's add a first command to our CLI
	pollCmd := cli.NewCommand(
		"poll",
		"Ask the user preferred hero",
		"A simple poll to ask the user about their preferred hero",
		startPoll,
	)
	myApp.CLI.AddCommand(pollCmd)

	// And finally, run the CLI
	myApp.CLI.Execute()
}

// Asking the user preferred hero
func startPoll(cmd *cobra.Command, args []string) error {
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
		myApp.Log.Term.Info().Msg("You don't like any hero?")
	}

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
		myApp.Log.Term.Info().Msgf("You picked %s", hero)
	}

	return nil
}
