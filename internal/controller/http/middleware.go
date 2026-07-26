package http

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JWTMiddleware(secretKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if (authHeader == "") || (!strings.HasPrefix(authHeader, "Bearer ")){
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing or malformed jwt"})
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token != secretKey {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}
		return c.Next()		
	}
}