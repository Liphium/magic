package auth_routes

import (
	"github.com/Liphium/magic/backend/views"
	auth_views "github.com/Liphium/magic/backend/views/auth"
	"github.com/gofiber/fiber/v2"
)

// Route: /auth/login
func loginRoute(c *fiber.Ctx) error {
	return views.RenderWithBase(c, auth_views.LoginView())
}
