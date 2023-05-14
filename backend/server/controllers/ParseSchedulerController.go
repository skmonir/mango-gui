package controllers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/skmonir/mango-gui/backend/judge-framework/parser"
	"github.com/skmonir/mango-gui/backend/judge-framework/services"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
)

func ScheduleParse(ctx *fiber.Ctx) error {
	scheduleParseRequest := struct {
		Url string `json:"url"`
	}{}
	err := ctx.BodyParser(&scheduleParseRequest)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	scheduleParseRequest.Url = utils.TransformUrl(strings.ToLower(scheduleParseRequest.Url))

	services.UpdateParseUrlHistory(scheduleParseRequest.Url)

	err = parser.ScheduleParse(scheduleParseRequest.Url)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}

func GetParseScheduledTasks(ctx *fiber.Ctx) error {
	appData := services.GetAppData()
	return ctx.Status(fiber.StatusOK).JSON(appData.ParseSchedulerTasks)
}

func RemoveParseScheduledTask(ctx *fiber.Ctx) error {
	taskId := ctx.Params("id")
	parser.RemoveParseSchedule(taskId)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}
