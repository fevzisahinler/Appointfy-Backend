package controllers

import (
	"errors"
	"fmt"
	"github.com/fevzisahinler/Appointfy-Backend/providers/cryptology"

	"github.com/fevzisahinler/Appointfy-Backend/db"
	"github.com/fevzisahinler/Appointfy-Backend/http/requests"
	"github.com/fevzisahinler/Appointfy-Backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Register(c *fiber.Ctx) error {
	requestUser := new(requests.CreateUserRequest)
	if err := c.BodyParser(requestUser); err != nil {
		fmt.Println("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid input",
		})
	}

	if err := requestUser.Validate(); err != nil {
		fmt.Println("Validation failed")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
		})
	}

	var existingUser models.User
	if err := db.DB.Where("username = ? OR email = ?", requestUser.Username, requestUser.Email).
		First(&existingUser).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("Failed to check existing user:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An unexpected error occurred",
		})
	}

	if existingUser.ID != 0 {
		fmt.Println("User already exists")
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "User already exists",
		})
	}

	hashedPassword, err := cryptology.HashPassword(requestUser.Password)
	if err != nil {
		fmt.Println("Failed to hash password:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An unexpected error occurred",
		})
	}

	userModel := &models.User{
		Username:    requestUser.Username,
		Password:    hashedPassword,
		Name:        requestUser.Name,
		Surname:     requestUser.Surname,
		Email:       requestUser.Email,
		PhoneNumber: requestUser.PhoneNumber,
	}

	if err := db.DB.Create(&userModel).Error; err != nil {
		fmt.Println("Failed to create user:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An unexpected error occurred",
		})
	}

	var defaultRole models.Role
	if err := db.DB.Where("role_name = ?", "Default").First(&defaultRole).Error; err != nil {
		fmt.Println("Could not find 'Default' role, user has no default role assigned")
	} else {
		if errAssoc := db.DB.Model(&userModel).Association("Roles").Append(&defaultRole); errAssoc != nil {
			fmt.Println("Failed to associate user with default role:", errAssoc)
		} else {
			fmt.Println("User associated with 'Default' role successfully")
		}
	}

	fmt.Println("User created successfully:", userModel.Username)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"data":    userModel,
	})
}
