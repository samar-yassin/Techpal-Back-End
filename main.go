package main

import (
	"CareerGuidance/database"
	"CareerGuidance/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	database.Connect()

	app := fiber.New()
	routes.Register(app)

	app.Listen(":8080")
}
