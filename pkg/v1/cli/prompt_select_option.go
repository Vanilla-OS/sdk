package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type itemSelectOption struct {
	text string
}

type modelSelectOption struct {
	items  []itemSelectOption
	cursor int
	chosen int
	done   bool
}

func (m *modelSelectOption) Init() tea.Cmd {
	return nil
}

func (m *modelSelectOption) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "enter":
			m.chosen = m.cursor
			m.done = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *modelSelectOption) View() string {
	if m.done {
		return m.items[m.chosen].text
	}

	var s string
	for i, item := range m.items {
		cursor := "  "
		if i == m.cursor {
			cursor = " ▸"
			s += fmt.Sprintf("%s \x1b[4m%s\x1b[0m\n", cursor, item.text)
		} else {
			s += fmt.Sprintf("%s %s\n", cursor, item.text)
		}
	}
	return s
}

// SelectOption prompts the user to select a single option from a list
// of options. The prompt is the message to display to the user, and
// options is the list of options to choose from. The function returns
// the selected option or an error if no option was selected.
//
// Example:
//
//	hero, err := myApp.CLI.SelectOption(
//		"What is your preferred hero?",
//		[]string{"Batman", "Ironman", "Spiderman", "Robin", "None"},
//	)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Printf("You selected: %s\n", hero)
func (c *Command) SelectOption(prompt string, options []string) (string, error) {
	items := make([]itemSelectOption, len(options))
	for i, opt := range options {
		items[i] = itemSelectOption{text: opt}
	}

	initialModel := modelSelectOption{
		items:  items,
		cursor: 0,
		chosen: -1,
		done:   false,
	}

	fmt.Printf("\x1b[36m? \x1b[37m%s\n", prompt)

	p := tea.NewProgram(&initialModel)
	if _, err := p.Run(); err != nil {
		return "", fmt.Errorf("failed to run selection program: %v", err)
	}

	if initialModel.chosen >= 0 {
		return initialModel.items[initialModel.chosen].text, nil
	}
	return "", fmt.Errorf("no selection made")
}
