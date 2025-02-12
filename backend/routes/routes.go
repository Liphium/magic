package routes

import (
	"os"

	auth_routes "github.com/Liphium/magic/backend/routes/auth"
	panel_routes "github.com/Liphium/magic/backend/routes/panel"
	"github.com/gofiber/fiber/v2"
)

func InitializeRoutes(router fiber.Router) {

	// Serve the static files for the frontend
	router.Static("/static/", "./static")

	// Add the skip auth endpoint (only when enabled)
	if os.Getenv("MAGIC_SKIP_AUTH") == "true" {
		router.Get("_dev/skip", skipAuth)
	}

	// Add all the routes
	router.Route("/", unauthorizedRouter)
	router.Route("/a", authorizedRouter)
}

func unauthorizedRouter(router fiber.Router) {
	router.Route("/auth", auth_routes.Unauthorized)
}

func authorizedRouter(router fiber.Router) {
	// TODO: Add auth middleware

	router.Route("/panel", panel_routes.Authorized)
}
