package services

import (
	"fmt"
	"time"

	"github.com/skmonir/mango-ui/backend/judge-framework/cacheServices"
	"github.com/skmonir/mango-ui/backend/judge-framework/fileService"
	"github.com/skmonir/mango-ui/backend/judge-framework/utils"

	"github.com/skmonir/mango-ui/backend/judge-framework/dto"
	"github.com/skmonir/mango-ui/backend/judge-framework/executor"
	"github.com/skmonir/mango-ui/backend/socket"
)

func RunTest(platform string, cid string, label string) dto.ProblemExecutionResult {
	execResult := GetProblemExecutionResult(platform, cid, label, true)

	// Step 1: Compile the source
	socket.PublishStatusMessage("test_status", "Compiling source code", "info")
	err := executor.Compile(platform, cid, label)
	if err != "" {
		socket.PublishStatusMessage("test_status", "Compilation error!", "error")
		execResult.CompilationError = err
		return execResult
	}
	socket.PublishStatusMessage("test_status", "Compilation successful!", "success")

	time.Sleep(500 * time.Millisecond)

	// Step 2: Check if binary is available for the source
	binaryPath := fileService.GetSourceBinaryPath(platform, cid, label)
	if !utils.IsFileExist(binaryPath) {
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
	socket.PublishExecutionResult(execResult)

	// Step 4: Run the binary and check testcases
	socket.PublishStatusMessage("test_status", "Running testcases", "info")
	execResult = executor.Execute(execResult)
	cacheServices.SaveExecutionResult(platform, cid, label, execResult)

	totalPassed, totalTests := 0, len(execResult.TestcaseExecutionDetailsList)
	for _, execDetails := range execResult.TestcaseExecutionDetailsList {
		if execDetails.TestcaseExecutionResult.Verdict == "OK" {
			totalPassed++
		}
	}
	testStatus := fmt.Sprintf("%v/%v Tests Passed", totalPassed, totalTests)
	if totalPassed == totalTests {
		socket.PublishStatusMessage("test_status", testStatus, "success")
	} else {
		socket.PublishStatusMessage("test_status", testStatus, "error")
	}

	return execResult
}

func GetProblemExecutionResult(platform string, cid string, label string, isForUI bool) dto.ProblemExecutionResult {
	fmt.Println("Fetching execution result for", platform, cid, label)

	maxRow, maxCol := 10000000, 10000000
	if isForUI {
		if ok, execResult := cacheServices.GetExecutionResult(platform, cid, label); ok {
			return execResult
		}
		maxRow, maxCol = 30, 50
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
