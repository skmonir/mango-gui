package context

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/skmonir/mango-gui/models"
	"github.com/skmonir/mango-gui/utils"
)

func GetConfigFolderPath() string {
	cfgPath := ""
	switch runtime.GOOS {
	case "linux":
		if os.Getenv("XDG_CONFIG_HOME") != "" {
			cfgPath = os.Getenv("XDG_CONFIG_HOME")
		} else {
			cfgPath = filepath.Join(os.Getenv("HOME"), ".mango")
		}
	case "windows":
		cfgPath = filepath.Join(os.Getenv("APPDATA"), "mango")
	case "darwin":
		cfgPath = filepath.Join(os.Getenv("HOME"), ".mango")
	default:
		cfgPath = ""
	}

	return cfgPath
}

func GetConfigFilePath() string {
	return filepath.Join(GetConfigFolderPath(), "config.json")
}

func IsConfigDirExist() bool {
	_, err := os.Stat(GetConfigFolderPath())
	if err != nil || os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateConfigDir() error {
	if !IsConfigDirExist() {
		if err := os.MkdirAll(GetConfigFolderPath(), os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func GetAppConfig() *AppConfig {
	var err error
	config := &AppConfig{}
	if !IsConfigExist() {
		if config, err = CreateDefaultConfig(); err != nil {
			return config
		}
	}

	cfgFilePath := GetConfigFilePath()

	data, e := ioutil.ReadFile(cfgFilePath)
	if e != nil {
		return config
	}

	json.Unmarshal(data, &config)

	return config
}

func (config *AppConfig) SaveConfig() error {
	data, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		return err
	}

	mu := sync.Mutex{}

	cfgPath := GetConfigFilePath()
	mu.Lock()
	err = ioutil.WriteFile(cfgPath, data, 0644)
	mu.Unlock()
	if err != nil {
		return err
	}

	return nil
}

func CreateDefaultConfig() (*AppConfig, error) {
	CreateConfigDir()

	config := &AppConfig{
		CompilationCommand: "g++",
		CompilationArgs:    "-std=c++17",
		SrcDir:             "src",
		TestDir:            "testcase",
		CurrentTheme:       "dark",
	}

	if err := config.SaveConfig(); err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("error while creating default AppConfig")
	}

	return config, nil
}

func (config *AppConfig) SetOnlineJudge(OJ string) error {
	if config == nil {
		return errors.New("could not find AppConfig")
	}

	config.OJ = OJ

	if err := config.SaveConfig(); err != nil {
		return err
	}

	return nil
}

func (config *AppConfig) SetContest(contestId string) error {
	if _, err := strconv.Atoi(contestId); err != nil {
		return errors.New("contest id not valid")
	}

	if config == nil {
		return errors.New("could not find AppConfig")
	}

	config.CurrentContestId = contestId

	if err := config.SaveConfig(); err != nil {
		return err
	}

	return nil
}

func (config *AppConfig) GetSourceDirPath() string {
	return filepath.Join(config.Workspace, strings.ToLower(config.OJ), config.CurrentContestId, config.SrcDir)
}

func (config *AppConfig) GetSourceFilePathWithExt(problemId string) string {
	return filepath.Join(config.GetSourceDirPath(), problemId+".cpp")
}

func (config *AppConfig) GetSourceFilePathWithoutExt(problemId string) string {
	return filepath.Join(config.GetSourceDirPath(), problemId)
}

func (config *AppConfig) GetTestcaseDirPath() string {
	return filepath.Join(config.Workspace, strings.ToLower(config.OJ), config.CurrentContestId, config.TestDir)
}

func (config *AppConfig) GetTestcaseFilePath(problemId string) string {
	return filepath.Join(config.GetTestcaseDirPath(), problemId+".json")
}

func (config *AppConfig) ResolveTescasePath(problemId string) error {
	testCaseDirPath := config.GetTestcaseDirPath()

	if err := utils.CreateFile(testCaseDirPath, problemId+".json"); err != nil {
		return err
	}

	return nil
}

func (config *AppConfig) GetProblemInfo(problemId string) (models.Problem, error) {
	var problemInfo models.Problem
	testpath := config.GetTestcaseFilePath(problemId)

	data, err := ioutil.ReadFile(testpath)
	if err != nil {
		return problemInfo, err
	}

	err = json.Unmarshal(data, &problemInfo)
	if err != nil {
		return problemInfo, err
	}

	return problemInfo, nil
}

func (config *AppConfig) GetProblemIdListForTester() []string {
	files, err := ioutil.ReadDir(config.GetTestcaseDirPath())
	if err != nil {
		return []string{}
	}

	problemIdList := []string{}
	for _, file := range files {
		filename := strings.Split(file.Name(), ".")[0]
		problemIdList = append(problemIdList, filename)
	}
	return problemIdList
}

func Configure() error {
	var err error

	if _, err = CreateDefaultConfig(); err != nil {
		return err
	}

	cfgPath := GetConfigFilePath()

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", cfgPath).Run()
	case "windows":
		exec.Command("cmd", fmt.Sprintf("/C start %v", cfgPath)).Run()
	case "darwin":
		err = exec.Command("open", cfgPath).Run()
	default:
		// ansi.Println(color.New(color.FgRed).Sprintf("unsupported os"))
	}

	return err
}

func IsConfigExist() bool {
	cfgFilePath := GetConfigFilePath()
	info, err := os.Stat(cfgFilePath)
	if err != nil || os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
