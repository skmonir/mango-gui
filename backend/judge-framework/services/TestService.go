package services

import (
	"fmt"
	"github.com/skmonir/mango-gui/backend/judge-framework/constants"
	"github.com/skmonir/mango-gui/backend/judge-framework/logger"
	"time"

	"github.com/skmonir/mango-gui/backend/judge-framework/cacheServices"
	"github.com/skmonir/mango-gui/backend/judge-framework/fileService"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"

	"github.com/skmonir/mango-gui/backend/judge-framework/dto"
	"github.com/skmonir/mango-gui/backend/judge-framework/executor"
	"github.com/skmonir/mango-gui/backend/socket"
)

func RunTest(platform string, cid string, label string) dto.ProblemExecutionResult {
	logger.Info(fmt.Sprintf("Run test request received for %v %v %v", platform, cid, label))
	// Step 0: Fetch new/previous execution result object and reset necessary fields
	execResult := GetProblemExecutionResult(platform, cid, label, true, false)
	execResult.CompilationError = ""
	for i := 0; i < len(execResult.TestcaseExecutionDetailsList); i++ {
		execResult.TestcaseExecutionDetailsList[i].Status = "none"
		execResult.TestcaseExecutionDetailsList[i].TestcaseExecutionResult = dto.TestcaseExecutionResult{}
	}

	// Step 1: Compile the source
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
	binaryPath := fileService.GetSourceBinaryPath(platform, cid, label)
	if !utils.IsFileExist(binaryPath) {
		logger.Error(fmt.Sprintf("Binary file not found for %v", binaryPath))
		socket.PublishStatusMessage("test_status", "Binary file not found!", "error")
		return execResult
	}

	// Step 3: Prepare testcases for execution
	prob := GetProblem(platform, cid, label)
	for i := 0; i < len(execResult.TestcaseExecutionDetailsList); i++ {
		execResult.TestcaseExecutionDetailsList[i].Status = "running"
		execResult.TestcaseExecutionDetailsList[i].Testcase.TimeLimit = prob.TimeLimit
		execResult.TestcaseExecutionDetailsList[i].Testcase.MemoryLimit = prob.MemoryLimit
	}
	socket.PublishExecutionResult(execResult, "test_exec_result_event")

	// Step 4: Run the binary and check testcases
	//socket.PublishStatusMessage("test_status", "Running testcases", "info")
	execResult = executor.Execute(execResult, "test_exec_result_event")
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

	if !isSkipCache {
		if ok, execResult := cacheServices.GetExecutionResult(platform, cid, label); ok {
			return execResult
		}
	}

	testcases := fileService.GetTestcasesFromFile(platform, cid, label, maxRow, maxCol)
	var testcaseExecutionDetailsList []dto.TestcaseExecutionDetails
	for i := 0; i < len(testcases); i++ {
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
