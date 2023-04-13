package fileServices

import (
	"errors"
	"fmt"
	"github.com/skmonir/mango-gui/backend/judge-framework/config"
	"github.com/skmonir/mango-gui/backend/judge-framework/logger"
	"github.com/skmonir/mango-gui/backend/judge-framework/models"
	"github.com/skmonir/mango-gui/backend/judge-framework/services/languageServices/sourceTemplateService"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
	"path/filepath"
	"strings"
	"time"
)

func CreateSourceFiles(problems []models.Problem) {
	logger.Info("Creating source files...")
	judgeConfig := config.GetJudgeConfigFromCache()
	for _, prob := range problems {
		saveSourceIntoFile(judgeConfig, prob)
	}
}

func saveSourceIntoFile(judgeConfig *config.JudgeConfig, problem models.Problem) {
	logger.Info(fmt.Sprintf("Saving source file for: %v", problem))
	folderPath := fmt.Sprintf("%v/%v/%v/source", strings.TrimRight(judgeConfig.WorkspaceDirectory, "/"), problem.Platform, problem.ContestId)
	fileName := problem.Label + judgeConfig.TestingLangConfigs[judgeConfig.ActiveTestingLang].FileExtension
	filePath := filepath.Join(folderPath, fileName)

	if !utils.IsFileExist(filePath) {
		template := getTemplate(judgeConfig, problem)
		utils.WriteFileContent(folderPath, fileName, []byte(template))
	}
}

func getTemplate(judgeConfig *config.JudgeConfig, problem models.Problem) string {
	body := sourceTemplateService.GetTemplateCode()

	body = strings.Replace(body, "{%AUTHOR%}", judgeConfig.Author, 1)
	body = strings.Replace(body, "{%CREATED_DATETIME%}", time.Now().Local().Format("2-Jan-2006 15:04:05"), 1)
	body = strings.Replace(body, "{%PROBLEM_NAME%}", problem.Label+" - "+problem.Name, 1)
	body = strings.Replace(body, "{%CLASS_NAME%}", problem.Label, 1)

	return body
}

func OpenSourceByMetadata(platform string, cid string, label string) error {
	filePath := GetSourceFilePath(platform, cid, label)
	if utils.IsFileExist(filePath) {
		return utils.OpenResourceInDefaultApplication(filePath)
	}
	logger.Error("Source file not found. Click Generate Source button.")
	return errors.New("Source file not found. Click Generate Source button.")
}

func GetCodeByMetadata(platform string, cid string, label string) map[string]string {
	filePath := GetSourceFilePath(platform, cid, label)
	code := ""
	if utils.IsFileExist(filePath) {
		code = utils.ReadFileContent(filePath, 123456, 123456)
	}
	return map[string]string{
		"lang": utils.GetLangNameByFileExt(filepath.Ext(filePath)),
		"code": code,
	}
}

func UpdateCodeByProblemPath(platform, cid, label, code string) {
	filePath := GetSourceFilePath(platform, cid, label)
	directory, filename := filepath.Split(filePath)
	utils.WriteFileContent(directory, filename, []byte(code))
}

func GetSourceFilePath(platform string, cid string, label string) string {
	judgeConfig := config.GetJudgeConfigFromCache()

	folderPath := filepath.Join(judgeConfig.WorkspaceDirectory, platform, cid, "source")
	fileName := label + judgeConfig.TestingLangConfigs[judgeConfig.ActiveTestingLang].FileExtension
	filePath := filepath.Join(folderPath, fileName)

	return filePath
}

func GenerateSourceByProblemPath(problem models.Problem) {
	problems := []models.Problem{problem}
	CreateSourceFiles(problems)
}
