package fileServices

import (
	"encoding/json"
	"fmt"
	"github.com/skmonir/mango-gui/backend/judge-framework/config"
	"github.com/skmonir/mango-gui/backend/judge-framework/models"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
	"strings"
)

func GetProblemListFromFile(platform string, cid string) []models.Problem {
	fmt.Println("Fetching problems from file...")

	judgeConfig := config.GetJudgeConfigFromCache()
	folderPath := fmt.Sprintf("%v/%v/%v", strings.TrimRight(judgeConfig.WorkspaceDirectory, "/"), platform, cid)
	fileName := "problems.json"
	filePath := folderPath + "/" + fileName

	var problems []models.Problem
	if utils.IsFileExist(filePath) {
		data := utils.ReadFileContent(filePath, 123456, 123456)

		if err := json.Unmarshal([]byte(data), &problems); err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("Fetched problems from file")
	return problems
}

func SaveProblemListIntoFile(platform string, cid string, problemList []models.Problem) {
	fmt.Println("Saving problems into file...")

	data, err := json.MarshalIndent(problemList, "", " ")
	if err != nil {
		fmt.Println(err)
		return
	}

	judgeConfig := config.GetJudgeConfigFromCache()
	folderPath := fmt.Sprintf("%v/%v/%v", strings.TrimRight(judgeConfig.WorkspaceDirectory, "/"), platform, cid)
	fileName := "problems.json"

	fmt.Println("Writing data into " + folderPath + "/" + fileName)
	fmt.Println("Saved problems into file")

	utils.WriteFileContent(folderPath, fileName, data)
}
