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
	ParseUrl              string
	TestContestUrl        string
	InputGenerateRequest  dto.TestcaseGenerateRequest
	OutputGenerateRequest dto.TestcaseGenerateRequest
}

func GetHistory() map[string]string {
	historyPath := filepath.Join(utils.GetAppDataDirectoryPath(), "history.json")
	historyData := utils.ReadFileContent(historyPath, constants.IO_MAX_ROW_FOR_TEST, constants.IO_MAX_COL_FOR_TEST)

	history := map[string]string{}
	if err := json.Unmarshal([]byte(historyData), &history); err != nil {
		logger.Error(err.Error())
	}
	return history
}

func UpdateHistory(key, value string) map[string]string {
	historyPath := filepath.Join(utils.GetAppDataDirectoryPath(), "history.json")
	historyData := utils.ReadFileContent(historyPath, constants.IO_MAX_ROW_FOR_TEST, constants.IO_MAX_COL_FOR_TEST)

	history := map[string]string{}
	if err := json.Unmarshal([]byte(historyData), &history); err != nil {
		logger.Error(err.Error())
	}
	return history
}
