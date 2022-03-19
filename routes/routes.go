package routes

import (
	"CareerGuidance/controllers"
	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App) {
	app.Get("/", controllers.Register)
}
