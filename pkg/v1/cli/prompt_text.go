package cli

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

// PromptText prompts the user to input a text, it supports customizing the
// prompt and the placeholder.
//
// Example:
//
//	response, err := myApp.CLI.PromptText(
//		"What is your name?",
//		"Bruce Wayne",
//	)
//	if err != nil {
//		fmt.Println(err)
//		return err
//	}
//	fmt.Printf("Hello %s!\n", response)
func (c *Command) PromptText(prompt, placeholder string) (string, error) {
	var response string
	inputPrompt := &survey.Input{
		Message: prompt,
		Default: placeholder,
	}
	err := survey.AskOne(inputPrompt, &response)
	if err != nil {
		return "", fmt.Errorf("input prompt failed: %v", err)
	}
	return response, nil
}
