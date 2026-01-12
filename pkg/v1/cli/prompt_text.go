package cli

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2024
	Description: Text prompt implementation using Bubble Tea (textinput).
*/

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type textInputModel struct {
	textInput textinput.Model
	err       error
	prompt    string
}

func initialTextInputModel(prompt, placeholder string) textInputModel {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40
	ti.Prompt = "➜ "

	return textInputModel{
		textInput: ti,
		prompt:    prompt,
	}
}

func (m textInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m textInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.err = fmt.Errorf("interrupted")
			return m, tea.Quit
		case tea.KeyEnter:
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.textInput.Width = msg.Width
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m textInputModel) View() string {
	var style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return fmt.Sprintf(
		"%s\n%s\n",
		style.Bold(true).Render(m.prompt),
		m.textInput.View(),
	)
}

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
	p := tea.NewProgram(initialTextInputModel(prompt, placeholder))
	m, err := p.Run()
	if err != nil {
		return "", err
	}

	if m, ok := m.(textInputModel); ok {
		if m.err != nil {
			return "", m.err
		}
		if m.textInput.Value() == "" {
			return placeholder, nil
		}
		return m.textInput.Value(), nil
	}

	return "", fmt.Errorf("could not retrieve manual input")
}
