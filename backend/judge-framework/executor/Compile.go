package executor

import (
	"bytes"
	"fmt"
	"github.com/skmonir/mango-gui/backend/judge-framework/logger"
	"github.com/skmonir/mango-gui/backend/judge-framework/services/languageServices"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/skmonir/mango-gui/backend/judge-framework/config"
)

func CompileSource(command string, showStdError bool) string {
	cmds := utils.ParseCommand(command)

	// ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	// defer cancel()
	// cmd := exec.CommandContext(ctx, cmds[0], cmds[1:]...)

	var stderr_buffer bytes.Buffer
	cmd := exec.Command(cmds[0], cmds[1:]...)
	if showStdError {
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stderr = &stderr_buffer
	}
	if err := cmd.Run(); err != nil {
		if stderr_buffer.String() != "" {
			logger.Error(stderr_buffer.String())
			return stderr_buffer.String()
		}
		logger.Error(err.Error())
		return err.Error()
	}
	return ""
}

func Compile(platform string, cid string, label string) string {
	logger.Info(fmt.Sprintf("Compiling source for %v %v %v", platform, cid, label))
	judgeConfig := config.GetJudgeConfigFromCache()

	sourceFolderPath := filepath.Join(judgeConfig.WorkspaceDirectory, platform, cid, "source")
	filePathWithoutExt := filepath.Join(sourceFolderPath, label)

	err, command := languageServices.GetCompilationCommand(filePathWithoutExt, judgeConfig.TestingLangConfigs[judgeConfig.ActiveTestingLang])
	if err != nil {
		logger.Error(err.Error())
		return err.Error()
	}

	return CompileSource(command, false)
}
