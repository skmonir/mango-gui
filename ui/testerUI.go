package ui

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/skmonir/mango-gui/context"
	"github.com/skmonir/mango-gui/system"
)

func GetTesterUI(MainWindow fyne.Window, ctx *context.AppCtx) *fyne.Container {
	ctx.TesterUi = &context.TesterUiCtx{}

	// var table *fyne.Container
	table := container.NewHBox()
	// table := container.NewScroll(container.NewHBox())
	var testerContainer *fyne.Container

	ctx.TesterUi.OnlineJudgeOptionSelect = widget.NewSelect(*ctx.OnlineJudgeOptions, nil)
	ctx.TesterUi.OnlineJudgeOptionSelect.PlaceHolder = "Select Online Judge"

	ctx.TesterUi.ContestIdInputField = widget.NewEntry()
	ctx.TesterUi.ContestIdInputField.SetPlaceHolder("Enter Contest/Problem ID")

	createContestButton := widget.NewButtonWithIcon("Create Contest/Problem", theme.ContentAddIcon(), func() {
		if err := ValidateTesterUiFields(ctx); err != nil {
			dialog.ShowError(err, MainWindow)
			return
		}
		problems := system.Create(ctx)
		if len(problems) == 0 {
			dialog.ShowError(errors.New("no problems found"), MainWindow)
		}
		ProcessProblems(ctx, problems)
	})

	ctx.TesterUi.ProblemNameListSelect = widget.NewSelect([]string{}, func(selected string) {
		testcaseList, ok := (*ctx.TesterUi.TestcaseList)[selected]
		if ok {
			testerContainer.Remove(table)
			table = container.NewGridWrap(fyne.NewSize(1000, 500), container.NewScroll(GetNewTable(testcaseList)))
			testerContainer.Add(table)
		}
	})
	ctx.TesterUi.ProblemNameListSelect.PlaceHolder = "Select Problem"
	ctx.TesterUi.ProblemNameListSelect.Hide()

	ctx.TesterUi.RunTestsButton = widget.NewButtonWithIcon("Run Tests", theme.ContentAddIcon(), func() {
		if err := RunTestClicked(ctx); err != nil {
			dialog.ShowError(err, MainWindow)
			return
		}
		testcaseList, ok := (*ctx.TesterUi.TestcaseList)[ctx.TesterUi.ProblemNameListSelect.Selected]
		if ok {
			testerContainer.Remove(table)
			table = container.NewGridWrap(fyne.NewSize(1000, 500), container.NewScroll(GetNewTable(testcaseList)))
			testerContainer.Add(table)
		}
	})
	ctx.TesterUi.RunTestsButton.Hide()

	testerContainer = container.New(layout.NewVBoxLayout(),
		container.NewGridWithColumns(
			4,
			ctx.TesterUi.OnlineJudgeOptionSelect,
			ctx.TesterUi.ContestIdInputField,
			createContestButton,
		),
		widget.NewSeparator(),
		container.NewGridWithColumns(
			2,
			ctx.TesterUi.ProblemNameListSelect,
			container.NewGridWithColumns(
				2,
				ctx.TesterUi.RunTestsButton,
				&widget.Label{},
			),
		),
		table,
	)

	return testerContainer
}
