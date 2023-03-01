package fileService

import (
	"encoding/json"
	"fmt"
	"github.com/skmonir/mango-ui/backend/judge-framework/cache"
	"github.com/skmonir/mango-ui/backend/judge-framework/config"
	"github.com/skmonir/mango-ui/backend/judge-framework/dto"
	"github.com/skmonir/mango-ui/backend/judge-framework/models"
	"github.com/skmonir/mango-ui/backend/judge-framework/utils"
	"path/filepath"
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
	inputDirectory, outputDirectory := getInputOutputDirectories(platform, cid, label)
	for id, testcase := range testcases {
		inputFilename := fmt.Sprintf("sample%v.txt", id+1)
		outputFilename := fmt.Sprintf("sample%v.txt", id+1)
		go utils.WriteFileContent(inputDirectory, inputFilename, []byte(testcase.Input))
		go utils.WriteFileContent(outputDirectory, outputFilename, []byte(testcase.Output))
	}
}

func GetTestcasesFromFile(platform string, cid string, label string) {
	inputDirectory, outputDirectory := getInputOutputDirectories(platform, cid, label)
	for _, filename := range utils.GetFileNamesInDirectory(inputDirectory) {
		fmt.Println(filename)
	}
	for _, filename := range utils.GetFileNamesInDirectory(outputDirectory) {
		fmt.Println(filename)
	}
}

func getInputOutputDirectories(platform string, cid string, label string) (string, string) {
	conf := config.GetJudgeConfigFromCache()
	inputDirectory := filepath.Join(conf.WorkspaceDirectory, platform, cid, "input", label)
	outputDirectory := filepath.Join(conf.WorkspaceDirectory, platform, cid, "output", label)
	return inputDirectory, outputDirectory
}
