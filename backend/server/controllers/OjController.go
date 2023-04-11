package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/skmonir/mango-gui/backend/judge-framework/services"
)

func Login(ctx *fiber.Ctx) error {
	credentials := struct {
		Platform      string `json:"platform"`
		HandleOrEmail string `json:"handleOrEmail"`
		Password      string `json:"password"`
	}{}
	err := ctx.BodyParser(&credentials)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err, handle := services.Login(credentials.Platform, credentials.HandleOrEmail, credentials.Password)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"handle":  handle,
	})
}

func SubmitCode(ctx *fiber.Ctx) error {
	platform := ctx.Params("platform")
	cid := ctx.Params("cid")
	label := ctx.Params("label")
	err, subId := services.Submit(platform, cid, label)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"subId":   subId,
	})
}
