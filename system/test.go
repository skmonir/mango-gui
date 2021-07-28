package system

import (
	"errors"
	"strings"

	"github.com/skmonir/mango-gui/context"
	"github.com/skmonir/mango-gui/models"
	"github.com/skmonir/mango-gui/utils"
)

func GetVerdict(testcase models.Testcase, executionResult *models.ExecutionResult) {
	errorMsg := ""
	if !executionResult.Status && executionResult.Error != nil {
		errorMsg = executionResult.Error.Error()
	}

	if strings.Contains(errorMsg, "segmentation fault") {
		executionResult.Verdict = "RE"
	} else if (executionResult.Runtime > (testcase.TimeLimit * 1000)) || strings.Contains(errorMsg, "killed") {
		executionResult.Verdict = "TLE"
		executionResult.Runtime = testcase.TimeLimit * 1000
	} else if utils.ConvertMemoryInMb(executionResult.Memory) > testcase.MemoryLimit {
		executionResult.Verdict = "MLE"
	} else if !executionResult.Status {
		executionResult.Verdict = "RE"
	} else if testcase.Output == executionResult.Output {
		executionResult.Verdict = "OK"
	} else {
		executionResult.Verdict = "WA"
	}
}

func RunTest(ctx *context.AppCtx, problemId string) ([]models.ExecutionResult, error) {
	ctx.ProgressBar.SetValue(0)

	var testResults []models.ExecutionResult

	if ctx.TesterUi.CurrentContestId == "" {
		return testResults, errors.New("contest id not found")
	}

	if problemId == "" {
		return testResults, errors.New("problem id not found")
	}

	ctx.Config.CurrentContestId = ctx.TesterUi.CurrentContestId
	ctx.Config.SaveConfig()

	problemInfo, err := ctx.Config.GetProblemInfo(problemId)
	if err != nil {
		return testResults, err
	}

	if err := CompileSource(*ctx.Config, problemId); err != nil {
		return testResults, err
	}

	testcases := problemInfo.Dataset

	if len(testcases) > 0 {
		ctx.ProgressBar.Max = float64(len(testcases))
	}

	for i, testcase := range testcases {
		executionResult := ExecuteSourceBinary(*ctx.Config, testcase, problemId)
		GetVerdict(testcase, &executionResult)
		testResults = append(testResults, executionResult)

		ctx.ProgressBar.SetValue(float64(i + 1))
	}

	return testResults, nil
}
