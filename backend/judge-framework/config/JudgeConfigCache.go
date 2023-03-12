package config

import (
	"github.com/skmonir/mango-gui/backend/judge-framework/logger"
	"sync"
)

var once sync.Once
var judgeConfig *JudgeConfig

func GetJudgeConfigFromCache() *JudgeConfig {
	if judgeConfig == nil {
		logger.Info("Config not available in cache")
		once.Do(func() {
			judgeConfig = GetJudgeConfigFromFile()
		})
	} else {
		logger.Info("Returning config from cache")
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
