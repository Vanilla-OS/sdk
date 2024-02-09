package cli

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

// ConfirmAction prompts the user to confirm an action, it supports customizing
// the prompt and the text for the "yes" and "no" options. If the user does not
// provide an answer, the default choice is used.
//
// Example:
//
//	confirm, err := myApp.CLI.ConfirmAction(
//		"Do you like Batman?",
//		"Yes", "No",
//		true,
//	)
//	if err != nil {
//		fmt.Println(err)
//		return err
//	}
//	if confirm {
//		fmt.Println("Everybody likes Batman!")
//	} else {
//		fmt.Println("You don't like Batman...")
//	}
func (c *Command) ConfirmAction(prompt, yesText, noText string, defaultChoice bool) (bool, error) {
	var confirm bool
	confirmationPrompt := &survey.Confirm{
		Message: prompt,
		Default: defaultChoice,
		Help:    fmt.Sprintf("Yes: %s, No: %s", yesText, noText),
	}
	err := survey.AskOne(confirmationPrompt, &confirm)
	if err != nil {
		return false, fmt.Errorf("confirmation failed: %v", err)
	}
	return confirm, nil
}
