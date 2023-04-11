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
	return ojClient.DoLogin(handleOrEmail, password)
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
	err = ojClient.Submit(problem, conf.ActiveLang, sourceResp["code"])
	if err != nil {
		socket.PublishStatusMessage("test_status", err.Error(), "error")
		return err, ""
	}
	return nil, ""
}
