package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/skmonir/mango-gui/backend/judge-framework/parser"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
)

func Parse(ctx *fiber.Ctx) error {
	encodedUrl := ctx.Params("encoded_url")
	parseResponseList := parser.Parse(utils.DecodeBase64(encodedUrl))
	return ctx.Status(fiber.StatusOK).JSON(parseResponseList)
}
