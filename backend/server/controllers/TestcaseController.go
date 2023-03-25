package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/skmonir/mango-gui/backend/judge-framework/services"
)

func GetCustomTestByPath(ctx *fiber.Ctx) error {
	getCustomTestRequest := struct {
		InputFilePath  string `json:"inputFilePath"`
		OutputFilePath string `json:"outputFilePath"`
	}{}
	err := ctx.BodyParser(&getCustomTestRequest)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	fmt.Println(getCustomTestRequest)
	testcase := services.GetTestcaseByPath(getCustomTestRequest.InputFilePath, getCustomTestRequest.OutputFilePath)
	return ctx.Status(fiber.StatusOK).JSON(testcase)
}

func AddCustomTest(ctx *fiber.Ctx) error {
	addCustomTestRequest := struct {
		Platform  string `json:"platform"`
		ContestId string `json:"contestId"`
		Label     string `json:"label"`
		Input     string `json:"input"`
		Output    string `json:"output"`
	}{}
	err := ctx.BodyParser(&addCustomTestRequest)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	fmt.Println(addCustomTestRequest)
	services.SaveCustomTestcaseIntoFile(
		addCustomTestRequest.Platform,
		addCustomTestRequest.ContestId,
		addCustomTestRequest.Label,
		addCustomTestRequest.Input,
		addCustomTestRequest.Output)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}

func UpdateCustomTest(ctx *fiber.Ctx) error {
	updateCustomTestRequest := struct {
		Platform       string `json:"platform"`
		ContestId      string `json:"contestId"`
		Label          string `json:"label"`
		InputFilePath  string `json:"inputFilePath"`
		OutputFilePath string `json:"outputFilePath"`
		Input          string `json:"input"`
		Output         string `json:"output"`
	}{}
	err := ctx.BodyParser(&updateCustomTestRequest)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	fmt.Println(updateCustomTestRequest)
	services.UpdateCustomTestcaseIntoFile(
		updateCustomTestRequest.Platform,
		updateCustomTestRequest.ContestId,
		updateCustomTestRequest.Label,
		updateCustomTestRequest.InputFilePath,
		updateCustomTestRequest.OutputFilePath,
		updateCustomTestRequest.Input,
		updateCustomTestRequest.Output)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}

func DeleteCustomTest(ctx *fiber.Ctx) error {
	updateCustomTestRequest := struct {
		Platform      string `json:"platform"`
		ContestId     string `json:"contestId"`
		Label         string `json:"label"`
		InputFilePath string `json:"inputFilePath"`
	}{}
	err := ctx.BodyParser(&updateCustomTestRequest)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	fmt.Println(updateCustomTestRequest)
	services.DeleteCustomTestcaseFromFile(
		updateCustomTestRequest.Platform,
		updateCustomTestRequest.ContestId,
		updateCustomTestRequest.Label,
		updateCustomTestRequest.InputFilePath)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}
