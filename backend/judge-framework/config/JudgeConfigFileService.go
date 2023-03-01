package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

func GetJudgeConfigFromFile() *JudgeConfig {
	fmt.Println("Getting config from fileService...")
	var err error
	config := JudgeConfig{}
	if !isConfigExist() {
		fmt.Println("Config doesn't exist...")
		if config, err = createDefaultConfig(); err != nil {
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

func getConfigFolderPath() string {
	configPath := ""
	switch runtime.GOOS {
	case "linux":
		if os.Getenv("XDG_CONFIG_HOME") != "" {
			configPath = os.Getenv("XDG_CONFIG_HOME")
		} else {
			configPath = filepath.Join(os.Getenv("HOME"), ".mango")
		}
	case "windows":
		configPath = filepath.Join(os.Getenv("APPDATA"), "mango")
	case "darwin":
		configPath = filepath.Join(os.Getenv("HOME"), ".mango")
	default:
		configPath = ""
	}

	return configPath
}

func getConfigFilePath() string {
	return filepath.Join(getConfigFolderPath(), "config.json")
}

func isConfigDirExist() bool {
	_, err := os.Stat(getConfigFolderPath())
	if err != nil || os.IsNotExist(err) {
		return false
	}
	return true
}

func createConfigDir() error {
	if !isConfigDirExist() {
		fmt.Println("Creating config directory " + getConfigFolderPath())
		if err := os.MkdirAll(getConfigFolderPath(), os.ModePerm); err != nil {
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

func createDefaultConfig() (JudgeConfig, error) {
	fmt.Println("Creating default conf...")
	if err := createConfigDir(); err != nil {
		return JudgeConfig{}, err
	}

	activeLang := LanguageConfig{
		Lang:               "c++",
		CompilationCommand: "g++",
		CompilationArgs:    "-std=c++20",
		FileExtension:      ".cpp",
	}
	conf := JudgeConfig{
		ActiveLanguage:  activeLang,
		LanguageConfigs: []LanguageConfig{activeLang},
	}

	if err := SaveConfigIntoJsonFile(conf); err != nil {
		fmt.Println(err.Error())
		return JudgeConfig{}, errors.New("error while creating default AppConfig")
	}

	return conf, nil
}
