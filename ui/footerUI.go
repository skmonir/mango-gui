package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/skmonir/mango-gui/context"
)

func GetFooterUI(MainWindow fyne.Window, ctx *context.AppCtx) *fyne.Container {
	ctx.ProgressBar = &widget.ProgressBar{
		Min:   0,
		Max:   100,
		Value: 0,
	}

	footerContainer := container.New(layout.NewVBoxLayout(),
		ctx.ProgressBar,
	)

	return footerContainer
}
