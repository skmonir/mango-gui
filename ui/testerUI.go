package ui

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

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

	// var table *fyne.Container
	table := container.NewHBox()
	// table := container.NewScroll(container.NewHBox())
	var testerContainer *fyne.Container

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

	ctx.TesterUi.ProblemNameListSelect = widget.NewSelect([]string{}, func(selected string) {
		testcaseList, ok := (*ctx.TesterUi.TestcaseList)[selected]
		if ok {
			testerContainer.Remove(table)
			// table = container.NewScroll(getNewTable(testcaseList))
			table = container.NewGridWrap(fyne.NewSize(1000, 500), container.NewScroll(getNewTable(testcaseList)))
			// table = container.NewScroll(getNewTable(testcaseList))
			testerContainer.Add(table)
		}
	})
	ctx.TesterUi.ProblemNameListSelect.PlaceHolder = "Select Problem"
	ctx.TesterUi.ProblemNameListSelect.Hide()

	ctx.TesterUi.RunTestsButton = widget.NewButtonWithIcon("Run Tests", theme.ContentAddIcon(), func() {
		// testerContainer.Remove(table)
		// new table
		// testerContainer.Add(table)
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
	testCaseMap := make(map[string][]models.ExecutionResult)

	for _, problem := range problems {
		*ctx.TesterUi.ProblemList = append(*ctx.TesterUi.ProblemList, problem)
		*ctx.TesterUi.ProblemNameList = append(*ctx.TesterUi.ProblemNameList, problem.Name)

		for _, testcase := range problem.Dataset {
			execResult := models.ExecutionResult{
				Test: testcase,
			}
			testCaseMap[problem.Name] = append(testCaseMap[problem.Name], execResult)
		}
	}

	ctx.TesterUi.TestcaseList = &testCaseMap

	ctx.TesterUi.ProblemNameListSelect.ClearSelected()
	ctx.TesterUi.ProblemNameListSelect.SetOptions(*ctx.TesterUi.ProblemNameList)

	if len(problems) > 0 {
		ctx.TesterUi.ProblemNameListSelect.Show()
		ctx.TesterUi.RunTestsButton.Show()
	}
}

func getNewTable(testcaseList []models.ExecutionResult) *fyne.Container {
	rowHeights := getRowHeights(testcaseList)

	rows := make([][]string, len(testcaseList))

	for i, testcase := range testcaseList {
		mxH := rowHeights[i]
		rows[i] = append(rows[i], getHeightAdjustedCell(mxH, strconv.Itoa(i)))
		rows[i] = append(rows[i], getHeightAdjustedCell(mxH, testcase.Test.Input))
		rows[i] = append(rows[i], getHeightAdjustedCell(mxH, testcase.Test.Output))
		rows[i] = append(rows[i], getHeightAdjustedCell(mxH, testcase.Output))
		rows[i] = append(rows[i], getHeightAdjustedCell(mxH, testcase.Verdict))
		rows[i] = append(rows[i], getHeightAdjustedCell(mxH, fmt.Sprintf("%v ms", testcase.Runtime)))
		rows[i] = append(rows[i], getHeightAdjustedCell(mxH, "0 KB"))
	}

	headings := []string{"#", "SAMPLE INPUT", "SAMPLE OUTPUT", "PROGRAM OUTPUT", "VERDICT", "TIME", "MEMORY"}

	return makeTable(headings, rows)
}

func getHeightAdjustedCell(height int, value string) string {
	for len(strings.Split(value, "\n")) < height {
		value += "\n "
	}
	return value
}

func getRowHeights(testcaseList []models.ExecutionResult) []int {
	var rowHeights []int
	for _, testcase := range testcaseList {
		mxH := math.Max(getLineCounts(testcase.Test.Input), getLineCounts(testcase.Test.Output))
		mxH = math.Max(mxH, getLineCounts(testcase.Output))
		rowHeights = append(rowHeights, int(mxH))
	}
	return rowHeights
}

func getLineCounts(s string) float64 {
	lines := strings.Split(s, "\n")
	return float64(len(lines))
}

func makeTable(headings []string, rows [][]string) *fyne.Container {

	columns := rowsToColumns(headings, rows)

	var objects []fyne.CanvasObject
	objects = append(objects, widget.NewSeparator())
	for k, col := range columns {
		box := container.NewVBox(widget.NewSeparator())
		box.Add(widget.NewLabelWithStyle(headings[k], fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
		box.Add(widget.NewSeparator())
		for _, val := range col {
			box.Add(widget.NewLabel(val))
			box.Add(widget.NewSeparator())
		}

		objects = append(objects, box)
		objects = append(objects, widget.NewSeparator())
	}
	return container.NewHBox(objects...)
}

func rowsToColumns(headings []string, rows [][]string) [][]string {
	columns := make([][]string, len(headings))
	for _, row := range rows {
		for colK := range row {
			columns[colK] = append(columns[colK], row[colK])
		}
	}
	return columns
}
