package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/skmonir/mango-gui/backend/judge-framework/fileService"
	"github.com/skmonir/mango-gui/backend/judge-framework/services"
)

func OpenSourceByPath(ctx *fiber.Ctx) error {
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

func OpenSourceByMetadata(ctx *fiber.Ctx) error {
	platform := ctx.Params("platform")
	cid := ctx.Params("cid")
	label := ctx.Params("label")
	if err := fileService.OpenSourceByMetadata(platform, cid, label); err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(err.Error())
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}

func GenerateSourceByProblemPath(ctx *fiber.Ctx) error {
	platform := ctx.Params("platform")
	cid := ctx.Params("cid")
	label := ctx.Params("label")

	if problem := services.GetProblem(platform, cid, label); problem.Status == "success" {
		fileService.GenerateSourceByProblemPath(problem)
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "success",
		})
	}
	return ctx.Status(fiber.StatusNotFound).JSON("Problem not found")
}
