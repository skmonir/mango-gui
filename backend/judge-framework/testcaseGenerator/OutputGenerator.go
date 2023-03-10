package testcaseGenerator

import (
	"fmt"
	"github.com/skmonir/mango-gui/backend/judge-framework/dto"
	"github.com/skmonir/mango-gui/backend/judge-framework/executor"
	"github.com/skmonir/mango-gui/backend/judge-framework/models"
	"github.com/skmonir/mango-gui/backend/judge-framework/services"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
	"path/filepath"
	"strings"
)

func GenerateOutput(request dto.TestcaseGenerateRequest) dto.ProblemExecutionResult {
	execResult := dto.ProblemExecutionResult{
		TestcaseExecutionDetailsList: []dto.TestcaseExecutionDetails{},
	}

	// Step-1: Compile generator source
	if err := compileScript(request.GeneratorScriptPath, false); err != "" {
		fmt.Println(err)
		execResult.CompilationError = err
		return execResult
	}

	// Step-2: Check if generator binary is created
	scriptBinaryPathWithoutExt := strings.TrimSuffix(request.GeneratorScriptPath, filepath.Ext(request.GeneratorScriptPath))
	if !utils.IsFileExist(scriptBinaryPathWithoutExt) {
		execResult.CompilationError = "Solution script binary not found!"
		return execResult
	}

	// Step-3: Prepare testcase detail and run the executor
	inputFiles := utils.GetFileNamesInDirectory(request.InputDirectoryPath)
	for _, inputFilename := range inputFiles {
		inputFilepath := filepath.Join(request.InputDirectoryPath, inputFilename)
		outputFilepath := filepath.Join(request.OutputDirectoryPath, strings.Replace(inputFilename, "in", "out", -1))

		execDetail := dto.TestcaseExecutionDetails{
			Status: "running",
			Testcase: models.Testcase{
				TimeLimit:          5,
				MemoryLimit:        512,
				InputFilePath:      inputFilepath,
				ExecOutputFilePath: outputFilepath,
				ExecutionCommand:   []string{scriptBinaryPathWithoutExt},
			},
		}
		execResult.TestcaseExecutionDetailsList = append(execResult.TestcaseExecutionDetailsList, execDetail)
	}

	execResult = executor.Execute(execResult, "output_generate_result_event")

	if len(request.ProblemUrl) > 0 {
		fmt.Println("Updating cache after output generation")
		ps := services.GetProblemListByUrl(request.ProblemUrl)
		if len(ps) > 0 {
			services.GetProblemExecutionResult(ps[0].Platform, ps[0].ContestId, ps[0].Label, true, true)
		} else {
			fmt.Println("No parsed problem found for", request.ProblemUrl)
		}
	}

	return execResult
}
