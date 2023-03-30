package sourceTemplateService

import (
	"github.com/skmonir/mango-gui/backend/judge-framework/config"
	"github.com/skmonir/mango-gui/backend/judge-framework/constants"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
	"path/filepath"
	"strings"
)

import _ "embed"

//go:embed template_CPP.txt
var templateCpp []byte

//go:embed template_Java.txt
var templateJava []byte

//go:embed template_Python.txt
var templatePython []byte

func CreateIfTemplatesAreNotAvailable() {
	appdataDirectory := utils.GetAppDataDirectoryPath()
	templateDirectory := filepath.Join(appdataDirectory, "source_templates")
	templateFiles := []string{"template_CPP.txt", "template_Java.txt", "template_Python.txt"}

	allFileAvailable := true
	for _, filename := range templateFiles {
		filePath := filepath.Join(templateDirectory, filename)
		allFileAvailable = allFileAvailable && utils.IsFileExist(filePath)
	}

	if !allFileAvailable {
		CreateTemplateFiles()
	}
}

func CreateTemplateFiles() {
	appdataDirectory := utils.GetAppDataDirectoryPath()
	templateDirectory := filepath.Join(appdataDirectory, "source_templates")
	utils.WriteFileContent(templateDirectory, "template_CPP.txt", templateCpp)
	utils.WriteFileContent(templateDirectory, "template_Java.txt", templateJava)
	utils.WriteFileContent(templateDirectory, "template_Python.txt", templatePython)
}

func GetTemplateCode() string {
	conf := config.GetJudgeConfigFromCache()

	body := ""
	if len(conf.LangConfigs[conf.ActiveLang].UserTemplatePath) > 0 {
		body = utils.ReadFileContent(conf.LangConfigs[conf.ActiveLang].UserTemplatePath, constants.SOURCE_MAX_ROW_FOR_TEST, constants.SOURCE_MAX_COL_FOR_TEST)
		if len(strings.Trim(body, " \n\t")) == 0 {
			body = GetDefaultTemplate(conf.ActiveLang)
		}
	} else {
		body = GetDefaultTemplate(conf.ActiveLang)
	}
	return body
}

func GetDefaultTemplate(lang string) string {
	CreateIfTemplatesAreNotAvailable()

	templateFilePath := utils.GetTemplateFilePathByLang(lang)
	return utils.ReadFileContent(templateFilePath, constants.SOURCE_MAX_ROW_FOR_TEST, constants.SOURCE_MAX_COL_FOR_TEST)
}
