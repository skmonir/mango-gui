package ui

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/skmonir/mango-gui/context"
	"github.com/skmonir/mango-gui/models"
)

func GetTesterUI(MainWindow fyne.Window, ctx *context.AppCtx) *fyne.Container {
	ctx.TesterUi = &context.TesterUiCtx{
		ProblemNameList: &[]string{},
	}

	ctx.TesterUi.OnlineJudgeOptions = widget.NewSelect([]string{
		"CodeForces",
		"AtCoder",
	}, nil)
	ctx.TesterUi.OnlineJudgeOptions.PlaceHolder = "Select Online Judge"

	ctx.TesterUi.ContestIdInputField = widget.NewEntry()
	ctx.TesterUi.ContestIdInputField.SetPlaceHolder("Enter Contest/Problem ID")

	createContestButton := widget.NewButtonWithIcon("Create Contest", theme.ContentAddIcon(), func() {
		GetProblemInfoList(ctx)
	})

	ctx.TesterUi.ProblemNameListSelect = widget.NewSelect(*ctx.TesterUi.ProblemNameList, nil)
	ctx.TesterUi.ProblemNameListSelect.PlaceHolder = "Select Problem"
	ctx.TesterUi.ProblemNameListSelect.Hide()

	ctx.TesterUi.RunTestsButton = widget.NewButtonWithIcon("Run Tests", theme.ContentAddIcon(), func() {
		// *ctx.TesterUi.ProblemNameList = append(*ctx.TesterUi.ProblemNameList, "C")
	})
	ctx.TesterUi.RunTestsButton.Hide()

	testerContainer := container.New(layout.NewVBoxLayout(),
		container.NewGridWithColumns(
			4,
			ctx.TesterUi.OnlineJudgeOptions,
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
	if ctx.TesterUi.OnlineJudgeOptions.Selected == "" {
		return errors.New("please select Online Judge before continuing")
	} else if ctx.TesterUi.ContestIdInputField.Text == "" {
		return errors.New("please enter valid Contest ID before continuing")
	}
	return nil
}

func GetProblemInfoList(ctx *context.AppCtx) []models.Problem {
	problemIdList := ctx.Config.GetProblemIdListForTester()
	var problemInfoList []models.Problem
	var problemNameList []string
	for _, problemId := range problemIdList {
		problemInfo, _ := ctx.Config.GetProblemInfo(problemId)
		problemInfoList = append(problemInfoList, problemInfo)
		problemNameList = append(problemNameList, problemInfo.Name)
	}

	if len(problemNameList) > 0 {
		ctx.TesterUi.ProblemNameListSelect.Show()
		ctx.TesterUi.ProblemNameListSelect.SetOptions(problemNameList)
		ctx.TesterUi.RunTestsButton.Show()
	}

	return problemInfoList
}
