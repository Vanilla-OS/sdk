package cli

import (
	"github.com/pterm/pterm"
)

type SpinnerModel struct {
	spinner  *pterm.SpinnerPrinter
	message  string
	finished bool
}

func newSpinnerModel(message string) *SpinnerModel {
	spinner, _ := pterm.DefaultSpinner.Start(message)
	return &SpinnerModel{
		spinner: spinner,
		message: message,
	}
}

func (m *SpinnerModel) UpdateMessage(message string) {
	m.spinner.UpdateText(message)
}

// Stop stops the spinner and marks it as finished.
//
// Example:
//
//	spinner := myApp.CLI.StartSpinner("Loading the batmobile...")
//	time.Sleep(3 * time.Second)
//	spinner.Stop()
func (m *SpinnerModel) Stop() {
	if !m.finished {
		m.spinner.Success()
		m.finished = true
	}
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
	model := newSpinnerModel(message)
	return model
}
