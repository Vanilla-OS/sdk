package cli

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type confirmModel struct {
	prompt        string
	yesText       string
	noText        string
	defaultChoice bool
	choice        *bool
}

func (m *confirmModel) Init() tea.Cmd {
	return nil
}

func (m *confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "y":
			trueChoice := true
			m.choice = &trueChoice
			return m, tea.Quit
		case "n":
			falseChoice := false
			m.choice = &falseChoice
			return m, tea.Quit
		case "enter":
			// We want to use the default choice if the user presses enter
			// without making a choice
			m.choice = &m.defaultChoice
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *confirmModel) View() string {
	yesText, noText := m.yesText, m.noText
	if m.defaultChoice {
		yesText = strings.ToUpper(m.yesText)
		noText = strings.ToLower(m.noText)
	} else {
		yesText = strings.ToLower(m.yesText)
		noText = strings.ToUpper(m.noText)
	}
	return fmt.Sprintf("%s (%s/%s): ", m.prompt, yesText, noText)
}

// ConfirmAction prompts the user to confirm an action with a customizable
// Yes/No prompt. To define the default choice, set defaultChoice accordingly,
// if true, the default choice will be the string in yesText, otherwise the
// one in noText.
//
// Example:
//
//	confirmed, err := cli.ConfirmAction(
//		"Do you want to continue?",
//		"Yes, continue",
//		"No, cancel",
//		true,
//	)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	if confirmed {
//		fmt.Println("Continuing...")
//	} else {
//		fmt.Println("Cancelled")
//	}
func (c *Command) ConfirmAction(prompt, yesText, noText string, defaultChoice bool) (bool, error) {
	model := confirmModel{
		prompt:        prompt,
		yesText:       yesText,
		noText:        noText,
		defaultChoice: defaultChoice,
	}

	p := tea.NewProgram(&model)
	if _, err := p.Run(); err != nil {
		return false, fmt.Errorf("failed to run confirmation program: %v", err)
	}

	if model.choice == nil {
		return false, fmt.Errorf("no choice made")
	}
	return *model.choice, nil
}
