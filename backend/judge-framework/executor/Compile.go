package executor

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"unicode"

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
	}
	cmd.Stderr = &stderr_buffer
	if err := cmd.Run(); err != nil {
		fmt.Println(stderr_buffer.String())
		return stderr_buffer.String()
	}
	return ""
}

func Compile(platform string, cid string, label string) string {
	fmt.Println("Compiling source...")
	judgeConfig := config.GetJudgeConfigFromCache()

	folderPath := fmt.Sprintf("%v/%v/%v/source", strings.TrimRight(judgeConfig.WorkspaceDirectory, "/"), platform, cid)
	filePathWithoutExt := folderPath + "/" + label
	filePathWithExt := folderPath + "/" + label + judgeConfig.ActiveLanguage.FileExtension

	if !utils.IsFileExist(filePathWithExt) {
		return "Source file not found!"
	}

	command := fmt.Sprintf("%v %v %v -o %v", judgeConfig.ActiveLanguage.CompilationCommand, judgeConfig.ActiveLanguage.CompilationArgs, filePathWithExt, filePathWithoutExt)

	return CompileSource(command, false)
}
