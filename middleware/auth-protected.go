package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Zhandos28/ticket-booking/config"
	"github.com/Zhandos28/ticket-booking/model"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func AuthProtected(db *gorm.DB, envConfig *config.EnvConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() == http.MethodOptions {
			return c.Next()
		}

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			log.Warnf("No Authorization header found")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "failed",
				"message": "Authorization header is required",
			})
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			log.Warnf("Authorization header is invalid")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "failed",
				"message": "Authorization header is invalid",
			})
		}

		tokenStr := tokenParts[1]
		secret := []byte(envConfig.JWTSecret)

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.GetSigningMethod("HS256").Alg() {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return secret, nil
		})

		if err != nil || !token.Valid {
			log.Warnf("Invalid token")

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "failed",
				"message": "Unauthorized",
			})
		}

		userID := token.Claims.(jwt.MapClaims)["id"]

		if err := db.Model(&model.User{}).First("id = ?", userID).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warnf("User not found in database")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "failed",
				"message": "Unauthorized",
			})
		}

		c.Locals("userID", userID)
		return c.Next()
	}
}
