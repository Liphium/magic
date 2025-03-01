package panel_routes

import (
	"github.com/Liphium/magic/backend/util/constants"
	"github.com/gofiber/fiber/v2"
)

// Route: /a/panel/logout
func logout(c *fiber.Ctx) error {

	// Clear the magic session cookie
	c.Cookie(&fiber.Cookie{
		Name:  constants.CookieMagicSession,
		Value: "",
	})

	return c.Redirect("/auth/login", fiber.StatusTemporaryRedirect)
}
