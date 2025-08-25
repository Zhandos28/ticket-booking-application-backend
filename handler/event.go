package handler

import (
	"context"
	"strconv"
	"time"

	"github.com/Zhandos28/ticket-booking/model"
	"github.com/gofiber/fiber/v2"
)

type EventHandler struct {
	repository model.EventRepository
}

func (h *EventHandler) CreateOne(c *fiber.Ctx) error {
	event := &model.Event{}

	context, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	if err := c.BodyParser(&event); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	event, err := h.repository.CreateOne(context, event)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"data":    event,
	})
}

func (h *EventHandler) GetMany(c *fiber.Ctx) error {
	context, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	events, err := h.repository.GetMany(context)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"data":    events,
	})
}

func (h *EventHandler) GetOne(c *fiber.Ctx) error {
	eventID, err := strconv.Atoi(c.Params("eventID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	context, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	event, err := h.repository.GetOne(context, uint(eventID))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"data":    event,
	})

}

func (h *EventHandler) UpdateOne(c *fiber.Ctx) error {
	eventID, err := strconv.Atoi(c.Params("eventID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}
	context, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	updatedData := make(map[string]interface{})
	if err := c.BodyParser(&updatedData); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	event, err := h.repository.UpdateOne(context, uint(eventID), updatedData)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"data":    event,
	})
}

func (h *EventHandler) DeleteOne(c *fiber.Ctx) error {
	eventID, err := strconv.Atoi(c.Params("eventID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	context, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	if err := h.repository.DeleteOne(context, uint(eventID)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{})
}

func NewEventHandler(router fiber.Router, repository model.EventRepository) {
	handler := EventHandler{repository: repository}

	router.Get("/", handler.GetMany)
	router.Post("/", handler.CreateOne)
	router.Get("/:eventID", handler.GetOne)
	router.Put("/:eventID", handler.UpdateOne)
	router.Delete("/:eventID", handler.DeleteOne)
}
