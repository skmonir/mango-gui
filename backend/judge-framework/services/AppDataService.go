package services

import (
	"encoding/json"
	"github.com/skmonir/mango-gui/backend/judge-framework/constants"
	"github.com/skmonir/mango-gui/backend/judge-framework/dto"
	"github.com/skmonir/mango-gui/backend/judge-framework/logger"
	"github.com/skmonir/mango-gui/backend/judge-framework/models"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
	"github.com/skmonir/mango-gui/backend/socket"
	"path/filepath"
	"strings"
	"time"
)

type QueryHistory struct {
	ParseUrl              string                      `json:"parseUrl"`
	TestContestUrl        string                      `json:"testContestUrl"`
	InputGenerateRequest  dto.TestcaseGenerateRequest `json:"inputGenerateRequest"`
	OutputGenerateRequest dto.TestcaseGenerateRequest `json:"outputGenerateRequest"`
}

type AppData struct {
	ParseSchedulerTasks []models.ParseSchedulerTask `json:"parseSchedulerTasks"`
	QueryHistories      QueryHistory                `json:"queryHistories"`
}

/* App Data Section */

func InitAppDataIfNotAvailable() {
	appDataPath := filepath.Join(utils.GetAppHomeDirectoryPath(), "appdata", "appdata.json")
	if !utils.IsFileExist(appDataPath) {
		appData := getDefaultAppData()
		UpdateAppDataIntoFile(appData)
	}
}

func GetAppData() AppData {
	appData := fetchAppDataFromFile()
	futureTasks, expiredTasks := getFilteredParseScheduledTasks(appData.ParseSchedulerTasks)

	if len(expiredTasks) > 0 {
		slightExpiredTasks := []models.ParseSchedulerTask{}
		for i := 0; i < len(expiredTasks); i++ {
			if utils.IsTimeInFuture(expiredTasks[i].StartTime.Add(30 * time.Minute)) {
				slightExpiredTasks = append(slightExpiredTasks, expiredTasks[i])
			}
		}
		appData.ParseSchedulerTasks = append(futureTasks, slightExpiredTasks...)
		if len(expiredTasks) != len(slightExpiredTasks) {
			UpdateAppDataIntoFile(appData)
		}
	} else {
		appData.ParseSchedulerTasks = futureTasks
	}

	return appData
}

func fetchAppDataFromFile() AppData {
	appData := getDefaultAppData()
	appDataPath := filepath.Join(utils.GetAppHomeDirectoryPath(), "appdata", "appdata.json")
	appDataStr := utils.ReadFileContent(appDataPath, constants.IO_MAX_ROW_FOR_TEST, constants.IO_MAX_COL_FOR_TEST)
	if len(appDataStr) == 0 {
		return appData
	}
	if err := json.Unmarshal([]byte(appDataStr), &appData); err != nil {
		logger.Error(err.Error())
	}
	return appData
}

func UpdateAppDataIntoFile(appData AppData) {
	appDataBytes, err := json.MarshalIndent(appData, "", " ")
	if err != nil {
		logger.Error(err.Error())
	}
	utils.WriteFileContent(filepath.Join(utils.GetAppHomeDirectoryPath(), "appdata"), "appdata.json", appDataBytes)
}

func getDefaultAppData() AppData {
	return AppData{
		QueryHistories:      getDefaultHistory(),
		ParseSchedulerTasks: []models.ParseSchedulerTask{},
	}
}

/* Query History Section */

func UpdateInputGenerateRequestHistory(req dto.TestcaseGenerateRequest) {
	appData := GetAppData()
	appData.QueryHistories.InputGenerateRequest = req
	UpdateAppDataIntoFile(appData)
}

func UpdateOutputGenerateRequestHistory(req dto.TestcaseGenerateRequest) {
	appData := GetAppData()
	appData.QueryHistories.OutputGenerateRequest = req
	UpdateAppDataIntoFile(appData)
}

func UpdateParseUrlHistory(url string) {
	appData := GetAppData()
	appData.QueryHistories.ParseUrl = url
	UpdateAppDataIntoFile(appData)
}

func UpdateTestContestUrlHistory(url string) {
	appData := GetAppData()
	appData.QueryHistories.TestContestUrl = url
	UpdateAppDataIntoFile(appData)
}

func getDefaultHistory() QueryHistory {
	return QueryHistory{
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

/* Parse Scheduler Task Section */

func GetFutureParseScheduledTasks() []models.ParseSchedulerTask {
	appData := GetAppData()
	futureTasks, _ := getFilteredParseScheduledTasks(appData.ParseSchedulerTasks)
	return futureTasks
}

func getFilteredParseScheduledTasks(tasks []models.ParseSchedulerTask) ([]models.ParseSchedulerTask, []models.ParseSchedulerTask) {
	futureTasks := []models.ParseSchedulerTask{}
	expiredTasks := []models.ParseSchedulerTask{}
	for i := 0; i < len(tasks); i++ {
		if utils.IsTimeInFuture(tasks[i].StartTime) {
			futureTasks = append(futureTasks, tasks[i])
		} else {
			if strings.Contains(tasks[i].Stage, "SCHEDULED") {
				tasks[i].Stage = "EXPIRED"
			} else if tasks[i].Stage == "RUNNING" {
				tasks[i].Stage = "ABORTED"
			}
			expiredTasks = append(expiredTasks, tasks[i])
		}
	}
	return futureTasks, expiredTasks
}

func AddParseScheduledTask(newTask models.ParseSchedulerTask) {
	appData := GetAppData()
	futureTasks, expiredTasks := getFilteredParseScheduledTasks(appData.ParseSchedulerTasks)
	futureTasks = append(futureTasks, newTask)
	appData.ParseSchedulerTasks = append(futureTasks, expiredTasks...)
	UpdateAppDataIntoFile(appData)

	socket.PublishParseScheduledTasks(appData.ParseSchedulerTasks)
}

func RemoveParseScheduledTasksByIds(taskIds []string) {
	appData := GetAppData()
	newTasks := []models.ParseSchedulerTask{}
	for i := 0; i < len(appData.ParseSchedulerTasks); i++ {
		if !utils.SliceContains(taskIds, appData.ParseSchedulerTasks[i].Id) {
			newTasks = append(newTasks, appData.ParseSchedulerTasks[i])
		}
	}
	appData.ParseSchedulerTasks = newTasks
	UpdateAppDataIntoFile(appData)

	socket.PublishParseScheduledTasks(newTasks)
}

func UpdateParseScheduledTaskStage(taskId, stage string) {
	appData := GetAppData()
	for i := 0; i < len(appData.ParseSchedulerTasks); i++ {
		if appData.ParseSchedulerTasks[i].Id == taskId {
			appData.ParseSchedulerTasks[i].Stage = stage
		}
	}
	UpdateAppDataIntoFile(appData)

	socket.PublishParseScheduledTasks(appData.ParseSchedulerTasks)
}

func UpdateParseScheduledTask(updatedScheduleTask models.ParseSchedulerTask) {
	appData := GetAppData()
	for i := 0; i < len(appData.ParseSchedulerTasks); i++ {
		if appData.ParseSchedulerTasks[i].Id == updatedScheduleTask.Id {
			appData.ParseSchedulerTasks[i] = updatedScheduleTask
		}
	}
	UpdateAppDataIntoFile(appData)

	socket.PublishParseScheduledTasks(appData.ParseSchedulerTasks)
}
