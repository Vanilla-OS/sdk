package cli

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type SpinnerModel struct {
	spinner  spinner.Model
	message  string
	quitting bool
}

func NewSpinnerModel(message string) SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return SpinnerModel{
		spinner: s,
		message: message,
	}
}

func (m SpinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.quitting {
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m SpinnerModel) View() string {
	str := fmt.Sprintf("%s %s\n", m.message, m.spinner.View())
	if m.quitting {
		return str + "\n"
	}
	return str
}

func (m *SpinnerModel) Start() {
	p := tea.NewProgram(m)
	go func() {
		_, err := p.Run()
		if err != nil {
			fmt.Printf("Error running spinner program: %v\n", err)
		}
	}()
}

func (c *Command) StartSpinner(message string) *SpinnerModel {
	model := NewSpinnerModel(message)
	p := tea.NewProgram(model)
	go func() {
		_, err := p.Run()
		if err != nil {
			fmt.Printf("Error running spinner program: %v\n", err)
		}
	}()
	return &model
}

func (m *SpinnerModel) Stop() {
	m.quitting = true
}
