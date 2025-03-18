package middleware

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fevzisahinler/Appointfy-Backend/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware() fiber.Handler {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Failed to load configuration: ", err)
	}

	secret := cfg.JwtSecretKey
	if secret == "" {
		fmt.Println("JWT secret key is missing in configuration")
	}

	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			fmt.Println("Missing authorization header")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Missing authorization header",
			})
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			fmt.Println("Invalid authorization header format")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Authorization header must start with 'Bearer '",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				fmt.Println("Unexpected signing method")
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrTokenMalformed) {
				fmt.Println("That's not even a token")
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "Invalid token format",
				})
			} else if errors.Is(err, jwt.ErrTokenExpired) {
				fmt.Println("Token has expired")
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "Token has expired",
				})
			} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
				fmt.Println("Token not valid yet")
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "Token not valid yet",
				})
			} else {
				fmt.Println("Couldn't handle this token")
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "Invalid or expired token",
					"error":   err.Error(),
				})
			}
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			fmt.Println("Failed to parse token claims")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Failed to parse token claims",
			})
		}

		c.Locals("user", claims)
		return c.Next()
	}
}
