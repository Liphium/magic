package panel_routes

import (
	"github.com/Liphium/magic/backend/views"
	panel_views "github.com/Liphium/magic/backend/views/panel"
	"github.com/gofiber/fiber/v2"
)

func Authorized(router fiber.Router) {
	router.Get("/", baseRoute)
}

// Route: /a/panel
func baseRoute(c *fiber.Ctx) error {
	return views.Render(c, panel_views.Base())
}
