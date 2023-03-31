package tgenScripts

import (
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
	"path/filepath"
)

import _ "embed"

//go:embed validator.txt
var validatorContent []byte

//go:embed generator.txt
var generatorContent []byte

//go:embed testlib.txt
var testlibContent []byte

func CreateGeneratorScriptsIfNotAvailable() {
	appdataDirectory := utils.GetAppDataDirectoryPath()
	scriptDirectory := filepath.Join(appdataDirectory, "tgen_scripts")
	scriptFiles := []string{"validator.cpp", "generator.cpp", "testlib.h"}

	allFileAvailable := true
	for _, filename := range scriptFiles {
		filePath := filepath.Join(scriptDirectory, filename)
		allFileAvailable = allFileAvailable && utils.IsFileExist(filePath)
	}

	if !allFileAvailable {
		createScriptFiles()
	}
}

func createScriptFiles() {
	appdataDirectory := utils.GetAppDataDirectoryPath()
	scriptDirectory := filepath.Join(appdataDirectory, "tgen_scripts")
	utils.WriteFileContent(scriptDirectory, "validator.cpp", validatorContent)
	utils.WriteFileContent(scriptDirectory, "generator.cpp", generatorContent)
	utils.WriteFileContent(scriptDirectory, "testlib.h", testlibContent)
}
