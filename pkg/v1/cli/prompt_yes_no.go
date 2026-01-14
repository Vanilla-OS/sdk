package cli

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Confirmation prompt implementation using Bubble Tea.
*/

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type confirmModel struct {
	prompt    string
	yesText   string
	noText    string
	choice    bool // true = yes, false = no
	err       error
	submitted bool
}

func initialConfirmModel(prompt, yesText, noText string, defaultChoice bool) confirmModel {
	return confirmModel{
		prompt:  prompt,
		yesText: yesText,
		noText:  noText,
		choice:  defaultChoice,
	}
}

func (m confirmModel) Init() tea.Cmd {
	return nil
}

func (m confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.err = fmt.Errorf("interrupted")
			return m, tea.Quit
		case "y", "Y":
			m.choice = true
		case "n", "N":
			m.choice = false
		case "left", "right", "h", "l", "tab":
			m.choice = !m.choice
		case "enter":
			m.submitted = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m confirmModel) View() string {
	var s strings.Builder
	s.WriteString("\n" + lipgloss.NewStyle().Bold(true).Render(m.prompt) + "\n\n")

	yStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginRight(2)
	nStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginRight(2)

	if m.choice {
		yStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).MarginRight(2)
	} else {
		nStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).MarginRight(2)
	}

	s.WriteString(yStyle.Render(m.yesText))
	s.WriteString(nStyle.Render(m.noText))
	s.WriteString("\n")

	return s.String()
}

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
	// If custom text provided, use it, assuming 'y' maps to yesText and 'n' to noText visual
	p := tea.NewProgram(initialConfirmModel(prompt, yesText, noText, defaultChoice))
	m, err := p.Run()
	if err != nil {
		return false, err
	}

	if m, ok := m.(confirmModel); ok {
		if m.err != nil {
			return false, m.err
		}
		if !m.submitted {
			return defaultChoice, nil
		}
		return m.choice, nil
	}

	return false, fmt.Errorf("could not retrieve confirmation")
}
