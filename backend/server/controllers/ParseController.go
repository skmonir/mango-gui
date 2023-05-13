package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/skmonir/mango-gui/backend/judge-framework/parser"
	"github.com/skmonir/mango-gui/backend/judge-framework/services"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
)

func Parse(ctx *fiber.Ctx) error {
	encodedUrl := ctx.Params("encoded_url")
	decodedUrl := utils.DecodeBase64(encodedUrl)
	services.UpdateParseUrlHistory(decodedUrl)
	parseResponseList := parser.Parse(utils.TransformUrl(decodedUrl))
	return ctx.Status(fiber.StatusOK).JSON(parseResponseList)
}
