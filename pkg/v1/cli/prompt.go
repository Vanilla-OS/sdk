package cli

import "github.com/manifoldco/promptui"

// SelectOption prompts the user to select an option from a list of options.
func (c *Command) SelectOption(label string, options []string) (string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: options,
	}

	_, result, err := prompt.Run()

	return result, err
}

// InputText prompts the user to input text.
func (c *Command) InputText(label string) (string, error) {
	prompt := promptui.Prompt{
		Label: label,
	}

	result, err := prompt.Run()

	return result, err
}

// ConfirmAction prompts the user to confirm an action.
func (c *Command) ConfirmAction(label string) (bool, error) {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}

	result, err := prompt.Run()

	return result == "y", err
}
