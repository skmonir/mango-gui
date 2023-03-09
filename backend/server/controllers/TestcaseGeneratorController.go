package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/skmonir/mango-gui/backend/judge-framework/dto"
	"github.com/skmonir/mango-gui/backend/judge-framework/testcaseGenerator"
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
	execRes := testcaseGenerator.Generate(req)
	return ctx.Status(fiber.StatusOK).JSON(execRes)
}
