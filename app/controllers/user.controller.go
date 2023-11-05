package controllers

import (
	"fmt"
	"main/common/dto"
	"main/models"
	"main/services"

	"github.com/gofiber/fiber/v2"
)

func FetchUser(ctx *fiber.Ctx) error {
	auth := ctx.Locals("auth").(*models.User)

	if auth == nil {
		return &fiber.Error{
			Code:    fiber.ErrUnauthorized.Code,
			Message: "User not found",
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"Email":     auth.Email,
		"ID":        auth.ID,
		"Name":      auth.Name,
		"CreatedAt": auth.CreatedAt,
		"UpdatedAt": auth.UpdatedAt,
	})
}

func UpdateUser(ctx *fiber.Ctx) error {
	auth := ctx.Locals("auth").(*models.User)

	if auth == nil {
		return &fiber.Error{
			Code:    fiber.ErrUnauthorized.Code,
			Message: "User not found",
		}
	}

	updatedUser := ctx.Locals("json").(*dto.UpdateUser)

	err := services.UpdateUser(auth.ID, dto.UpdateUser{Name: updatedUser.Name})

	if err != nil {
		return &fiber.Error{
			Code:    fiber.ErrInternalServerError.Code,
			Message: "User update failed",
		}
	}

	fmt.Println(updatedUser)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User updated successfully",
	})
}