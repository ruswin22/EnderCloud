package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"endercloud-backend/internal/handlers"
)

func main() {
	app := fiber.New()

	app.Post("/create-server", handlers.CreateServer)

	log.Println("EnderCloud backend running on :8080")
	app.Listen(":8080")
}
