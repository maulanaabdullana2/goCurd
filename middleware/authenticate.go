package middleware

import (
	"fiber-crud/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")

		if tokenString == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "No token provided")
		}

		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		fmt.Println("Received token:", tokenString)

		claims, err := utils.ParseTokenString(tokenString)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
		}

		c.Locals("userID", claims.Subject)
		return c.Next()
	}
}
