package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/skmonir/mango-gui/backend/judge-framework/parser"
	"github.com/skmonir/mango-gui/backend/judge-framework/services"
)

func GetAppData(ctx *fiber.Ctx) error {
	appData := services.GetAppData()
	for _, scheduledTask := range appData.ParseSchedulerTasks {
		if scheduledTask.Stage == "SCHEDULED" {
			_ = parser.ScheduleTaskInScheduler(scheduledTask)
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(appData)
}
