package testcaseGeneratorServices

import (
	"fmt"
	"github.com/skmonir/mango-gui/backend/judge-framework/dto"
	"github.com/skmonir/mango-gui/backend/judge-framework/executor"
	"github.com/skmonir/mango-gui/backend/judge-framework/models"
	"github.com/skmonir/mango-gui/backend/judge-framework/services"
	"github.com/skmonir/mango-gui/backend/judge-framework/services/languageServices"
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
	scriptBinaryPath := languageServices.GetBinaryFilePathByFilePath(request.GeneratorScriptPath)
	executionCommand := languageServices.GetExecutionCommandByFilePath(request.GeneratorScriptPath)
	if !utils.IsFileExist(scriptBinaryPath) {
		execResult.CompilationError = "Solution script binary not found!"
		return execResult
	}

	if !utils.IsDirExist(request.OutputDirectoryPath) {
		_ = utils.CreateDir(request.OutputDirectoryPath)
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
				ExecutionCommand:   executionCommand,
			},
		}
		execResult.TestcaseExecutionDetailsList = append(execResult.TestcaseExecutionDetailsList, execDetail)
	}

	execResult = executor.Execute(execResult, "output_generate_result_event", true)

	services.UpdateProblemExecutionResultInCacheByUrl(request.ParsedProblemUrl)

	return execResult
}
