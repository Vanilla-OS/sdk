package cli

import (
	"github.com/pterm/pterm"
)

// ProgressBarModel represents a progress bar model
type ProgressBarModel struct {
	progressBar *pterm.ProgressbarPrinter
	finished    bool
}

// newProgressBarModel creates a new progress bar model
func newProgressBarModel() *ProgressBarModel {
	return &ProgressBarModel{}
}

// UpdateMessage updates the message of the progress bar
func (m *ProgressBarModel) UpdateMessage(message string) {
	m.progressBar.UpdateTitle(message)
}

// UpdateProgress increments the progress of the progress bar only if it
// has not reached the total.
//
// Example:
//
//	progressBar := myApp.CLI.StartProgressBar("Loading the batmobile...", 100)
//	progressBar.Increment(50)
func (m *ProgressBarModel) Increment(progress int) {
	if m.progressBar.Total != m.progressBar.Current {
		m.progressBar.Add(progress)
	} else {
		m.Stop()
	}
}

// Stop stops the progress bar and marks it as finished.
//
// Example:
//
//	progressBar := myApp.CLI.StartProgressBar("Loading the batmobile...", 100)
//	progressBar.Increment(50)
//	progressBar.UpdateMessage("Failed to load the batmobile")
//	progressBar.Stop()
func (m *ProgressBarModel) Stop() {
	if !m.finished {
		m.progressBar.Stop()
		m.finished = true
	}
}

// StartProgressBar starts a progress bar with a message and a total
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
	model := newProgressBarModel()
	model.progressBar, _ = pterm.DefaultProgressbar.WithTotal(total).WithTitle(message).Start()
	return model
}
