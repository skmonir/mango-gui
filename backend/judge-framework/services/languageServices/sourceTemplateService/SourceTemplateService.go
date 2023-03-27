package sourceTemplateService

import (
	"github.com/skmonir/mango-gui/backend/judge-framework/config"
	"github.com/skmonir/mango-gui/backend/judge-framework/constants"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
	"strings"
)

import _ "embed"

//go:embed template_CPP.txt
var templateCpp []byte

//go:embed template_Java.txt
var templateJava []byte

//go:embed template_Python.txt
var templatePython []byte

func GetTemplateCode() string {
	conf := config.GetJudgeConfigFromCache()

	body := ""
	if len(conf.LangConfigs[conf.ActiveLang].TemplatePath) > 0 {
		body = utils.ReadFileContent(conf.LangConfigs[conf.ActiveLang].TemplatePath, constants.SOURCE_MAX_ROW_FOR_TEST, constants.SOURCE_MAX_COL_FOR_TEST)
		if len(strings.Trim(body, " \n\t")) == 0 {
			body = GetDefaultTemplate(conf.ActiveLang)
		}
	} else {
		body = GetDefaultTemplate(conf.ActiveLang)
	}
	return body
}

func GetDefaultTemplate(lang string) string {
	if lang == "cpp" {
		return string(templateCpp)
	} else if lang == "java" {
		return string(templateJava)
	} else if lang == "python" {
		return string(templatePython)
	}
	return ""
}
