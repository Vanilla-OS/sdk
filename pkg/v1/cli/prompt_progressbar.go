package cli

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Progress bar prompt implementation using Bubble Tea (progress).
*/

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type progressMsg float64
type titleMsg string
type stopMsg struct{}

type progressComponent struct {
	progress progress.Model
	total    int
	current  int
	message  string
}

func (m progressComponent) Init() tea.Cmd {
	return nil
}

func (m progressComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	case progressMsg:
		var cmd tea.Cmd
		if m.current < m.total {
			pct := float64(m.current) / float64(m.total)
			cmd = m.progress.SetPercent(pct)
			return m, cmd
		}
	case titleMsg:
		m.message = string(msg)
	case stopMsg:
		return m, tea.Quit
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}
	return m, nil
}

func (m progressComponent) View() string {
	pad := strings.Repeat(" ", 2)
	return "\n" +
		pad + m.message + "\n" +
		pad + m.progress.View() + "\n\n"
}

type ProgressBarModel struct {
	program *tea.Program
	total   int
	current int
}

// StartProgressBar starts a progress bar with a message and a total.
// The progress bar is stopped automatically when it reaches the total or
// manually by calling the Stop method on the returned model.
//
// Example:
//
//	progressBar := myApp.CLI.StartProgressBar("Loading the batmobile...", 100)
//	for i := 0; i < 100; i++ {
//		progressBar.Increment(1)
//		time.Sleep(50 * time.Millisecond)
//	}
func (c Command) StartProgressBar(message string, total int) *ProgressBarModel {
	p := progress.New(
		progress.WithGradient("#277eff", "#e0388d"),
		progress.WithoutPercentage(),
	)
	m := progressComponent{
		progress: p,
		total:    total,
		message:  message,
	}

	prog := tea.NewProgram(m)

	go func() {
		if _, err := prog.Run(); err != nil {
			fmt.Println("Error running progress bar:", err)
		}
	}()

	return &ProgressBarModel{
		program: prog,
		total:   total,
		current: 0,
	}
}

func (m *ProgressBarModel) Increment(inc int) {
	m.current += inc
	if m.current >= m.total {
		m.current = m.total
		m.Stop()
		return
	}

	// Bubbles progress takes a 0-1 float.
	pct := float64(m.current) / float64(m.total)
	m.program.Send(progressMsg(pct))
}

// UpdateMessage updates the title logic.
func (m *ProgressBarModel) UpdateMessage(msg string) {
	m.program.Send(titleMsg(msg))
}

func (m *ProgressBarModel) Stop() {
	m.program.Send(stopMsg{})
	// Allow cleanup
	time.Sleep(100 * time.Millisecond)
}
