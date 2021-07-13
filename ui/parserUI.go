package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/skmonir/mango-gui/context"
	"github.com/skmonir/mango-gui/system/taskParser"
)

func GetParserUI(MainWindow fyne.Window, ctx *context.AppCtx) *fyne.Container {
	ctx.ParserUi = &context.ParserUiCtx{
		ParsedProblemStatus: &[]string{},
	}

	ctx.ParserUi.OnlineJudgeOptions = widget.NewSelect([]string{
		"CodeForces",
		"AtCoder",
	}, func(s string) {})
	ctx.ParserUi.OnlineJudgeOptions.PlaceHolder = "Select Online Judge"

	ctx.ParserUi.ContestIdInputField = widget.NewEntry()
	ctx.ParserUi.ContestIdInputField.SetPlaceHolder("Enter Contest/Problem ID")

	parseButton := widget.NewButtonWithIcon("Parse", theme.DownloadIcon(), func() {
		ctx.ProgressBar.SetValue(0)
		if err := taskParser.Parse(ctx); err != nil {
			dialog.ShowError(err, MainWindow)
		} else {
			ctx.HeaderUi.CurrentContestField.SetText(ctx.Config.CurrentContestId)
		}
	})

	createButton := widget.NewButtonWithIcon("Create", theme.DocumentCreateIcon(), func() {
		ctx.ProgressBar.SetValue(0)
		fmt.Println("Some code TODO")
	})

	parsedProblemList := widget.NewList(
		func() int {
			return len(*ctx.ParserUi.ParsedProblemStatus)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(index widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText((*ctx.ParserUi.ParsedProblemStatus)[index])
		},
	)

	ctx.ParserUi.ParsedProblemListContainer = container.NewGridWrap(fyne.NewSize(1000, 550), parsedProblemList)

	parserContainer := container.New(layout.NewVBoxLayout(),
		container.NewGridWithColumns(
			4,
			ctx.ParserUi.OnlineJudgeOptions,
			ctx.ParserUi.ContestIdInputField,
			parseButton,
			createButton,
		),
		widget.NewSeparator(),
		ctx.ParserUi.ParsedProblemListContainer,
	)

	return parserContainer
}

// func makeList(rows []string) *fyne.Container {
// 	var objects []fyne.CanvasObject
// 	for _, val := range rows {
// 		btn := widget.NewButtonWithIcon(val, theme.ConfirmIcon(), func() {})
// 		objects = append(objects, btn)
// 	}
// 	return container.NewVBox(
// 		container.NewGridWithColumns(1, objects...),
// 	)
// }

// create table https://www.gitmemory.com/issue/fyne-io/fyne/157/476708251
