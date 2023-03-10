package testcaseGenerator

import (
	"fmt"
	"github.com/skmonir/mango-gui/backend/judge-framework/config"
	"github.com/skmonir/mango-gui/backend/judge-framework/dto"
	"github.com/skmonir/mango-gui/backend/judge-framework/executor"
	"github.com/skmonir/mango-gui/backend/judge-framework/models"
	"github.com/skmonir/mango-gui/backend/judge-framework/services"
	"github.com/skmonir/mango-gui/backend/judge-framework/testcaseGenerator/tgenScripts"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func GenerateInput(request dto.TestcaseGenerateRequest) dto.ProblemExecutionResult {
	var execResult dto.ProblemExecutionResult

	tgenScripts.CreateIfScriptsAreNotAvailable()

	if request.GenerationProcess == "tgen_script" {
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
	request.GeneratorScriptPath = filepath.Join(folderPath, "generator.cpp")

	return runGenerator(request, false)
}

func runValidator(tgenScript string) string {
	folderPath := getScriptDirectory()

	// Step-1: Compile validator source
	if err := compileScript(filepath.Join(folderPath, "validator.cpp"), false); err != "" {
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
	filePathWithoutExt := filepath.Join(folderPath, "validator")

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

	if execResult.CompilationError == "" && len(request.ProblemUrl) > 0 {
		fmt.Println("Updating cache after input generation")
		ps := services.GetProblemListByUrl(request.ProblemUrl)
		if len(ps) > 0 {
			services.GetProblemExecutionResult(ps[0].Platform, ps[0].ContestId, ps[0].Label, true, true)
		} else {
			fmt.Println("No parsed problem found for", request.ProblemUrl)
		}
	}

	return execResult
}

func generateInput(request dto.TestcaseGenerateRequest) dto.ProblemExecutionResult {
	scriptBinaryPathWithoutExt := strings.TrimSuffix(request.GeneratorScriptPath, filepath.Ext(request.GeneratorScriptPath))
	if !utils.IsFileExist(scriptBinaryPathWithoutExt) {
		return dto.ProblemExecutionResult{
			CompilationError: "Generator script binary not found!",
		}
	}

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
				ExecutionCommand:   []string{scriptBinaryPathWithoutExt, strconv.Itoa(request.TestPerFile), strconv.FormatInt(paramId, 10)},
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
	appdataDirectory := utils.GetAppDataDirectoryPath()
	scriptDirectory := filepath.Join(appdataDirectory, "tgen_scripts")
	return scriptDirectory
}
