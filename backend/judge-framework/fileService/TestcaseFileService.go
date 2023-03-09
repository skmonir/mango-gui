package fileService

import (
	"fmt"
	"github.com/skmonir/mango-gui/backend/judge-framework/constants"
	"path/filepath"
	"strings"

	"github.com/skmonir/mango-gui/backend/judge-framework/config"
	"github.com/skmonir/mango-gui/backend/judge-framework/models"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
)

func GetTestcaseByPath(inputFilePath string, outputFilePath string) models.Testcase {
	input := utils.ReadFileContent(inputFilePath, constants.IO_MAX_ROW_FOR_TEST, constants.IO_MAX_COL_FOR_TEST)
	output := utils.ReadFileContent(outputFilePath, constants.IO_MAX_ROW_FOR_TEST, constants.IO_MAX_COL_FOR_TEST)
	return models.Testcase{
		Input:  input,
		Output: output,
	}
}

func SaveCustomTestcaseIntoFile(inputDirectory, outputDirectory, input, output string, maxCustomTestId int) {
	inputFilename := fmt.Sprintf("01_custom_input_%03d.txt", maxCustomTestId)
	outputFilename := fmt.Sprintf("01_custom_output_%03d.txt", maxCustomTestId)
	utils.WriteFileContent(inputDirectory, inputFilename, []byte(input))
	utils.WriteFileContent(outputDirectory, outputFilename, []byte(output))
}

func UpdateCustomTestcaseIntoFile(inputFilePath string, outputFilePath string, input string, output string) {
	inputDirectory, inputFilename := filepath.Split(inputFilePath)
	outputDirectory, outputFilename := filepath.Split(outputFilePath)
	go utils.WriteFileContent(inputDirectory, inputFilename, []byte(input))
	go utils.WriteFileContent(outputDirectory, outputFilename, []byte(output))
}

func SaveTestcasesIntoFiles(platform string, cid string, label string, testcases []models.Testcase) {
	inputDirectory, outputDirectory := GetInputOutputDirectories(platform, cid, label)
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
			ExecutionCommand: []string{sourceBinaryPath},
		})
	}
	fmt.Println("Fetched testcases from file.")
	return testcases
}

func GetInputOutputFilePaths(platform string, cid string, label string) ([]string, []string) {
	inputDirectory, outputDirectory := GetInputOutputDirectories(platform, cid, label)

	var inputPaths []string
	var outputPaths []string
	inputFiles := utils.GetFileNamesInDirectory(inputDirectory)
	for _, inputFilename := range inputFiles {
		inputFilepath := filepath.Join(inputDirectory, inputFilename)
		outputFilepath := filepath.Join(outputDirectory, strings.Replace(inputFilename, "in", "out", -1))

		inputPaths = append(inputPaths, inputFilepath)
		outputPaths = append(outputPaths, outputFilepath)
	}

	return inputPaths, outputPaths
}

func GetInputOutputDirectories(platform string, cid string, label string) (string, string) {
	conf := config.GetJudgeConfigFromCache()
	inputDirectory := filepath.Join(conf.WorkspaceDirectory, platform, cid, "input", label)
	outputDirectory := filepath.Join(conf.WorkspaceDirectory, platform, cid, "output", label)
	return inputDirectory, outputDirectory
}

func getInputTypeFromId(id string) string {
	inputType := ""
	if id == "00" {
		inputType = "sample"
	} else if id == "01" {
		inputType = "custom"
	} else if id == "02" {
		inputType = "random"
	}
	return inputType
}
