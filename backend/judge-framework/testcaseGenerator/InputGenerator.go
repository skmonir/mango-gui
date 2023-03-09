package testcaseGenerator

import (
	"fmt"
	"github.com/skmonir/mango-gui/backend/judge-framework/config"
	"github.com/skmonir/mango-gui/backend/judge-framework/dto"
	"github.com/skmonir/mango-gui/backend/judge-framework/executor"
	"github.com/skmonir/mango-gui/backend/judge-framework/models"
	"github.com/skmonir/mango-gui/backend/judge-framework/services"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func Generate(request dto.TestcaseGenerateRequest) dto.ProblemExecutionResult {
	var execResult dto.ProblemExecutionResult

	if request.GenerationProcess == "" {
		execResult = generateWithTgenScript(request)
	} else {
		execResult = runGenerator(request, false)
	}

	return execResult
}

func generateWithTgenScript(request dto.TestcaseGenerateRequest) dto.ProblemExecutionResult {
	if msg := runValidator(request.TgenScriptContent); msg != "OK" {
		return dto.ProblemExecutionResult{
			CompilationError: msg,
		}
	}

	folderPath := getScriptDirectory()
	request.GeneratorScriptPath = folderPath + "generator.cpp"

	return runGenerator(request, true)
}

func runValidator(tgenScript string) string {
	folderPath := getScriptDirectory()

	// Step-1: Compile validator source
	if err := compileScript(folderPath+"validator.cpp", true); err != "" {
		fmt.Println(err)
		return err
	}

	// Step-2: Validate the script
	if err := validateScript(tgenScript); err != "" {
		fmt.Println(err)
		return err
	}

	return ""
}

func validateScript(tgenScript string) string {
	folderPath := getScriptDirectory()
	filePathWithoutExt := folderPath + "validator"

	execResult := dto.ProblemExecutionResult{
		TestcaseExecutionDetailsList: []dto.TestcaseExecutionDetails{
			{
				Status: "none",
				Testcase: models.Testcase{
					Input:            utils.TrimIO(tgenScript + "\nEND"),
					TimeLimit:        5,
					MemoryLimit:      512,
					ExecutionCommand: []string{filePathWithoutExt},
				},
			},
		},
	}

	execResult = executor.Execute(execResult, "input_generate_result_event")

	return execResult.TestcaseExecutionDetailsList[0].TestcaseExecutionResult.Output
}

func runGenerator(request dto.TestcaseGenerateRequest, skipIfCompiled bool) dto.ProblemExecutionResult {
	// Step-1: Compile generator source
	if err := compileScript(request.GeneratorScriptPath, skipIfCompiled); err != "" {
		fmt.Println(err)
		return dto.ProblemExecutionResult{
			CompilationError: err,
		}
	}

	// Step-2: Validate the script
	execResult := generateInput(request)

	if len(request.ProblemUrl) > 0 {
		ps := services.GetProblemListByUrl(request.ProblemUrl)
		if len(ps) > 0 {
			services.GetProblemExecutionResult(ps[0].Platform, ps[0].ContestId, ps[0].Label, true, true)
		}
	}

	return execResult
}

func generateInput(request dto.TestcaseGenerateRequest) dto.ProblemExecutionResult {
	filePathWithoutExt := strings.TrimSuffix(request.GeneratorScriptPath, filepath.Ext(request.GeneratorScriptPath))

	execResult := dto.ProblemExecutionResult{
		TestcaseExecutionDetailsList: []dto.TestcaseExecutionDetails{},
	}

	paramId := time.Now().UnixNano() / int64(time.Millisecond)
	sn := request.SerialFrom - 1

	for i := 0; i < request.FileNum; i++ {
		sn++
		paramId++
		execOutputFilePath := fmt.Sprintf("%v_%03d.txt", filepath.Join(request.InputDirectoryPath, request.FileName), sn)
		execDetail := dto.TestcaseExecutionDetails{
			Status: "running",
			Testcase: models.Testcase{
				Input:              utils.TrimIO(request.TgenScriptContent + "\nEND"),
				TimeLimit:          5,
				MemoryLimit:        512,
				ExecOutputFilePath: execOutputFilePath,
				ExecutionCommand:   []string{filePathWithoutExt, strconv.Itoa(request.TestPerFile), strconv.FormatInt(paramId, 10)},
			},
		}
		execResult.TestcaseExecutionDetailsList = append(execResult.TestcaseExecutionDetailsList, execDetail)
	}

	execResult = executor.Execute(execResult, "input_generate_result_event")

	return execResult
}

func compileScript(filePathWithExt string, skipIfCompiled bool) string {
	fmt.Println("Compiling " + filePathWithExt)
	judgeConfig := config.GetJudgeConfigFromCache()

	filePathWithoutExt := strings.TrimSuffix(filePathWithExt, filepath.Ext(filePathWithExt))

	if skipIfCompiled && utils.IsFileExist(filePathWithoutExt) {
		return ""
	}
	if !utils.IsFileExist(filePathWithExt) {
		return filePathWithExt + ": file not found!"
	}

	command := fmt.Sprintf("%v %v %v -o %v", judgeConfig.ActiveLanguage.CompilationCommand, judgeConfig.ActiveLanguage.CompilationArgs, filePathWithExt, filePathWithoutExt)

	return executor.CompileSource(command, false)
}

func getScriptDirectory() string {
	projectBasePath, _ := os.Getwd()
	return projectBasePath + "/backend/judge-framework/testcaseGenerator/scripts/"
}
