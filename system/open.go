package system

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/skmonir/mango-gui/context"
)

// https://gist.github.com/hyg/9c4afcd91fe24316cbf0
func OpenSourceList(config context.AppConfig, problemIdList []string) error {
	var err error
	for _, problemId := range problemIdList {
		sourcePath := config.GetSourceFilePathWithExt(problemId)

		switch runtime.GOOS {
		case "linux":
			err = exec.Command("xdg-open", sourcePath).Run()
		case "windows":
			exec.Command("cmd", fmt.Sprintf("/C start %v", sourcePath)).Run()
		case "darwin":
			err = exec.Command("open", sourcePath).Run()
		default:
			err = errors.New("unsupported os")
		}
	}
	return err
}

func Open(ctx *context.AppCtx, contestId string, problemId string) error {
	if contestId == "" {
		return errors.New("contest id not found")
	}

	if problemId == "" {
		return errors.New("problem id not found")
	}

	ctx.Config.CurrentContestId = contestId
	ctx.Config.SaveConfig()

	if err := OpenSourceList(*ctx.Config, []string{problemId}); err != nil {
		return errors.New("error while opening source")
	}
	return nil
}
