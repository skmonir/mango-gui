package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/skmonir/mango-gui/context"
)

func GetMainWindowContainer(app fyne.App, MainWindow fyne.Window, ctx *context.AppCtx) *fyne.Container {
	headerUI := GetHeaderUI(app, MainWindow, ctx)
	footerUI := GetFooterUI(MainWindow, ctx)

	tabs := container.NewAppTabs(
		// container.NewTabItem("Tab 1", widget.NewLabel("Hello")),
		// container.NewTabItemWithIcon("Home", theme.HomeIcon(), widget.NewLabel("Home tab")),
		container.NewTabItemWithIcon("Parser", theme.DownloadIcon(), GetParserUI(MainWindow, ctx)),
		// container.NewTabItemWithIcon("Tester", theme.DownloadIcon(), ui.GetTesterUI(MainWindow)),
		// container.NewTabItemWithIcon("Input Generator", theme.SettingsIcon(), container.New(layout.NewVBoxLayout(), content, centered)),
		// container.NewTabItemWithIcon("Output Generator", theme.SettingsIcon(), container.New(layout.NewVBoxLayout(), content, centered)),
		// container.NewTabItemWithIcon("Settings", theme.SettingsIcon(), container.New(layout.NewVBoxLayout(), content, centered)),
		// container.NewTabItemWithIcon("About", theme.InfoIcon(), widget.NewLabel("About")),
	)

	mainWindowContainer := container.NewVBox(
		headerUI,
		tabs,
		widget.NewSeparator(),
		layout.NewSpacer(),
		footerUI,
	)

	return mainWindowContainer
}
