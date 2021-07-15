package system

import (
	"github.com/skmonir/mango-gui/context"
	"github.com/skmonir/mango-gui/models"
)

func Create(ctx *context.AppCtx) []models.Problem {
	ctx.ProgressBar.SetValue(0)

	var parser Parser
	var problemInfoList []models.Problem

	if ctx.TesterUi.OnlineJudgeOptionSelect.Selected == "CodeForces" {
		parser = CodeforcesParser{}
	}

	contestId, problemId, err := parser.ParseContestAndProblemId(ctx.TesterUi.ContestIdInputField.Text)
	if err != nil || contestId == "" {
		return problemInfoList
	}

	ctx.Config.CurrentContestId = contestId
	ctx.Config.OJ = ctx.TesterUi.OnlineJudgeOptionSelect.Selected
	ctx.Config.SaveConfig()

	problemIdList := []string{}
	if problemId == "" {
		problemIdList = ctx.Config.GetProblemIdListForTester()
	} else {
		problemIdList = append(problemIdList, problemId)
	}

	if len(problemIdList) > 0 {
		ctx.ProgressBar.Max = float64(len(problemIdList))
	}

	for i, problemId := range problemIdList {
		problemInfo, err := ctx.Config.GetProblemInfo(problemId)
		if err == nil {
			problemInfoList = append(problemInfoList, problemInfo)
		}
		ctx.ProgressBar.SetValue(float64(i + 1))
	}

	return problemInfoList
}
