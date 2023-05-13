package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/skmonir/mango-gui/backend/judge-framework/logger"
	"github.com/skmonir/mango-gui/backend/judge-framework/models"
	"github.com/skmonir/mango-gui/backend/judge-framework/services"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
)

func GetProblemList(ctx *fiber.Ctx) error {
	encodedUrl := ctx.Params("encoded_url")
	decodedUrl := utils.DecodeBase64(encodedUrl)
	services.UpdateTestContestUrlHistory(decodedUrl)
	problems := services.GetProblemListByUrl(utils.TransformUrl(decodedUrl))
	return ctx.Status(fiber.StatusOK).JSON(problems)
}

func GetProblem(ctx *fiber.Ctx) error {
	platform := ctx.Params("platform")
	cid := ctx.Params("cid")
	label := ctx.Params("label")
	prob := services.GetProblem(platform, cid, label)
	return ctx.Status(fiber.StatusOK).JSON(prob)
}

func AddCustomProblem(ctx *fiber.Ctx) error {
	problem := models.Problem{}
	err := ctx.BodyParser(&problem)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	logger.Info(fmt.Sprintf("Received request to add custom problem: %v", problem))
	problems := services.AddCustomProblem(problem)
	return ctx.Status(fiber.StatusOK).JSON(problems)
}
