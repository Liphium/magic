package routes

import (
	"os"

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

	// Add the skip auth endpoint (only when enabled)
	if os.Getenv("MAGIC_SKIP_AUTH") == "true" {
		router.Get("_dev/skip", skipAuth)
	}

	// Add all authorized routes
	router.Route("/a", authorizedRouter)
}

func authorizedRouter(router fiber.Router) {

}
