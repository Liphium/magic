package routes

import "github.com/gofiber/fiber/v2"

// Route: /a/_dev/skip
func skipAuth(c *fiber.Ctx) error {

	c.Cookie(&fiber.Cookie{
		Name:  "auth",
		Value: "admin",
	})

	return c.SendString("Wrote cookie.")
}
