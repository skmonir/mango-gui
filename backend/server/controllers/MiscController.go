package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/skmonir/mango-gui/backend/judge-framework/services"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
)

func GetInputOutputDirectoriesByUrl(ctx *fiber.Ctx) error {
	encodedUrl := ctx.Params("encoded_url")
	inputDirectory, outputDirectory := services.GetInputOutputDirectoryByUrl(utils.DecodeBase64(encodedUrl))
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"inputDirectory":  inputDirectory,
		"outputDirectory": outputDirectory,
	})
}

func CheckDirectoryPathValidity(ctx *fiber.Ctx) error {
	encodedPath := ctx.Params("encoded_path")
	isExist := utils.IsDirExist(utils.DecodeBase64(encodedPath))
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"isExist": isExist,
	})
}

func CheckFilePathValidity(ctx *fiber.Ctx) error {
	encodedPath := ctx.Params("encoded_path")
	isExist := utils.IsFileExist(utils.DecodeBase64(encodedPath))
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"isExist": isExist,
	})
}
