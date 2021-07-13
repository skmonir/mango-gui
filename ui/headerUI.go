package ui

import (
	"fmt"

	"fyne.io/fyne/v2/dialog"
	"github.com/skmonir/mango-gui/context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func GetHeaderUI(app fyne.App, MainWindow fyne.Window, ctx *context.AppCtx) *fyne.Container {
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

	headerUI.CurrentOnlineJudge = widget.NewLabel(ctx.Config.OJ)
	headerUI.CurrentContestField = widget.NewLabel(ctx.Config.CurrentContestId)

	themeToggler := widget.NewCheck("Dark Mode", func(isChecked bool) {
		toggleAppTheme(app, ctx)
	})

	CurrentContestLabel := container.NewGridWithColumns(3,
		widget.NewForm(widget.NewFormItem("OJ", headerUI.CurrentOnlineJudge)),
		widget.NewForm(widget.NewFormItem("Contest ID", headerUI.CurrentContestField)),
		themeToggler,
	)

	headerContainer := container.New(layout.NewVBoxLayout(),
		container.NewGridWithColumns(
			2,
			headerUI.WorkSpaceDirChooser,
			CurrentContestLabel,
		),
		widget.NewSeparator(),
		widget.NewSeparator(),
	)

	initAppTheme(app, themeToggler, ctx)

	ctx.HeaderUi = &headerUI

	return headerContainer
}

func initAppTheme(app fyne.App, themeToggler *widget.Check, ctx *context.AppCtx) {
	if ctx.Config.CurrentTheme == "dark" {
		themeToggler.SetChecked(true)
		app.Settings().SetTheme(theme.DarkTheme())
	} else {
		app.Settings().SetTheme(theme.LightTheme())
	}
}

func toggleAppTheme(app fyne.App, ctx *context.AppCtx) {
	fmt.Println("fff")
	if ctx.Config.CurrentTheme == "light" {
		ctx.Config.CurrentTheme = "dark"
		app.Settings().SetTheme(theme.DarkTheme())
	} else {
		ctx.Config.CurrentTheme = "light"
		app.Settings().SetTheme(theme.LightTheme())
	}
	ctx.Config.SaveConfig()
}
