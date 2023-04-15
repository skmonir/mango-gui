package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/skmonir/mango-gui/backend/judge-framework/services"
	"github.com/skmonir/mango-gui/backend/socket"
)

func TestProblem(ctx *fiber.Ctx) error {
	platform := ctx.Params("platform")
	cid := ctx.Params("cid")
	label := ctx.Params("label")
	probExecResult := services.RunTest(platform, cid, label)
	return ctx.Status(fiber.StatusOK).JSON(probExecResult)
}

func GetExecutionResult(ctx *fiber.Ctx) error {
	platform := ctx.Params("platform")
	cid := ctx.Params("cid")
	label := ctx.Params("label")
	probExecResult := services.GetProblemExecutionResult(platform, cid, label, true, false)
	socket.PublishPreviousRunStatus(&probExecResult)
	return ctx.Status(fiber.StatusOK).JSON(probExecResult)
}
