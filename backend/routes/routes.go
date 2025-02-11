package routes

import (
	"github.com/Liphium/magic/backend/views"
	"github.com/gofiber/fiber/v2"
)

func InitializeRoutes(router fiber.Router) {

	// Serve the static files for the frontend
	router.Static("/static/", "./static")

	// Serve a basic first page
	router.Get("/", func(c *fiber.Ctx) error {
		return views.Render(c, views.Home())
	})
}
