package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/skmonir/mango-gui/backend/judge-framework/dto"
	"github.com/skmonir/mango-gui/backend/judge-framework/services"
	"github.com/skmonir/mango-gui/backend/judge-framework/services/testcaseGeneratorServices"
	"strings"
)

func GenerateRandomTests(ctx *fiber.Ctx) error {
	req := dto.TestcaseGenerateRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	fmt.Println(req)
	req.ParsedProblemUrl = strings.ToLower(req.ParsedProblemUrl)
	services.UpdateInputGenerateRequestHistory(req)
	execRes := testcaseGeneratorServices.GenerateInput(req)
	return ctx.Status(fiber.StatusOK).JSON(execRes)
}

func GenerateOutputs(ctx *fiber.Ctx) error {
	req := dto.TestcaseGenerateRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	fmt.Println(req)
	services.UpdateOutputGenerateRequestHistory(req)
	req.ParsedProblemUrl = strings.ToLower(req.ParsedProblemUrl)
	execRes := testcaseGeneratorServices.GenerateOutput(req)
	return ctx.Status(fiber.StatusOK).JSON(execRes)
}
