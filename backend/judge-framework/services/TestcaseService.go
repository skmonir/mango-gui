package services

import (
	"github.com/skmonir/mango-gui/backend/judge-framework/models"
	"github.com/skmonir/mango-gui/backend/judge-framework/services/fileServices"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
	"strconv"
	"strings"
)

func GetTestcaseByPath(inputFilePath string, outputFilePath string) models.Testcase {
	return fileServices.GetTestcaseByPath(inputFilePath, outputFilePath)
}

func SaveCustomTestcaseIntoFile(platform string, cid string, label string, input string, output string) {
	inputDirectory, outputDirectory := fileServices.GetInputOutputDirectories(platform, cid, label)
	inputFiles := utils.GetFileNamesInDirectory(inputDirectory)
	maxCustomTestId := -1
	for _, filename := range inputFiles {
		if strings.Contains(filename, "01_custom_input_") {
			serialStr := strings.Replace(filename, "01_custom_input_", "", -1)
			serialStr = strings.Replace(serialStr, ".txt", "", -1)
			if serial, err := strconv.Atoi(serialStr); err == nil {
				if serial > maxCustomTestId {
					maxCustomTestId = serial
				}
			}
		}
	}
	maxCustomTestId++

	fileServices.SaveCustomTestcaseIntoFile(inputDirectory, outputDirectory, input, output, maxCustomTestId)
	GetProblemExecutionResult(platform, cid, label, true, true)
}

func UpdateCustomTestcaseIntoFile(platform, cid, label, inputFilePath, outputFilePath, input, output string) {
	fileServices.UpdateCustomTestcaseIntoFile(inputFilePath, outputFilePath, input, output)
	GetProblemExecutionResult(platform, cid, label, true, true)
}

func DeleteCustomTestcaseFromFile(platform, cid, label, inputFilePath string) error {
	if err := utils.RemoveFile(inputFilePath); err != nil {
		return err
	}
	GetProblemExecutionResult(platform, cid, label, true, true)
	return nil
}

func GetInputOutputDirectoryByUrl(url string) (string, string) {
	probs := GetProblemListByUrl(url)
	if len(probs) == 0 {
		return "", ""
	}
	return fileServices.GetInputOutputDirectories(probs[0].Platform, probs[0].ContestId, probs[0].Label)
}
