package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/skmonir/mango-gui/backend/judge-framework/constants"
	"github.com/skmonir/mango-gui/backend/judge-framework/logger"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
)

func GetJudgeConfigFromFile() *JudgeConfig {
	var err error
	config := JudgeConfig{}
	configFilePath := getConfigFilePath()

	if !utils.IsFileExist(configFilePath) {
		if config, err = CreateDefaultConfig(); err != nil {
			return &config
		}
	}

	data, e := ioutil.ReadFile(configFilePath)
	if e != nil {
		fmt.Println("Error: ", err)
		return &config
	}

	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Println("Error: ", err)
		return &config
	}

	fmt.Println("Returning app config from file")

	return &config
}

func SaveConfigIntoJsonFile(config JudgeConfig) error {
	data, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		return err
	}

	mu := sync.Mutex{}

	configFilePath := getConfigFilePath()
	mu.Lock()
	err = ioutil.WriteFile(configFilePath, data, 0644)
	mu.Unlock()
	if err != nil {
		return err
	}
	return nil
}

func getConfigFilePath() string {
	return filepath.Join(getConfigDirectoryPath(), "config.json")
}

func isConfigDirExist() bool {
	_, err := os.Stat(getConfigDirectoryPath())
	if err != nil || os.IsNotExist(err) {
		return false
	}
	return true
}

func createConfigDir() error {
	if !isConfigDirExist() {
		fmt.Println("Creating config directory " + getConfigDirectoryPath())
		if err := os.MkdirAll(getConfigDirectoryPath(), os.ModePerm); err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func getConfigDirectoryPath() string {
	return filepath.Join(utils.GetAppHomeDirectoryPath(), "appdata")
}

func CreateDefaultConfig() (JudgeConfig, error) {
	if err := createConfigDir(); err != nil {
		return JudgeConfig{}, err
	}

	conf := JudgeConfig{
		AppVersion:   constants.APP_VERSION,
		ActiveLang:   "cpp",
		ActiveLangId: "",
		LangConfigs: map[string]LanguageConfig{
			"cpp": {
				Lang:                "CPP",
				CompilationCommand:  "g++",
				CompilationFlags:    map[bool]string{true: "-std=gnu++17", false: "-std=c++20"}[utils.IsOsWindows()],
				FileExtension:       ".cpp",
				DefaultTemplatePath: utils.GetDefaultTemplateFilePathByLang("cpp"),
			},
			"java": {
				Lang:                "Java",
				CompilationCommand:  "javac",
				CompilationFlags:    "-encoding UTF-8 -J-Xmx2048m",
				ExecutionCommand:    "java",
				ExecutionFlags:      "-XX:+UseSerialGC -Xss64m -Xms64m -Xmx2048m",
				FileExtension:       ".java",
				DefaultTemplatePath: utils.GetDefaultTemplateFilePathByLang("java"),
			},
			"python": {
				Lang:                "Python",
				CompilationCommand:  map[bool]string{true: "py", false: "python3"}[utils.IsOsWindows()],
				ExecutionCommand:    map[bool]string{true: "py", false: "python3"}[utils.IsOsWindows()],
				FileExtension:       ".py",
				DefaultTemplatePath: utils.GetDefaultTemplateFilePathByLang("python"),
			},
		},
		EditorPreference: EditorPreferences{
			Theme:    "monokai",
			FontSize: "14",
			TabSize:  "4",
		},
	}

	if err := SaveConfigIntoJsonFile(conf); err != nil {
		return JudgeConfig{}, err
	}
	logger.Info("Default app config is created")

	return conf, nil
}
