package services

import (
	"main/common/dto"
	"main/config"
	conn "main/config/database"
	"main/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string, salt int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), salt)
	return string(bytes), err
}

// func comparePassword(hashedPassword string, password string) bool {
// 	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
// 	return err == nil
// }

func CreateUser(data dto.SignupUser) error {
	res := conn.Database.
		Where(&models.User{Email: data.Email}).
		First(&models.User{})

	if res.Error == nil {
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: "User already exists",
		}
	}

	env, _ := config.LoadEnvData()

	hashedPassword, err := hashPassword(data.Password, env.Salt)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: err.Error(),
		}
	}

	user := models.User{
		Email:    data.Email,
		Password: hashedPassword,
		Name:     data.Name,
	}

	if err := conn.Database.Create(&user).Error; err != nil {
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: err.Error(),
		}
	}

	return nil
}
