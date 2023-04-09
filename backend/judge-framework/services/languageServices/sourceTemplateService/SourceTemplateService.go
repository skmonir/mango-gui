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

func CreateDefaultTemplatesIfNotAvailable() {
	templateDirectory := filepath.Join(utils.GetAppHomeDirectoryPath(), "source_templates")
	templateFiles := []string{"template_CPP.cpp", "template_Java.java", "template_Python.py"}

	allFileAvailable := true
	for _, filename := range templateFiles {
		filePath := filepath.Join(templateDirectory, filename)
		allFileAvailable = allFileAvailable && utils.IsFileExist(filePath)
	}

	if !allFileAvailable {
		CreateDefaultTemplateFiles()
	}
}

func CreateDefaultTemplateFiles() {
	templateDirectory := filepath.Join(utils.GetAppHomeDirectoryPath(), "source_templates")
	utils.WriteFileContent(templateDirectory, "template_CPP.cpp", templateCpp)
	utils.WriteFileContent(templateDirectory, "template_Java.java", templateJava)
	utils.WriteFileContent(templateDirectory, "template_Python.py", templatePython)
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
	CreateDefaultTemplatesIfNotAvailable()

	templateFilePath := utils.GetDefaultTemplateFilePathByLang(lang)
	return utils.ReadFileContent(templateFilePath, constants.SOURCE_MAX_ROW_FOR_TEST, constants.SOURCE_MAX_COL_FOR_TEST)
}
