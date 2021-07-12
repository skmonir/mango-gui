package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/skmonir/mango-gui/context"
	"github.com/skmonir/mango-gui/ui"
)

func main() {
	app := app.New()
	MainWindow := app.NewWindow("Mango - Task Parser & Tester")

	var ctx *context.AppCtx = &context.AppCtx{
		Config: context.GetAppConfig(),
	}

	mainWindowContainer := ui.GetMainWindowContainer(MainWindow, ctx)

	MainWindow.SetContent(mainWindowContainer)
	MainWindow.Resize(fyne.NewSize(1000, 720))
	MainWindow.SetFixedSize(true)
	MainWindow.CenterOnScreen()
	MainWindow.ShowAndRun()
}
