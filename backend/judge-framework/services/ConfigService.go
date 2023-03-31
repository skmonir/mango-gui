package services

import (
	"github.com/skmonir/mango-gui/backend/judge-framework/config"
	"github.com/skmonir/mango-gui/backend/socket"
)

func UpdateEditorPreference(preferences config.EditorPreferences) config.EditorPreferences {
	judgeConfig := config.GetJudgeConfigFromCache()
	judgeConfig.EditorPreference = preferences
	judgeConfig = config.UpdateJudgeConfigIntoCache(*judgeConfig)
	socket.PublishAppConfig(*judgeConfig)
	return judgeConfig.EditorPreference
}
