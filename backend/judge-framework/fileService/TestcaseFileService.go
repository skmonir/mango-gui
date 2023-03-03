package fileService

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/skmonir/mango-ui/backend/judge-framework/cache"
	"github.com/skmonir/mango-ui/backend/judge-framework/config"
	"github.com/skmonir/mango-ui/backend/judge-framework/dto"
	"github.com/skmonir/mango-ui/backend/judge-framework/models"
	"github.com/skmonir/mango-ui/backend/judge-framework/utils"
)

func GetProblemExecutionResult(platform string, cid string, label string) dto.ProblemExecutionResult {
	key := fmt.Sprintf("ProblemExecutionResult:%v.%v.%v", platform, cid, label)

	globalCache := cache.GetGlobalCache()
	execResStr := globalCache.Get(key)

	var probExecRes dto.ProblemExecutionResult
	if execResStr != "" {
		if err := json.Unmarshal([]byte(execResStr), &probExecRes); err != nil {
			fmt.Println(err)
		} else {
			return probExecRes
		}
	}

	return probExecRes
}

func SaveTestcasesIntoFiles(platform string, cid string, label string, testcases []models.Testcase) {
	inputDirectory, outputDirectory, userOutputDirectory := getInputOutputDirectories(platform, cid, label)
	for id, testcase := range testcases {
		inputFilename := fmt.Sprintf("sample%v.txt", id+1)
		outputFilename := fmt.Sprintf("sample%v.txt", id+1)
		userOutputFilename := fmt.Sprintf("sample%v.txt", id+1)
		go utils.WriteFileContent(inputDirectory, inputFilename, []byte(testcase.Input))
		go utils.WriteFileContent(outputDirectory, outputFilename, []byte(testcase.Output))
		go utils.CreateFile(userOutputDirectory, userOutputFilename)
	}
}

func GetTestcasesFromFile(platform string, cid string, label string) []models.Testcase {
	fmt.Println("Fetching testcases from file...")

	inputpaths, outputpaths, userOutputpaths := GetInputOutputFilePaths(platform, cid, label)
	sourceBinaryPath := GetSourceBinaryPath(platform, cid, label)

	var testcases []models.Testcase
	for i := 0; i < len(inputpaths); i++ {
		input := utils.ReadFileContent(inputpaths[i], 30, 50)
		output := utils.ReadFileContent(outputpaths[i], 30, 50)
		testcases = append(testcases, models.Testcase{
			Input:              input,
			Output:             output,
			InputFilePath:      inputpaths[i],
			OutputFilePath:     outputpaths[i],
			UserOutputFilePath: userOutputpaths[i],
			SourceBinaryPath:   sourceBinaryPath,
		})
	}
	fmt.Println("Fetched testcases from file.")
	return testcases
}

func GetInputOutputFilePaths(platform string, cid string, label string) ([]string, []string, []string) {
	inputDirectory, outputDirectory, userOutputDirectory := getInputOutputDirectories(platform, cid, label)
	inputFiles := utils.GetFileNamesInDirectory(inputDirectory)

	var inputpaths []string
	var outputpaths []string
	var userOutputpaths []string
	for i := 0; i < len(inputFiles); i++ {
		inputFilepath := filepath.Join(inputDirectory, inputFiles[i])
		outputFilepath := filepath.Join(outputDirectory, inputFiles[i])
		userOutputFilepath := filepath.Join(userOutputDirectory, inputFiles[i])

		inputpaths = append(inputpaths, inputFilepath)
		outputpaths = append(outputpaths, outputFilepath)
		userOutputpaths = append(userOutputpaths, userOutputFilepath)
	}

	return inputpaths, outputpaths, userOutputpaths
}

func getInputOutputDirectories(platform string, cid string, label string) (string, string, string) {
	conf := config.GetJudgeConfigFromCache()
	inputDirectory := filepath.Join(conf.WorkspaceDirectory, platform, cid, "input", label)
	outputDirectory := filepath.Join(conf.WorkspaceDirectory, platform, cid, "output", label)
	userOutputDirectory := filepath.Join(conf.WorkspaceDirectory, platform, cid, "user_output", label)
	return inputDirectory, outputDirectory, userOutputDirectory
}
