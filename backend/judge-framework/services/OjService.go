package services

import (
	"github.com/skmonir/mango-gui/backend/judge-framework/config"
	"github.com/skmonir/mango-gui/backend/judge-framework/oj/client"
	"github.com/skmonir/mango-gui/backend/judge-framework/services/fileServices"
	"github.com/skmonir/mango-gui/backend/socket"
)

func Login(platform, handleOrEmail, password string) (err error, handle string) {
	err, ojClient := client.GetClientByPlatform(platform)
	if err != nil {
		return err, ""
	}

	err, handle = ojClient.DoLogin(handleOrEmail, password)
	if err != nil {
		return err, ""
	}

	UpdateJudgeAccountInfo(platform, handleOrEmail, password, handle)
	return
}

func Submit(platform, cid, pid string) (error, string) {
	err, ojClient := client.GetClientByPlatform(platform)
	if err != nil {
		return err, ""
	}
	problem := GetProblem(platform, cid, pid)
	conf := config.GetJudgeConfigFromCache()
	sourceResp := fileServices.GetCodeByMetadata(platform, cid, pid)
	socket.PublishStatusMessage("test_status", "Submitting code...", "info")
	err = ojClient.Submit(problem, conf.JudgeAccInfo[platform].SubmissionLangId, sourceResp["code"])
	if err != nil {
		socket.PublishStatusMessage("test_status", err.Error(), "error")
		return err, ""
	}
	return nil, ""
}
