package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/Zhandos28/ticket-booking/middleware"
	"github.com/Zhandos28/ticket-booking/service"

	"github.com/Zhandos28/ticket-booking/config"
	"github.com/Zhandos28/ticket-booking/db"
	"github.com/Zhandos28/ticket-booking/handler"
	"github.com/Zhandos28/ticket-booking/repository"
	"github.com/gofiber/fiber/v2"
)

func main() {
	envConfig := config.NewEnvConfig()

	database := db.Init(envConfig, db.Migrator)
	app := fiber.New(fiber.Config{
		AppName:      "TicketBooking",
		ServerHeader: "Fiber",
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	eventRepository := repository.NewEventRepository(database)
	ticketRepository := repository.NewTicketRepository(database)
	authRepository := repository.NewAuthRepository(database)

	// Service
	authService := service.NewAuthService(authRepository, envConfig)

	server := app.Group("/api")
	handler.NewAuthHandler(server.Group("/auth"), authService)

	privateRoutes := server.Use(middleware.AuthProtected(database, envConfig))

	handler.NewEventHandler(privateRoutes.Group("/events"), eventRepository)
	handler.NewTicketHandler(privateRoutes.Group("/tickets"), ticketRepository)

	app.Listen(fmt.Sprintf(":%s", envConfig.ServerPort))
}
