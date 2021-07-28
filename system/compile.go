package system

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/skmonir/mango-gui/context"
	"github.com/skmonir/mango-gui/utils"
)

func getCompilationCommand(config context.AppConfig, problemId string) (string, error) {
	filePathWithExt := config.GetSourceFilePathWithExt(problemId)
	filePathWithoutExt := config.GetSourceFilePathWithoutExt(problemId)

	if !utils.IsFileExist(filePathWithExt) {
		if filePathWithExt = config.GetSourceFilePathWithExt(strings.ToLower(problemId)); !utils.IsFileExist(filePathWithExt) {
			return "", errors.New("source file not found")
		}
	}

	command := fmt.Sprintf("%v %v %v -o %v", config.CompilationCommand, config.CompilationArgs, filePathWithExt, filePathWithoutExt)

	return command, nil
}

func CompileSource(config context.AppConfig, problemId string) error {
	command, err := getCompilationCommand(config, problemId)
	if err != nil {
		return err
	}

	cmds := utils.ParseCommand(command)

	cmd := exec.Command(cmds[0], cmds[1:]...)

	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		return errors.New(stdErr.String())
	}

	return nil
}
