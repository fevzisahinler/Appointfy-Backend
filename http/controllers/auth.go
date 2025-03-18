package controllers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/fevzisahinler/Appointfy-Backend/config"
	"github.com/fevzisahinler/Appointfy-Backend/db"
	"github.com/fevzisahinler/Appointfy-Backend/models"
	"github.com/fevzisahinler/Appointfy-Backend/providers/cryptology"
)

func Login(c *fiber.Ctx) error {
	type LoginRequest struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	loginReq := new(LoginRequest)
	if err := c.BodyParser(loginReq); err != nil {
		fmt.Println("Failed to parse login request")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid input",
		})
	}

	var user models.User
	if err := db.DB.Where("username = ?", loginReq.Username).First(&user).Error; err != nil {
		fmt.Println("User not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid credentials",
		})
	}

	if err := cryptology.CheckPasswordHash(loginReq.Password, user.Password); err != nil {
		fmt.Println("Invalid password")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid credentials",
		})
	}

	cfg, _ := config.LoadConfig()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()

	tokenString, err := token.SignedString([]byte(cfg.JwtSecretKey))
	if err != nil {
		fmt.Println("Failed to generate token")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Could not login",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": fmt.Sprintf("Welcome, %s! Login successful.", user.Username),
		"data": fiber.Map{
			"token": tokenString,
		},
	})
}
