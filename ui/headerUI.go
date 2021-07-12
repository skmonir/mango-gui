package ui

import (
	"fyne.io/fyne/v2/dialog"
	"github.com/skmonir/mango-gui/context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func GetHeaderUI(MainWindow fyne.Window, ctx *context.AppCtx) *fyne.Container {
	headerUI := context.HeaderUiCtx{}

	workspaceDirPath := "Select workspace directory"
	if ctx.Config.Workspace != "" {
		workspaceDirPath = ctx.Config.Workspace
	}

	headerUI.WorkSpaceDirChooser = widget.NewButtonWithIcon(workspaceDirPath, theme.FolderIcon(), func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, MainWindow)
			} else if uri != nil {
				ctx.Config.Workspace = uri.Path()
				ctx.Config.SaveConfig()
				headerUI.WorkSpaceDirChooser.SetText(uri.Path())
			}
		}, MainWindow)
	})

	headerUI.CurrentContestField = widget.NewLabel(ctx.Config.CurrentContestId)
	CurrentContestLabel := widget.NewForm(widget.NewFormItem("Current Working Contest ID", headerUI.CurrentContestField))

	headerContainer := container.New(layout.NewVBoxLayout(),
		container.NewGridWithColumns(
			2,
			headerUI.WorkSpaceDirChooser,
			CurrentContestLabel,
		),
		widget.NewSeparator(),
		widget.NewSeparator(),
	)

	ctx.HeaderUi = &headerUI

	return headerContainer
}
