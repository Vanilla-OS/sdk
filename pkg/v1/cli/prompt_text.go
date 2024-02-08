package cli

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type promptTextModel struct {
	textInput textinput.Model
	prompt    string
}

func initialPromptTextModel(prompt, placeholder string) promptTextModel {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()

	return promptTextModel{
		textInput: ti,
		prompt:    prompt,
	}
}

func (m promptTextModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m promptTextModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			return m, tea.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m promptTextModel) View() string {
	return fmt.Sprintf(
		"%s\n\n%s\n",
		m.prompt,
		m.textInput.View(),
	)
}

// PromptText asks the user for text input with an optional placeholder.
//
// Example:
//
//	text, err := cli.PromptText("What is your name?", "Bruce Wayne")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Printf("Hello, %s!\n", text)
func (c *Command) PromptText(prompt, placeholder string) (string, error) {
	model := initialPromptTextModel(prompt, placeholder)
	p := tea.NewProgram(model)

	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to start program: %v", err)
		return "", err
	}

	finalPromptModel, ok := finalModel.(promptTextModel)
	if !ok {
		return "", fmt.Errorf("could not assert final model to promptTextModel")
	}

	return finalPromptModel.textInput.Value(), nil
}
