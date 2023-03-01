package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/skmonir/mango-ui/backend/judge-framework/config"
	"github.com/skmonir/mango-ui/backend/judge-framework/fileService"
	"github.com/skmonir/mango-ui/backend/judge-framework/parser"
	"github.com/skmonir/mango-ui/backend/judge-framework/services"
	"github.com/skmonir/mango-ui/backend/judge-framework/utils"
	"github.com/skmonir/mango-ui/backend/socket"
)

func SetRoutes(app fiber.Router) {
	app.Get("/ws", websocket.New(func(conn *websocket.Conn) {
		socket.CreateNewSocketConnection(conn)
	}))

	app.Get("/parse/:encoded_url", parse)

	app.Get("/problem/:encoded_url", getProblemList)

	app.Get("/problem/:platform/:cid/:label", getProblem)

	app.Get("/config", getConfig)

	app.Put("/config", updateConfig)

	app.Put("/code", getCode)

	app.Put("/source", openSource)

	app.Get("/test/:platform/:cid/:label", testProblem)

	app.Get("/execresult/:platform/:cid/:label", getExecutionResult)
}

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

func getCode(ctx *fiber.Ctx) error {
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

func openSource(ctx *fiber.Ctx) error {
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
	fileService.Open(openSourceRequest.FilePath)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}

func testProblem(ctx *fiber.Ctx) error {
	//platform := ctx.Params("platform")
	//cid := ctx.Params("cid")
	//label := ctx.Params("label")
	//prob := tester.RunTest(platform, cid, label) // change
	return ctx.Status(fiber.StatusOK).JSON("")
}

func getExecutionResult(ctx *fiber.Ctx) error {
	platform := ctx.Params("platform")
	cid := ctx.Params("cid")
	label := ctx.Params("label")
	fileService.GetTestcasesFromFile(platform, cid, label)
	return ctx.Status(fiber.StatusOK).JSON("")
}
