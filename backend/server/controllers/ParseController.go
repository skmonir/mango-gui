package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/skmonir/mango-gui/backend/judge-framework/parser"
	"github.com/skmonir/mango-gui/backend/judge-framework/services"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
	"strings"
)

func Parse(ctx *fiber.Ctx) error {
	encodedUrl := ctx.Params("encoded_url")
	decodedUrl := utils.DecodeBase64(encodedUrl)
	services.UpdateParseUrlHistory(decodedUrl)
	parseResponseList := parser.Parse(decodedUrl)
	return ctx.Status(fiber.StatusOK).JSON(parseResponseList)
}

func ScheduleParse(ctx *fiber.Ctx) error {
	encodedUrl := ctx.Params("encoded_url")
	decodedUrl := utils.DecodeBase64(encodedUrl)
	services.UpdateParseUrlHistory(decodedUrl)
	startTime, err := parser.ScheduleParse(decodedUrl)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	startTimeStr := strings.TrimSpace(strings.Split(startTime.String(), "+")[0])
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":   true,
		"url":       decodedUrl,
		"startTime": startTimeStr,
	})
}
