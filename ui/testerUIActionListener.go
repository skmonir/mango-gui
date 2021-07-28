package ui

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/skmonir/mango-gui/system"
	"github.com/skmonir/mango-gui/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/skmonir/mango-gui/context"
	"github.com/skmonir/mango-gui/models"
)

func ValidateTesterUiFields(ctx *context.AppCtx) error {
	if ctx.TesterUi.OnlineJudgeOptionSelect.Selected == "" {
		return errors.New("please select Online Judge before continuing")
	} else if ctx.TesterUi.ContestIdInputField.Text == "" {
		return errors.New("please enter valid Contest ID before continuing")
	}
	return nil
}

func ProcessProblems(ctx *context.AppCtx, problems []models.Problem) {
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

func RunTestClicked(ctx *context.AppCtx) error {
	problemName := ctx.TesterUi.ProblemNameListSelect.Selected
	problemId := strings.Split(problemName, ".")[0]

	testResults, err := system.RunTest(ctx, problemId)
	if err != nil {
		return err
	}

	(*ctx.TesterUi.TestcaseList)[problemName] = testResults

	return nil
}

func GetNewTable(testcaseList []models.ExecutionResult) *fyne.Container {
	rowHeights := GetRowHeights(testcaseList)

	rows := make([][]string, len(testcaseList))

	for i, testcase := range testcaseList {
		mxH := rowHeights[i]
		rows[i] = append(rows[i], GetHeightAdjustedCell(mxH, strconv.Itoa(i)))
		rows[i] = append(rows[i], GetHeightAdjustedCell(mxH, testcase.Test.Input))
		rows[i] = append(rows[i], GetHeightAdjustedCell(mxH, testcase.Test.Output))
		rows[i] = append(rows[i], GetHeightAdjustedCell(mxH, testcase.Output))
		rows[i] = append(rows[i], GetHeightAdjustedCell(mxH, testcase.Verdict))
		rows[i] = append(rows[i], GetHeightAdjustedCell(mxH, fmt.Sprintf("%v ms", testcase.Runtime)))
		rows[i] = append(rows[i], GetHeightAdjustedCell(mxH, utils.ParseMemoryInKb(testcase.Memory)))
	}

	headings := []string{"#", "SAMPLE INPUT", "SAMPLE OUTPUT", "PROGRAM OUTPUT", "VERDICT", "TIME", "MEMORY"}

	return MakeTable(headings, rows)
}

func GetHeightAdjustedCell(height int, value string) string {
	for len(strings.Split(value, "\n")) < height {
		value += "\n "
	}
	return value
}

func GetRowHeights(testcaseList []models.ExecutionResult) []int {
	var rowHeights []int
	for _, testcase := range testcaseList {
		mxH := math.Max(GetLineCounts(testcase.Test.Input), GetLineCounts(testcase.Test.Output))
		mxH = math.Max(mxH, GetLineCounts(testcase.Output))
		rowHeights = append(rowHeights, int(mxH))
	}
	return rowHeights
}

func GetLineCounts(s string) float64 {
	lines := strings.Split(s, "\n")
	return float64(len(lines))
}

func MakeTable(headings []string, rows [][]string) *fyne.Container {

	columns := RowsToColumns(headings, rows)

	var objects []fyne.CanvasObject
	objects = append(objects, widget.NewSeparator())
	for k, col := range columns {
		box := container.NewVBox(widget.NewSeparator())
		box.Add(widget.NewLabelWithStyle(headings[k], fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
		box.Add(widget.NewSeparator())

		textAlign := fyne.TextAlignCenter
		if 0 < k && k < 4 {
			textAlign = fyne.TextAlignLeading
		}
		for _, val := range col {
			box.Add(widget.NewLabelWithStyle(val, textAlign, fyne.TextStyle{}))
			box.Add(widget.NewSeparator())
		}

		objects = append(objects, box)
		objects = append(objects, widget.NewSeparator())
	}
	return container.NewHBox(objects...)
}

func RowsToColumns(headings []string, rows [][]string) [][]string {
	columns := make([][]string, len(headings))
	for _, row := range rows {
		for colK := range row {
			columns[colK] = append(columns[colK], row[colK])
		}
	}
	return columns
}
