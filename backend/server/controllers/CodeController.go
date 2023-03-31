package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/skmonir/mango-gui/backend/judge-framework/services/fileServices"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
	"path/filepath"
)

func GetCodeByFilePath(ctx *fiber.Ctx) error {
	codeRequest := struct {
		FilePath string `json:"filePath"`
	}{}
	err := ctx.BodyParser(&codeRequest)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	fmt.Println(codeRequest)
	code := utils.ReadFileContent(codeRequest.FilePath, 123456, 123456)
	response := map[string]string{
		"lang": utils.GetLangNameByFileExt(filepath.Ext(codeRequest.FilePath)),
		"code": code,
	}
	fmt.Println(response)
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func GetCodeByProblemPath(ctx *fiber.Ctx) error {
	platform := ctx.Params("platform")
	cid := ctx.Params("cid")
	label := ctx.Params("label")
	response := fileServices.GetCodeByMetadata(platform, cid, label)
	fmt.Println(response)
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func UpdateCodeByFilePath(ctx *fiber.Ctx) error {
	updateRequest := struct {
		FilePath string `json:"filePath"`
		Code     string `json:"code"`
	}{}
	err := ctx.BodyParser(&updateRequest)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	directory, filename := filepath.Split(updateRequest.FilePath)
	utils.WriteFileContent(directory, filename, []byte(updateRequest.Code))
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}

func UpdateCodeByProblemPath(ctx *fiber.Ctx) error {
	platform := ctx.Params("platform")
	cid := ctx.Params("cid")
	label := ctx.Params("label")

	updateRequest := struct {
		Code string `json:"code"`
	}{}
	err := ctx.BodyParser(&updateRequest)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	fileServices.UpdateCodeByProblemPath(platform, cid, label, updateRequest.Code)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}
