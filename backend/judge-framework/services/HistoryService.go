package services

import (
	"encoding/json"
	"github.com/skmonir/mango-gui/backend/judge-framework/constants"
	"github.com/skmonir/mango-gui/backend/judge-framework/dto"
	"github.com/skmonir/mango-gui/backend/judge-framework/logger"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
	"path/filepath"
)

type history struct {
	ParseUrl              string                      `json:"parseUrl"`
	TestContestUrl        string                      `json:"testContestUrl"`
	InputGenerateRequest  dto.TestcaseGenerateRequest `json:"inputGenerateRequest"`
	OutputGenerateRequest dto.TestcaseGenerateRequest `json:"outputGenerateRequest"`
}

func GetHistory() string {
	historyData := fetchHistoryFromFile()
	historyJson, err := json.Marshal(historyData)
	if err != nil {
		logger.Error(err.Error())
	}
	return string(historyJson)
}

func UpdateInputGenerateRequestHistory(req dto.TestcaseGenerateRequest) {
	history := fetchHistoryFromFile()
	history.InputGenerateRequest = req
	updateHistory(history)
}

func UpdateOutputGenerateRequestHistory(req dto.TestcaseGenerateRequest) {
	history := fetchHistoryFromFile()
	history.OutputGenerateRequest = req
	updateHistory(history)
}

func UpdateParseUrlHistory(url string) {
	history := fetchHistoryFromFile()
	history.ParseUrl = url
	updateHistory(history)
}

func UpdateTestContestUrlHistory(url string) {
	history := fetchHistoryFromFile()
	history.TestContestUrl = url
	updateHistory(history)
}

func InitHistory() {
	historyPath := filepath.Join(utils.GetAppDataDirectoryPath(), "history.json")
	if !utils.IsFileExist(historyPath) {
		updateHistory(getDefaultHistory())
	}
}

func fetchHistoryFromFile() history {
	hist := getDefaultHistory()
	historyPath := filepath.Join(utils.GetAppDataDirectoryPath(), "history.json")
	historyData := utils.ReadFileContent(historyPath, constants.IO_MAX_ROW_FOR_TEST, constants.IO_MAX_COL_FOR_TEST)
	if len(historyData) == 0 {
		return hist
	}
	if err := json.Unmarshal([]byte(historyData), &hist); err != nil {
		logger.Error(err.Error())
	}
	return hist
}

func updateHistory(hist history) {
	historyJson, err := json.Marshal(hist)
	if err != nil {
		logger.Error(err.Error())
	}
	utils.WriteFileContent(utils.GetAppDataDirectoryPath(), "history.json", historyJson)
}

func getDefaultHistory() history {
	return history{
		InputGenerateRequest: dto.TestcaseGenerateRequest{
			FileNum:           1,
			FileMode:          "write",
			FileName:          "02_random_input",
			TestPerFile:       0,
			SerialFrom:        1,
			GenerationProcess: "tgen_script",
		},
	}
}
