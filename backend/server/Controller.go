package server

import (
	"fmt"
	"github.com/skmonir/mango-gui/backend/socket"

	"github.com/gofiber/fiber/v2"
	"github.com/skmonir/mango-gui/backend/judge-framework/config"
	"github.com/skmonir/mango-gui/backend/judge-framework/fileService"
	"github.com/skmonir/mango-gui/backend/judge-framework/parser"
	"github.com/skmonir/mango-gui/backend/judge-framework/services"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
)

func parse(ctx *fiber.Ctx) error {
	encodedUrl := ctx.Params("encoded_url")
	parseResponseList := parser.Parse(utils.DecodeBase64(encodedUrl))
	return ctx.Status(fiber.StatusOK).JSON(parseResponseList)
}

func getProblemList(ctx *fiber.Ctx) error {
	encodedUrl := ctx.Params("encoded_url")
	problems := services.GetProblemListByUrl(utils.DecodeBase64(encodedUrl))
	return ctx.Status(fiber.StatusOK).JSON(problems)
}

func getProblem(ctx *fiber.Ctx) error {
	platform := ctx.Params("platform")
	cid := ctx.Params("cid")
	label := ctx.Params("label")
	prob := services.GetProblem(platform, cid, label)
	return ctx.Status(fiber.StatusOK).JSON(prob)
}

func getConfig(ctx *fiber.Ctx) error {
	judgeConfig := config.GetJudgeConfigFromCache()
	return ctx.Status(fiber.StatusOK).JSON(judgeConfig)
}

func updateConfig(ctx *fiber.Ctx) error {
	var configToUpdate config.JudgeConfig
	err := ctx.BodyParser(&configToUpdate)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	judgeConfig := config.UpdateJudgeConfigIntoCache(configToUpdate)
	return ctx.Status(fiber.StatusOK).JSON(judgeConfig)
}

func openSourceByPath(ctx *fiber.Ctx) error {
	openSourceRequest := struct {
		FilePath string `json:"filePath"`
	}{}
	err := ctx.BodyParser(&openSourceRequest)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	fmt.Println(openSourceRequest)
	fileService.OpenSourceByPath(openSourceRequest.FilePath)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}

func openSourceByMetadata(ctx *fiber.Ctx) error {
	platform := ctx.Params("platform")
	cid := ctx.Params("cid")
	label := ctx.Params("label")
	fileService.OpenSourceByMetadata(platform, cid, label)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}

func getCustomTestByPath(ctx *fiber.Ctx) error {
	getCustomTestRequest := struct {
		InputFilePath  string `json:"inputFilePath"`
		OutputFilePath string `json:"outputFilePath"`
	}{}
	err := ctx.BodyParser(&getCustomTestRequest)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	fmt.Println(getCustomTestRequest)
	testcase := services.GetTestcaseByPath(getCustomTestRequest.InputFilePath, getCustomTestRequest.OutputFilePath)
	return ctx.Status(fiber.StatusOK).JSON(testcase)
}

func addCustomTest(ctx *fiber.Ctx) error {
	addCustomTestRequest := struct {
		Platform  string `json:"platform"`
		ContestId string `json:"contestId"`
		Label     string `json:"label"`
		Input     string `json:"input"`
		Output    string `json:"output"`
	}{}
	err := ctx.BodyParser(&addCustomTestRequest)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	fmt.Println(addCustomTestRequest)
	services.SaveCustomTestcaseIntoFile(
		addCustomTestRequest.Platform,
		addCustomTestRequest.ContestId,
		addCustomTestRequest.Label,
		addCustomTestRequest.Input,
		addCustomTestRequest.Output)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}

func updateCustomTest(ctx *fiber.Ctx) error {
	updateCustomTestRequest := struct {
		Platform       string `json:"platform"`
		ContestId      string `json:"contestId"`
		Label          string `json:"label"`
		InputFilePath  string `json:"inputFilePath"`
		OutputFilePath string `json:"outputFilePath"`
		Input          string `json:"input"`
		Output         string `json:"output"`
	}{}
	err := ctx.BodyParser(&updateCustomTestRequest)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	fmt.Println(updateCustomTestRequest)
	services.UpdateCustomTestcaseIntoFile(
		updateCustomTestRequest.Platform,
		updateCustomTestRequest.ContestId,
		updateCustomTestRequest.Label,
		updateCustomTestRequest.InputFilePath,
		updateCustomTestRequest.OutputFilePath,
		updateCustomTestRequest.Input,
		updateCustomTestRequest.Output)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}

func deleteCustomTest(ctx *fiber.Ctx) error {
	updateCustomTestRequest := struct {
		Platform      string `json:"platform"`
		ContestId     string `json:"contestId"`
		Label         string `json:"label"`
		InputFilePath string `json:"inputFilePath"`
	}{}
	err := ctx.BodyParser(&updateCustomTestRequest)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	fmt.Println(updateCustomTestRequest)
	services.DeleteCustomTestcaseFromFile(
		updateCustomTestRequest.Platform,
		updateCustomTestRequest.ContestId,
		updateCustomTestRequest.Label,
		updateCustomTestRequest.InputFilePath)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}

func testProblem(ctx *fiber.Ctx) error {
	platform := ctx.Params("platform")
	cid := ctx.Params("cid")
	label := ctx.Params("label")
	probExecResult := services.RunTest(platform, cid, label)
	return ctx.Status(fiber.StatusOK).JSON(probExecResult)
}

func getExecutionResult(ctx *fiber.Ctx) error {
	platform := ctx.Params("platform")
	cid := ctx.Params("cid")
	label := ctx.Params("label")
	probExecResult := services.GetProblemExecutionResult(platform, cid, label, true, false)
	socket.PublishPreviousRunStatus(probExecResult)
	return ctx.Status(fiber.StatusOK).JSON(probExecResult)
}

func getInputOutputDirectoriesByUrl(ctx *fiber.Ctx) error {
	encodedUrl := ctx.Params("encoded_url")
	inputDirectory, outputDirectory := services.GetInputOutputDirectoryByUrl(utils.DecodeBase64(encodedUrl))
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"inputDirectory":  inputDirectory,
		"outputDirectory": outputDirectory,
	})
}
