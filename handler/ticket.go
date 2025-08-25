package handler

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/skip2/go-qrcode"

	"github.com/Zhandos28/ticket-booking/model"
	"github.com/Zhandos28/ticket-booking/repository"
	"github.com/gofiber/fiber/v2"
)

type TicketHandler struct {
	repository *repository.TicketRepository
}

func (h *TicketHandler) GetMany(c *fiber.Ctx) error {
	context, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	userID := uint(c.Locals("userID").(float64))

	tickets, err := h.repository.GetMany(context, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"data":    tickets,
	})
}

func (h *TicketHandler) GetOne(c *fiber.Ctx) error {
	ticketID, err := strconv.Atoi(c.Params("ticketID"))
	if err != nil {

	}

	userID := uint(c.Locals("userID").(float64))
	context, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	ticket, err := h.repository.GetOne(context, userID, uint(ticketID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	var QRCode []byte
	QRCode, err = qrcode.Encode(
		fmt.Sprintf("ticketID:%v,ownerID:%v", ticketID, userID),
		qrcode.Medium,
		256,
	)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"data": fiber.Map{
			"ticket": ticket,
			"qrcode": QRCode,
		},
	})
}

func (h *TicketHandler) CreateOne(c *fiber.Ctx) error {
	ticket := &model.Ticket{}

	if err := c.BodyParser(&ticket); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	userID := uint(c.Locals("userID").(float64))
	context, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	ticket, err := h.repository.CreateOne(context, userID, ticket)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"data":    ticket,
	})
}

func (h *TicketHandler) ValidateOne(c *fiber.Ctx) error {
	validateTicket := model.ValidateTicket{}
	if err := c.BodyParser(&validateTicket); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	validateData := make(map[string]interface{})
	validateData["entered"] = true

	context, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	ticket, err := h.repository.UpdateOne(context, validateTicket.OwnerID, validateTicket.TicketID, validateData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"data":    ticket,
	})
}
func NewTicketHandler(router fiber.Router, repository *repository.TicketRepository) {
	handler := TicketHandler{repository: repository}

	router.Get("/", handler.GetMany)
	router.Get("/:ticketID", handler.GetOne)
	router.Post("/", handler.CreateOne)
	router.Post("/validate", handler.ValidateOne)
}
