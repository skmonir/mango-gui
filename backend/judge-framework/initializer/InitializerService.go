package initializer

import (
	"github.com/skmonir/mango-gui/backend/judge-framework/config"
	"github.com/skmonir/mango-gui/backend/judge-framework/constants"
	"github.com/skmonir/mango-gui/backend/judge-framework/logger"
	"github.com/skmonir/mango-gui/backend/judge-framework/services"
	"github.com/skmonir/mango-gui/backend/judge-framework/services/languageServices/sourceTemplateService"
	"github.com/skmonir/mango-gui/backend/judge-framework/services/testcaseGeneratorServices/tgenScripts"
	"github.com/skmonir/mango-gui/backend/socket"
	"time"
)

func InitializeJudgeFramework() {
	// Invalidate and Prepare app config
	invalidateAndPrepareAppConfig()
	socket.PublishStatusMessage("init_app_event", "Initializing...(1/5)", "info")
	time.Sleep(300 * time.Millisecond)

	// Prepare scripts for Input/Output Generator
	tgenScripts.CreateGeneratorScriptsIfNotAvailable()
	socket.PublishStatusMessage("init_app_event", "Initializing...(2/5)", "info")
	time.Sleep(300 * time.Millisecond)

	// Prepare template codes
	sourceTemplateService.CreateDefaultTemplatesIfNotAvailable()
	socket.PublishStatusMessage("init_app_event", "Initializing...(3/5)", "info")
	time.Sleep(300 * time.Millisecond)

	// Invalidate the log files
	logger.InvalidateLogFiles()
	socket.PublishStatusMessage("init_app_event", "Initializing...(4/5)", "info")
	time.Sleep(300 * time.Millisecond)

	// Init history file
	services.InitHistory()
	socket.PublishStatusMessage("init_app_event", "Initializing...(5/5)", "info")
	time.Sleep(300 * time.Millisecond)
}

func invalidateAndPrepareAppConfig() {
	conf := config.GetJudgeConfigFromCache()
	if conf.AppVersion != constants.APP_VERSION {
		// invalidation code goes here
	}
}
