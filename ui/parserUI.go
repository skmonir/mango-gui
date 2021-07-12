package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/skmonir/mango-gui/context"
	"github.com/skmonir/mango-gui/system/taskParser"
)

func GetParserUI(MainWindow fyne.Window, ctx *context.AppCtx) *fyne.Container {
	parserUI := context.ParserUiCtx{}

	var parserContainer *fyne.Container
	var ListContainer *fyne.Container
	var list = []string{}

	parserUI.OnlineJudgeOptions = widget.NewSelect([]string{
		"CodeForces",
		"AtCoder",
	}, func(s string) {})
	parserUI.OnlineJudgeOptions.PlaceHolder = "Select Online Judge"

	parserUI.ContestIdInputField = widget.NewEntry()
	parserUI.ContestIdInputField.SetPlaceHolder("Enter Contest/Problem ID")

	parseButton := widget.NewButtonWithIcon("Parse", theme.DownloadIcon(), func() {
		// // system.Parse(cfg, "1521", &list, ListContainer)
		// ctx.Config.CurrentContestId = "1234"
		// ctx.HeaderUi.CurrentContestField.SetText("1234")
		// ctx.ProgressBar.SetValue(100)
		if err := taskParser.Parse(ctx); err != nil {
			//something
		} else {
			ctx.HeaderUi.CurrentContestField.SetText(ctx.Config.CurrentContestId)
		}
	})

	createButton := widget.NewButtonWithIcon("Create", theme.DocumentCreateIcon(), func() {
		fmt.Println("Some code TODO")
	})

	components := widget.NewList(
		func() int {
			return len(list)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			fmt.Println(i)
			o.(*widget.Label).SetText(list[i])
		},
	)

	ListContainer = container.NewGridWrap(fyne.NewSize(1000, 550), components)

	parserContainer = container.New(layout.NewVBoxLayout(),
		container.NewGridWithColumns(
			4,
			parserUI.OnlineJudgeOptions,
			parserUI.ContestIdInputField,
			parseButton,
			createButton,
		),
		widget.NewSeparator(),
		ListContainer,
	)

	ctx.ParserUi = &parserUI
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
