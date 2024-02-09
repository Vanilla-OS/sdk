package cli

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

// SelectOption prompts the user to select an option from a list of options.
//
// Example:
//
//	selected, err := myApp.CLI.SelectOption(
//		"What is your preferred hero?",
//		[]string{"Batman", "Ironman", "Spiderman", "Robin", "None"},
//	)
//	if err != nil {
//		fmt.Println(err)
//		return err
//	}
//	fmt.Printf("You selected %s!\n", selected)
func (c *Command) SelectOption(prompt string, options []string) (string, error) {
	var selected string
	selectPrompt := &survey.Select{
		Message:  prompt,
		Options:  options,
		PageSize: 5,
	}
	err := survey.AskOne(selectPrompt, &selected)
	if err != nil {
		return "", fmt.Errorf("selection failed: %v", err)
	}
	return selected, nil
}
