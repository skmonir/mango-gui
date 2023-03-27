package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/skmonir/mango-gui/backend/judge-framework/config"
)

func GetConfig(ctx *fiber.Ctx) error {
	judgeConfig := config.GetJudgeConfigFromCache()
	return ctx.Status(fiber.StatusOK).JSON(judgeConfig)
}

func UpdateConfig(ctx *fiber.Ctx) error {
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

func ResetConfig(ctx *fiber.Ctx) error {
	conf, err := config.CreateDefaultConfig()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON("Oops! Something went wrong!")
	}
	return ctx.Status(fiber.StatusOK).JSON(conf)
}
