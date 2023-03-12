package executor

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"unicode"

	"github.com/skmonir/mango-gui/backend/judge-framework/logger"

	"github.com/skmonir/mango-gui/backend/judge-framework/config"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
)

// parseCommand parses a command line and handle arguments in quotes.
// https://github.com/vrischmann/shlex/blob/master/shlex.go
func parseCommand(s string) (res []string) {
	var buf bytes.Buffer
	insideQuotes := false
	for _, r := range s {
		switch {
		case unicode.IsSpace(r) && !insideQuotes:
			if buf.Len() > 0 {
				res = append(res, buf.String())
				buf.Reset()
			}
		case r == '"' || r == '\'':
			if insideQuotes {
				res = append(res, buf.String())
				buf.Reset()
				insideQuotes = false
				continue
			}
			insideQuotes = true
		default:
			buf.WriteRune(r)
		}
	}
	if buf.Len() > 0 {
		res = append(res, buf.String())
	}
	return
}

func CompileSource(command string, showStdError bool) string {
	cmds := parseCommand(command)

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

	folderPath := filepath.Join(judgeConfig.WorkspaceDirectory, platform, cid, "source")
	filePathWithoutExt := filepath.Join(folderPath, label)
	filePathWithExt := filePathWithoutExt + judgeConfig.ActiveLanguage.FileExtension

	if !utils.IsFileExist(filePathWithExt) {
		logger.Error("Source file not found!")
		return "Source file not found!"
	}

	command := fmt.Sprintf("%v %v %v -o %v", judgeConfig.ActiveLanguage.CompilationCommand, judgeConfig.ActiveLanguage.CompilationArgs, filePathWithExt, filePathWithoutExt)

	return CompileSource(command, false)
}
