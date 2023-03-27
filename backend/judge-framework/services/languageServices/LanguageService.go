package languageServices

import (
	"errors"
	"fmt"
	"github.com/skmonir/mango-gui/backend/judge-framework/config"
	"github.com/skmonir/mango-gui/backend/judge-framework/logger"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
	"path/filepath"
	"strings"
)

func GetCompilationCommand(programFilePathWithoutExt string, lang config.LanguageConfig) (error, string) {
	programFilePathWithExt := programFilePathWithoutExt + lang.FileExtension

	if !utils.IsFileExist(programFilePathWithExt) {
		return errors.New("Source file not found!"), ""
	}

	command := ""
	if lang.Lang == "CPP" {
		command = fmt.Sprintf("%v %v %v -o %v%v", lang.CompilationCommand, lang.CompilationFlags, programFilePathWithExt, programFilePathWithoutExt, utils.GetBinaryFileExt())
	} else if lang.Lang == "Java" {
		command = fmt.Sprintf("%v %v %v", lang.CompilationCommand, lang.CompilationFlags, programFilePathWithExt)
	} else if lang.Lang == "Python" {
		command = fmt.Sprintf("%v -m py_compile %v", lang.CompilationCommand, programFilePathWithExt)
	}

	logger.Info("Prepared compilation command: " + command)
	return nil, command
}

func GetExecutionCommandByFilePath(filePathWithExt string) []string {
	conf := config.GetJudgeConfigFromCache()

	fileExt := filepath.Ext(filePathWithExt)
	filePathWithoutExt := strings.TrimSuffix(filePathWithExt, fileExt)

	command := ""
	if fileExt == ".cpp" || fileExt == ".cc" {
		command = fmt.Sprintf("%v%v", filePathWithoutExt, utils.GetBinaryFileExt())
	} else if fileExt == ".java" {
		command = fmt.Sprintf("%v %v %v", conf.LangConfigs[conf.ActiveLang].ExecutionCommand, conf.LangConfigs[conf.ActiveLang].ExecutionFlags, filePathWithoutExt)
	} else if fileExt == ".py" {
		command = fmt.Sprintf("%v %v", conf.LangConfigs[conf.ActiveLang].ExecutionCommand, filePathWithExt)
	}

	return utils.ParseCommand(command)
}

func GetLangConfigFromFileExt(ext string) config.LanguageConfig {
	conf := config.GetJudgeConfigFromCache()
	lang := ""
	if ext == ".cpp" || ext == ".cc" {
		lang = "cpp"
	} else if ext == ".java" {
		lang = "java"
	} else if ext == ".py" {
		lang = "python"
	}
	return conf.LangConfigs[lang]
}

func GetBinaryFilePathByFilePath(filePathWithExt string) string {
	fileExt := filepath.Ext(filePathWithExt)
	filePathWithoutExt := strings.TrimSuffix(filePathWithExt, fileExt)

	scriptBinaryPath := ""
	if fileExt == ".cpp" || fileExt == ".cc" {
		scriptBinaryPath = filePathWithoutExt + utils.GetBinaryFileExt()
	} else if fileExt == ".java" {
		scriptBinaryPath = filePathWithoutExt + utils.GetBinaryFileExt()
	} else if fileExt == ".py" {
		scriptBinaryPath = filePathWithExt
	}

	return scriptBinaryPath
}
