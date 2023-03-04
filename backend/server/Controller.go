package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/skmonir/mango-ui/backend/judge-framework/config"
	"github.com/skmonir/mango-ui/backend/judge-framework/fileService"
	"github.com/skmonir/mango-ui/backend/judge-framework/parser"
	"github.com/skmonir/mango-ui/backend/judge-framework/services"
	"github.com/skmonir/mango-ui/backend/judge-framework/utils"
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

func getCodeByPath(ctx *fiber.Ctx) error {
	codeRequest := struct {
		FilePath string `json:"filePath"`
	}{}
	err := ctx.BodyParser(&codeRequest)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	fmt.Println(codeRequest)
	code := utils.ReadFileContent(codeRequest.FilePath, 123456, 123456)
	return ctx.Status(fiber.StatusOK).JSON(code)
}

func getCodeByMetadata(ctx *fiber.Ctx) error {
	platform := ctx.Params("platform")
	cid := ctx.Params("cid")
	label := ctx.Params("label")
	code := fileService.GetCodeByMetadata(platform, cid, label)
	return ctx.Status(fiber.StatusOK).JSON(code)
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
	probExecResult := services.GetProblemExecutionResult(platform, cid, label, true)
	return ctx.Status(fiber.StatusOK).JSON(probExecResult)
}
