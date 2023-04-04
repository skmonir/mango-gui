package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/skmonir/mango-gui/backend/socket"
	"log"
)

func RunServer() {
	go socket.RunSocketHub()

	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	SetRoutes(app.Group("/api/v1"))

	log.Fatal(app.Listen(":3456"))
}
