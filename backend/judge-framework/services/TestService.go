package services

import (
	"fmt"
	"github.com/skmonir/mango-gui/backend/judge-framework/constants"
	"github.com/skmonir/mango-gui/backend/judge-framework/logger"
	"github.com/skmonir/mango-gui/backend/judge-framework/services/cacheServices"
	"github.com/skmonir/mango-gui/backend/judge-framework/services/fileServices"
	"github.com/skmonir/mango-gui/backend/judge-framework/services/languageServices"
	"time"

	"github.com/skmonir/mango-gui/backend/judge-framework/utils"

	"github.com/skmonir/mango-gui/backend/judge-framework/dto"
	"github.com/skmonir/mango-gui/backend/judge-framework/executor"
	"github.com/skmonir/mango-gui/backend/socket"
)

func RunTest(platform string, cid string, label string) dto.ProblemExecutionResult {
	logger.Info(fmt.Sprintf("Run test request received for %v %v %v", platform, cid, label))
	// Step 0: Fetch new/previous execution result object and reset necessary fields
	execResult := GetProblemExecutionResult(platform, cid, label, true, false)
	if len(execResult.TestcaseExecutionDetailsList) == 0 {
		logger.Error("No input files to test!")
		socket.PublishStatusMessage("test_status", "No input files to test!", "error")
		return execResult
	}
	execResult.CompilationError = ""
	for i := 0; i < len(execResult.TestcaseExecutionDetailsList); i++ {
		execResult.TestcaseExecutionDetailsList[i].Status = "none"
		execResult.TestcaseExecutionDetailsList[i].TestcaseExecutionResult = dto.TestcaseExecutionResult{}
	}

	// Step 1: Compile the source (except Python)
	socket.PublishStatusMessage("test_status", "Compiling source code", "info")
	err := executor.Compile(platform, cid, label)
	if err != "" {
		logger.Info("Compilation error!")
		socket.PublishStatusMessage("test_status", "Compilation error!", "error")
		execResult.CompilationError = err
		cacheServices.SaveExecutionResult(platform, cid, label, execResult)
		return execResult
	}
	logger.Info("Compilation successful!")
	socket.PublishStatusMessage("test_status", "Compilation successful!", "success")
	socket.PublishExecutionResult(execResult, "test_exec_result_event")

	time.Sleep(500 * time.Millisecond)

	// Step 2: Check if binary is available for the source
	binaryPath := GetProblemSourceBinaryPath(platform, cid, label)
	if !utils.IsFileExist(binaryPath) {
		logger.Error(fmt.Sprintf("Binary file not found for %v", binaryPath))
		socket.PublishStatusMessage("test_status", "Binary file not found!", "error")
		return execResult
	}

	// Step 3: Prepare testcases for execution
	prob := GetProblem(platform, cid, label)
	for i := 0; i < len(execResult.TestcaseExecutionDetailsList); i++ {
		execResult.TestcaseExecutionDetailsList[i].Testcase.ExecutionCommand = GetProblemBinaryExecCommand(platform, cid, label)
		execResult.TestcaseExecutionDetailsList[i].Status = "running"
		execResult.TestcaseExecutionDetailsList[i].Testcase.TimeLimit = prob.TimeLimit
		execResult.TestcaseExecutionDetailsList[i].Testcase.MemoryLimit = prob.MemoryLimit
	}
	socket.PublishExecutionResult(execResult, "test_exec_result_event")

	// Step 4: Run the binary and check testcases
	execResult = executor.Execute(execResult, "test_exec_result_event", false)
	cacheServices.SaveExecutionResult(platform, cid, label, execResult)

	logger.Info("Execution complete")
	return execResult
}

func GetProblemExecutionResult(platform string, cid string, label string, isForUI bool, isSkipCache bool) dto.ProblemExecutionResult {
	logger.Info(fmt.Sprintf("Fetching execution result for %v %v %v %v %v", platform, cid, label, isForUI, isSkipCache))

	maxRow, maxCol := constants.IO_MAX_ROW_FOR_TEST, constants.IO_MAX_COL_FOR_TEST
	if isForUI {
		maxRow, maxCol = constants.IO_MAX_ROW_FOR_UI, constants.IO_MAX_COL_FOR_UI
	}

	sourceBinaryPath := GetProblemSourceBinaryPath(platform, cid, label)

	if !isSkipCache {
		if ok, execResult := cacheServices.GetExecutionResult(platform, cid, label); ok {
			for i := 0; i < len(execResult.TestcaseExecutionDetailsList); i++ {
				execResult.TestcaseExecutionDetailsList[i].Testcase.SourceBinaryPath = sourceBinaryPath
			}
			return execResult
		}
	}

	testcases := fileServices.GetTestcasesFromFile(platform, cid, label, maxRow, maxCol)
	var testcaseExecutionDetailsList []dto.TestcaseExecutionDetails
	for i := 0; i < len(testcases); i++ {
		testcases[i].SourceBinaryPath = sourceBinaryPath
		testcases[i].ExecutionCommand = GetProblemBinaryExecCommand(platform, cid, label)
		testcaseExecutionDetailsList = append(testcaseExecutionDetailsList, dto.TestcaseExecutionDetails{
			Status:   "none",
			Testcase: testcases[i],
		})
	}
	execResult := dto.ProblemExecutionResult{
		CompilationError:             "",
		TestcaseExecutionDetailsList: testcaseExecutionDetailsList,
	}

	cacheServices.SaveExecutionResult(platform, cid, label, execResult)

	return execResult
}

func UpdateProblemExecutionResultInCacheByUrl(url string) {
	if len(url) > 0 {
		fmt.Println("Updating cache after output generation")
		ps := GetProblemListByUrl(url)
		if len(ps) > 0 {
			GetProblemExecutionResult(ps[0].Platform, ps[0].ContestId, ps[0].Label, true, true)
		} else {
			fmt.Println("No parsed problem found for", url)
		}
	}
}

func GetProblemSourceBinaryPath(platform string, cid string, label string) string {
	filePath := fileServices.GetSourceFilePath(platform, cid, label)
	return languageServices.GetBinaryFilePathByFilePath(filePath)
}

func GetProblemBinaryExecCommand(platform string, cid string, label string) []string {
	filePath := fileServices.GetSourceFilePath(platform, cid, label)
	return languageServices.GetExecutionCommandByFilePath(filePath)
}
