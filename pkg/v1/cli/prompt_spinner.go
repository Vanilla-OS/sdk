package cli

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Spinner prompt implementation using Bubble Tea (spinner).
*/

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SpinnerModel struct {
	program *tea.Program
	quit    chan struct{}
}

type spinnerComponent struct {
	spinner  spinner.Model
	message  string
	quitting bool
}

type quitMsg struct{}

type updateMsg string

func initialSpinnerComponent(message string) spinnerComponent {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return spinnerComponent{
		spinner: s,
		message: message,
	}
}

func (m spinnerComponent) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m spinnerComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case quitMsg:
		m.quitting = true
		return m, tea.Quit
	case updateMsg:
		m.message = string(msg)
		return m, nil
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m spinnerComponent) View() string {
	if m.quitting {
		return ""
	}
	return fmt.Sprintf("\n %s %s\n\n", m.spinner.View(), m.message)
}

// StartSpinner starts a spinner with a message.
// The spinner can be stopped by calling the Stop method on the returned model.
//
// Example:
//
//	spinner := myApp.CLI.StartSpinner("Loading the batmobile...")
//	time.Sleep(3 * time.Second)
//	spinner.Stop()
func (c Command) StartSpinner(message string) *SpinnerModel {
	p := tea.NewProgram(initialSpinnerComponent(message))
	quit := make(chan struct{})

	go func() {
		if _, err := p.Run(); err != nil {
			fmt.Println("Error running spinner:", err)
		}
		close(quit)
	}()

	return &SpinnerModel{
		program: p,
	}
}

// UpdateMessage updates the spinner message dynamically.
//
// Example:
//
//	spinner.UpdateMessage("Loading the batcave...")
func (m *SpinnerModel) UpdateMessage(message string) {
	if m.program != nil {
		m.program.Send(updateMsg(message))
	}
}

func (m *SpinnerModel) Stop() {
	if m.program != nil {
		m.program.Send(quitMsg{})
		// Wait for clear
		time.Sleep(100 * time.Millisecond)
	}
}
