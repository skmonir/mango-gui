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

func GetTesterUI(MainWindow fyne.Window, ctx *context.AppCtx) *fyne.Container {
	ctx.TesterUi = &context.TesterUiCtx{}

	table := container.NewScroll(container.NewHBox())
	var testerContainer *fyne.Container

	ctx.TesterUi.OnlineJudgeOptionSelect = widget.NewSelect(*ctx.OnlineJudgeOptions, nil)
	ctx.TesterUi.OnlineJudgeOptionSelect.PlaceHolder = "Select Online Judge"

	ctx.TesterUi.ContestIdInputField = widget.NewEntry()
	ctx.TesterUi.ContestIdInputField.SetPlaceHolder("Enter Contest/Problem ID")

	createContestButton := widget.NewButtonWithIcon("Load Contest/Problem", theme.FolderOpenIcon(), func() {
		if err := ValidateTesterUiFields(ctx); err != nil {
			dialog.ShowError(err, MainWindow)
			return
		}
		problems := system.Create(ctx)
		if len(problems) == 0 {
			dialog.ShowError(errors.New("no problems found"), MainWindow)
			return
		}
		ProcessProblems(ctx, problems)
		testerContainer.Remove(table)
		table = container.NewScroll(container.NewHBox())
		testerContainer.Add(table)
	})

	ctx.TesterUi.ProblemNameListSelect = widget.NewSelect([]string{}, func(selected string) {
		testcaseList, ok := (*ctx.TesterUi.TestcaseList)[selected]
		if ok {
			testerContainer.Remove(table)
			table = GetNewTable(testcaseList)
			testerContainer.Add(table)
		}
	})
	ctx.TesterUi.ProblemNameListSelect.PlaceHolder = "Select Problem"
	ctx.TesterUi.ProblemNameListSelect.Hide()

	ctx.TesterUi.RunTestsButton = widget.NewButtonWithIcon("Run Tests", theme.MediaPlayIcon(), func() {
		if ctx.TesterUi.ProblemNameListSelect.Selected == "" {
			return
		}
		if err := RunTestClicked(ctx); err != nil {
			dialog.ShowError(err, MainWindow)
			return
		}
		testcaseList, ok := (*ctx.TesterUi.TestcaseList)[ctx.TesterUi.ProblemNameListSelect.Selected]
		if ok {
			testerContainer.Remove(table)
			table = GetNewTable(testcaseList)
			testerContainer.Add(table)
		}
	})
	ctx.TesterUi.RunTestsButton.Hide()

	// ctx.TesterUi.AddTestButton = widget.NewButton("Add Test", func() {
	// })
	// ctx.TesterUi.AddTestButton.Hide()

	// ctx.TesterUi.RemoveTestButton = widget.NewButton("Delete Test", func() {
	// })
	// ctx.TesterUi.RemoveTestButton.Hide()

	ctx.TesterUi.OpenSourceButton = widget.NewButton("Open Source", func() {
		selected := ctx.TesterUi.ProblemNameListSelect.Selected
		if selected == "" {
			return
		}
		problemId := strings.Split(selected, ".")[0]
		problemId = strings.Trim(problemId, " \n")

		if err := system.Open(ctx, ctx.TesterUi.CurrentContestId, problemId); err != nil {
			dialog.ShowError(err, MainWindow)
		}
	})
	ctx.TesterUi.OpenSourceButton.Hide()

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
				4,
				ctx.TesterUi.RunTestsButton,
				// ctx.TesterUi.AddTestButton,
				// ctx.TesterUi.RemoveTestButton,
				ctx.TesterUi.OpenSourceButton,
			),
		),
		table,
	)

	return testerContainer
}
