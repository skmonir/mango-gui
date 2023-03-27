package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

func GetJudgeConfigFromFile() *JudgeConfig {
	fmt.Println("Getting config from file...")
	var err error
	config := JudgeConfig{}
	if !isConfigExist() {
		fmt.Println("Config doesn't exist...")
		if config, err = CreateDefaultConfig(); err != nil {
			return &config
		}
	}

	configFilePath := getConfigFilePath()

	fmt.Println("Reading config from " + configFilePath)
	data, e := ioutil.ReadFile(configFilePath)
	if e != nil {
		fmt.Println("Error: ", err)
		return &config
	}

	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Println("Error: ", err)
		return &config
	}

	fmt.Println(config)

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
	return filepath.Join(utils.GetAppDataDirectoryPath(), "config.json")
}

func isConfigDirExist() bool {
	_, err := os.Stat(utils.GetAppDataDirectoryPath())
	if err != nil || os.IsNotExist(err) {
		return false
	}
	return true
}

func createConfigDir() error {
	if !isConfigDirExist() {
		fmt.Println("Creating config directory " + utils.GetAppDataDirectoryPath())
		if err := os.MkdirAll(utils.GetAppDataDirectoryPath(), os.ModePerm); err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func isConfigExist() bool {
	cfgFilePath := getConfigFilePath()
	info, err := os.Stat(cfgFilePath)
	if err != nil || os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func CreateDefaultConfig() (JudgeConfig, error) {
	fmt.Println("Creating default conf...")
	if err := createConfigDir(); err != nil {
		return JudgeConfig{}, err
	}

	conf := JudgeConfig{
		ActiveLang: "cpp",
		LangConfigs: map[string]LanguageConfig{
			"cpp": {
				Lang:               "CPP",
				CompilationCommand: "g++",
				CompilationFlags:   "-std=c++20",
				FileExtension:      ".cpp",
			},
			"java": {
				Lang:               "Java",
				CompilationCommand: "javac",
				CompilationFlags:   "-encoding UTF-8 -J-Xmx2048m",
				ExecutionCommand:   "java",
				ExecutionFlags:     "-XX:+UseSerialGC -Xss64m -Xms64m -Xmx2048m",
				FileExtension:      ".java",
			},
			"python": {
				Lang:             "Python",
				ExecutionCommand: "python3",
				FileExtension:    ".py",
			},
		},
	}

	if err := SaveConfigIntoJsonFile(conf); err != nil {
		fmt.Println(err.Error())
		return JudgeConfig{}, errors.New("error while creating default AppConfig")
	}

	return conf, nil
}
