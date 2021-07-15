package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"github.com/skmonir/mango-gui/context"
	"github.com/skmonir/mango-gui/ui"
)

func initAppVars(app fyne.App, ctx *context.AppCtx) {
	ctx.Config = context.GetAppConfig()
	if ctx.Config != nil && ctx.Config.CurrentTheme == "light" {
		app.Settings().SetTheme(theme.LightTheme())
	} else {
		app.Settings().SetTheme(theme.DarkTheme())
	}
	ctx.OnlineJudgeOptions = &[]string{
		"CodeForces",
		// "AtCoder",
	}
}

func main() {
	app := app.New()
	MainWindow := app.NewWindow("Mango - Task Parser and Tester")

	var ctx *context.AppCtx = &context.AppCtx{}

	initAppVars(app, ctx)

	mainWindowContainer := ui.GetMainWindowContainer(app, MainWindow, ctx)

	MainWindow.SetContent(mainWindowContainer)
	MainWindow.Resize(fyne.NewSize(1000, 720))
	MainWindow.SetFixedSize(true)
	MainWindow.CenterOnScreen()
	MainWindow.ShowAndRun()
}
