package config

import (
	"fmt"
	"sync"
)

var once sync.Once
var judgeConfig *JudgeConfig

func GetJudgeConfigFromCache() *JudgeConfig {
	if judgeConfig == nil {
		fmt.Println("Config not available in cache")
		once.Do(func() {
			judgeConfig = GetJudgeConfigFromFile()
		})
	} else {
		fmt.Println("Returning config from cache")
	}
	return judgeConfig
}

func UpdateJudgeConfigIntoCache(config JudgeConfig) *JudgeConfig {
	if err := SaveConfigIntoJsonFile(config); err != nil {
		return judgeConfig
	}
	conf := GetJudgeConfigFromFile()
	judgeConfig = conf
	return conf
}
