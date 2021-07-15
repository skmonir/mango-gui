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
	"github.com/skmonir/mango-gui/models"
	"github.com/skmonir/mango-gui/system"
)

func GetTesterUI(MainWindow fyne.Window, ctx *context.AppCtx) *fyne.Container {
	ctx.TesterUi = &context.TesterUiCtx{}

	ctx.TesterUi.OnlineJudgeOptionSelect = widget.NewSelect(*ctx.OnlineJudgeOptions, nil)
	ctx.TesterUi.OnlineJudgeOptionSelect.PlaceHolder = "Select Online Judge"

	ctx.TesterUi.ContestIdInputField = widget.NewEntry()
	ctx.TesterUi.ContestIdInputField.SetPlaceHolder("Enter Contest/Problem ID")

	createContestButton := widget.NewButtonWithIcon("Create Contest/Problem", theme.ContentAddIcon(), func() {
		if err := validateTesterUiFields(ctx); err != nil {
			dialog.ShowError(err, MainWindow)
			return
		}
		problems := system.Create(ctx)
		if len(problems) == 0 {
			dialog.ShowError(errors.New("no problems found"), MainWindow)
		}
		processProblems(ctx, problems)
	})

	ctx.TesterUi.ProblemNameListSelect = widget.NewSelect([]string{}, nil)
	ctx.TesterUi.ProblemNameListSelect.PlaceHolder = "Select Problem"
	ctx.TesterUi.ProblemNameListSelect.Hide()

	ctx.TesterUi.RunTestsButton = widget.NewButtonWithIcon("Run Tests", theme.ContentAddIcon(), func() {
		// *ctx.TesterUi.ProblemNameList = append(*ctx.TesterUi.ProblemNameList, "C")
	})
	ctx.TesterUi.RunTestsButton.Hide()

	testerContainer := container.New(layout.NewVBoxLayout(),
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
	)

	return testerContainer
}

func validateTesterUiFields(ctx *context.AppCtx) error {
	if ctx.TesterUi.OnlineJudgeOptionSelect.Selected == "" {
		return errors.New("please select Online Judge before continuing")
	} else if ctx.TesterUi.ContestIdInputField.Text == "" {
		return errors.New("please enter valid Contest ID before continuing")
	}
	return nil
}

func processProblems(ctx *context.AppCtx, problems []models.Problem) {
	ctx.TesterUi.ProblemList = &[]models.Problem{}
	ctx.TesterUi.ProblemNameList = &[]string{}

	for _, problem := range problems {
		*ctx.TesterUi.ProblemList = append(*ctx.TesterUi.ProblemList, problem)
		*ctx.TesterUi.ProblemNameList = append(*ctx.TesterUi.ProblemNameList, problem.Name)
	}

	ctx.TesterUi.ProblemNameListSelect.ClearSelected()
	ctx.TesterUi.ProblemNameListSelect.SetOptions(*ctx.TesterUi.ProblemNameList)

	if len(problems) > 0 {
		ctx.TesterUi.ProblemNameListSelect.Show()
		ctx.TesterUi.RunTestsButton.Show()
	}
}
