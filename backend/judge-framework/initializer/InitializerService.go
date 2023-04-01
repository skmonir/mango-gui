package initializer

import (
	"github.com/skmonir/mango-gui/backend/judge-framework/config"
	"github.com/skmonir/mango-gui/backend/judge-framework/constants"
	"github.com/skmonir/mango-gui/backend/judge-framework/logger"
	"github.com/skmonir/mango-gui/backend/judge-framework/services"
	"github.com/skmonir/mango-gui/backend/judge-framework/services/languageServices/sourceTemplateService"
	"github.com/skmonir/mango-gui/backend/judge-framework/services/testcaseGeneratorServices/tgenScripts"
)

func InitializeJudgeFramework() {
	// Invalidate and Prepare app config
	invalidateAndPrepareAppConfig()

	// Prepare scripts for Input/Output Generator
	tgenScripts.CreateGeneratorScriptsIfNotAvailable()

	// Prepare template codes
	sourceTemplateService.CreateDefaultTemplatesIfNotAvailable()

	// Invalidate the log files
	logger.InvalidateLogFiles()

	// Init history file
	services.InitHistory()
}

func invalidateAndPrepareAppConfig() {
	conf := config.GetJudgeConfigFromCache()
	if conf.AppVersion != constants.APP_VERSION {
		// invalidation code goes here
	}
}
