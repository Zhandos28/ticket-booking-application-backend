package handler

import (
	"context"
	"errors"
	"time"

	"github.com/Zhandos28/ticket-booking/model"
	"github.com/Zhandos28/ticket-booking/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

type AuthHandler struct {
	service model.AuthService
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	creds := model.AuthCredentials{}

	context, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	if err := c.BodyParser(&creds); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	if err := validate.Struct(creds); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	token, user, err := h.service.Login(context, &creds)
	if err != nil {
		if errors.Is(err, service.InvalidCredentialsError) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "failed",
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully logged in",
		"data": fiber.Map{
			"token": token,
			"user":  user,
		},
	})

}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	creds := model.AuthCredentials{}

	context, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	if err := c.BodyParser(&creds); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	if err := validate.Struct(creds); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "failed",
			"message": "please provide a valid email or password",
		})
	}

	token, user, err := h.service.Register(context, &creds)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"data": fiber.Map{
			"token": token,
			"user":  user,
		},
	})
}

func NewAuthHandler(router fiber.Router, service model.AuthService) *AuthHandler {
	handler := &AuthHandler{service: service}

	router.Post("/login", handler.Login)
	router.Post("/register", handler.Register)

	return handler
}
