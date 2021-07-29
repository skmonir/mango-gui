package ui

import (
	"errors"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/skmonir/mango-gui/context"
	"github.com/skmonir/mango-gui/system"
)

func GetParserUI(MainWindow fyne.Window, ctx *context.AppCtx) *fyne.Container {
	ctx.ParserUi = &context.ParserUiCtx{
		ParsedProblemStatus: &[]string{},
	}

	ctx.ParserUi.OnlineJudgeOptionSelect = widget.NewSelect(*ctx.OnlineJudgeOptions, func(s string) {})
	ctx.ParserUi.OnlineJudgeOptionSelect.PlaceHolder = "Select Online Judge"

	ctx.ParserUi.ContestIdInputField = widget.NewEntry()
	ctx.ParserUi.ContestIdInputField.SetPlaceHolder("Enter Contest/Problem ID")

	parseButton := widget.NewButtonWithIcon("Parse Samples", theme.DownloadIcon(), func() {
		ctx.ProgressBar.SetValue(0)
		if err := validateParserUiFields(ctx); err != nil {
			dialog.ShowError(err, MainWindow)
			return
		}
		if err := system.Parse(ctx); err != nil {
			dialog.ShowError(err, MainWindow)
		}
	})

	createButton := widget.NewButtonWithIcon("Create Sources", theme.FileTextIcon(), func() {
		ctx.ProgressBar.SetValue(0)
		if err := validateParserUiFields(ctx); err != nil {
			dialog.ShowError(err, MainWindow)
			return
		}
		if err := system.Source(ctx); err != nil {
			dialog.ShowError(err, MainWindow)
		}
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
	parsedProblemList.OnSelected = func(id int) {
		selected := (*ctx.ParserUi.ParsedProblemStatus)[id]
		prefix := "[CREATED] Successfully created source "
		if strings.HasPrefix(selected, prefix) {
			srcFileName := strings.TrimPrefix(selected, prefix)
			srcFileName = strings.TrimSuffix(srcFileName, ". Click to open the source file.")
			srcFileName = strings.Trim(srcFileName, " \n")
			message := "Do you want to open " + srcFileName + " in default editor?"
			dialog.ShowConfirm("Source Open Confirmation", message, func(response bool) {
				if response {
					if err := system.Open(ctx, ctx.ParserUi.CurrentContestId, strings.TrimSuffix(srcFileName, ".cpp")); err != nil {
						dialog.ShowError(err, MainWindow)
					}
				}
			}, MainWindow)
		}
	}

	ctx.ParserUi.ParsedProblemListContainer = container.NewGridWrap(fyne.NewSize(1000, 550), parsedProblemList)

	parserContainer := container.New(layout.NewVBoxLayout(),
		container.NewGridWithColumns(
			4,
			ctx.ParserUi.OnlineJudgeOptionSelect,
			ctx.ParserUi.ContestIdInputField,
			parseButton,
			createButton,
		),
		widget.NewSeparator(),
		ctx.ParserUi.ParsedProblemListContainer,
	)

	return parserContainer
}

func validateParserUiFields(ctx *context.AppCtx) error {
	if ctx.ParserUi.OnlineJudgeOptionSelect.Selected == "" {
		return errors.New("please select Online Judge before continuing")
	} else if ctx.ParserUi.ContestIdInputField.Text == "" {
		return errors.New("please enter valid Contest ID before continuing")
	}
	return nil
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
