package services

import (
	"github.com/skmonir/mango-gui/backend/judge-framework/config"
	"github.com/skmonir/mango-gui/backend/judge-framework/models"
	"github.com/skmonir/mango-gui/backend/socket"
)

func UpdateEditorPreference(preferences config.EditorPreferences) config.EditorPreferences {
	judgeConfig := config.GetJudgeConfigFromCache()
	judgeConfig.EditorPreference = preferences
	judgeConfig = config.UpdateJudgeConfigIntoCache(*judgeConfig)
	socket.PublishAppConfig(*judgeConfig)
	return judgeConfig.EditorPreference
}

func UpdateJudgeAccountInfo(platform, handleOrEmail, password, handle string) {
	conf := config.GetJudgeConfigFromCache()
	submissionLangId := conf.JudgeAccInfo[platform].SubmissionLangId
	conf.JudgeAccInfo[platform] = models.JudgeAccountInfo{
		Handle:           handle,
		HandleOrEmail:    handleOrEmail,
		Password:         password,
		SubmissionLangId: submissionLangId,
	}
	conf = config.UpdateJudgeConfigIntoCache(*conf)
	socket.PublishAppConfig(*conf)
}
