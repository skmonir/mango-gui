package fileService

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/skmonir/mango-ui/backend/judge-framework/config"
	"github.com/skmonir/mango-ui/backend/judge-framework/models"
	"github.com/skmonir/mango-ui/backend/judge-framework/utils"
)

func SaveTestcasesIntoFiles(platform string, cid string, label string, testcases []models.Testcase) {
	inputDirectory, outputDirectory := getInputOutputDirectories(platform, cid, label)
	for id, testcase := range testcases {
		inputFilename := fmt.Sprintf("00_sample_input_%03d.txt", id)
		outputFilename := fmt.Sprintf("00_sample_output_%03d.txt", id)
		go utils.WriteFileContent(inputDirectory, inputFilename, []byte(testcase.Input))
		go utils.WriteFileContent(outputDirectory, outputFilename, []byte(testcase.Output))
	}
}

func GetTestcasesFromFile(platform string, cid string, label string, maxRow int, maxCol int) []models.Testcase {
	fmt.Println("Fetching testcases from file...")

	inputPaths, outputPaths := GetInputOutputFilePaths(platform, cid, label)
	sourceBinaryPath := GetSourceBinaryPath(platform, cid, label)

	var testcases []models.Testcase
	for i := 0; i < len(inputPaths); i++ {
		input := utils.ReadFileContent(inputPaths[i], maxRow, maxCol)
		output := utils.ReadFileContent(outputPaths[i], maxRow, maxCol)
		testcases = append(testcases, models.Testcase{
			Input:            input,
			Output:           output,
			InputFilePath:    inputPaths[i],
			OutputFilePath:   outputPaths[i],
			SourceBinaryPath: sourceBinaryPath,
		})
	}
	fmt.Println("Fetched testcases from file.")
	return testcases
}

func GetInputOutputFilePaths(platform string, cid string, label string) ([]string, []string) {
	inputDirectory, outputDirectory := getInputOutputDirectories(platform, cid, label)

	var inputPaths []string
	var outputPaths []string
	inputFiles := utils.GetFileNamesInDirectory(inputDirectory)
	sort.Strings(inputFiles)
	for i := 0; i < len(inputFiles); i++ {
		inputFilepath := filepath.Join(inputDirectory, inputFiles[i])
		outputFilepath := filepath.Join(outputDirectory, strings.Replace(inputFiles[i], "in", "out", -1))

		inputPaths = append(inputPaths, inputFilepath)
		outputPaths = append(outputPaths, outputFilepath)
	}

	return inputPaths, outputPaths
}

func getInputOutputDirectories(platform string, cid string, label string) (string, string) {
	conf := config.GetJudgeConfigFromCache()
	inputDirectory := filepath.Join(conf.WorkspaceDirectory, platform, cid, "input", label)
	outputDirectory := filepath.Join(conf.WorkspaceDirectory, platform, cid, "output", label)
	return inputDirectory, outputDirectory
}
