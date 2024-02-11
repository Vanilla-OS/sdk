package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	vApp "github.com/vanilla-os/sdk/pkg/v1/app"
	vAppTypes "github.com/vanilla-os/sdk/pkg/v1/app/types"
	"github.com/vanilla-os/sdk/pkg/v1/system"
	sysTypes "github.com/vanilla-os/sdk/pkg/v1/system/types"
)

var (
	myApp        *vApp.App
	currentPID   int
	processList  []sysTypes.Process
	processTable *widget.Table
)

func main() {
	myApp, _ = vApp.NewApp(vAppTypes.AppOptions{
		RDNN:    "com.vanillaos.taskmanager",
		Name:    "Task Manager",
		Version: "1.0.0",
	})
	myFyneApp := app.NewWithID("Task Manager")
	window := NewWindow(myFyneApp)

	var err error
	processList, err = system.GetProcessList()
	if err != nil {
		dialog.ShowError(err, window)
		myApp.Log.Term.Error().Msgf("Error: %v", err)
		return
	}

	processTable = createProcessTable(&processList)
	toolbar := createToolbar(&processList, window)
	content := container.NewBorder(toolbar, nil, nil, nil, processTable)

	myApp.Log.Term.Info().Msg("Starting Task Manager")

	window.SetContent(content)
	window.Resize(fyne.NewSize(600, 500))
	window.ShowAndRun()

	myApp.Log.Term.Info().Msg("Task Manager closed")
}

func NewWindow(app fyne.App) fyne.Window {
	myWindow := app.NewWindow("Task Manager")
	return myWindow
}

func createProcessTable(processList *[]sysTypes.Process) *widget.Table {
	table := widget.NewTable(
		func() (int, int) { return len(*processList), 5 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			process := (*processList)[id.Row]
			switch id.Col {
			case 0:
				cell.(*widget.Label).SetText(fmt.Sprintf("%d", process.PID))
			case 1:
				cell.(*widget.Label).SetText(process.Name)
			case 2:
				cell.(*widget.Label).SetText(process.State)
			case 3:
				cell.(*widget.Label).SetText(fmt.Sprintf("%d", process.PPID))
			case 4:
				cell.(*widget.Label).SetText(fmt.Sprintf("%d", process.Priority))
			}
		},
	)
	table.SetColumnWidth(0, 50)  // PID
	table.SetColumnWidth(1, 200) // Name
	table.SetColumnWidth(2, 100) // State
	table.SetColumnWidth(3, 50)  // PPID
	table.SetColumnWidth(4, 50)  // Priority

	table.OnSelected = func(id widget.TableCellID) {
		currentPID = (*processList)[id.Row].PID
	}
	return table
}

func createToolbar(processList *[]sysTypes.Process, window fyne.Window) *fyne.Container {
	procsLabel := widget.NewLabel("Processes: 0")
	updateProcsLabel := func() {
		procsLabel.SetText(fmt.Sprintf("Processes: %d", len(*processList)))
	}
	refreshFn := func() {
		var err error
		*processList, err = system.GetProcessList()
		myApp.Log.Term.Info().Msgf("Refreshed process list: %d", len(*processList))
		if err != nil {
			dialog.ShowError(err, window)
			myApp.Log.Term.Error().Msgf("Error: %v", err)
			return
		}
		processTable.Refresh()
		updateProcsLabel()
	}
	refreshAction := widget.NewToolbarAction(
		theme.ViewRefreshIcon(),
		refreshFn,
	)
	deleteAction := widget.NewToolbarAction(
		theme.DeleteIcon(),
		func() {
			if currentPID != 0 {
				err := system.KillProcess(currentPID)
				if err != nil {
					dialog.ShowError(err, window)
					myApp.Log.Term.Error().Msgf("Error: %v", err)
				}
				refreshFn()
			}
		},
	)

	toolbar := widget.NewToolbar(
		refreshAction,
		deleteAction,
	)
	hbox := container.NewHBox(
		toolbar,
		layout.NewSpacer(),
		procsLabel,
	)

	updateProcsLabel()

	return hbox
}
