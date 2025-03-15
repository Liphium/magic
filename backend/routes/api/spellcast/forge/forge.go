package spellcast_forge_routes

import (
	"github.com/Liphium/magic/backend/database"
	"github.com/Liphium/magic/backend/util/constants"
	"github.com/gofiber/fiber/v2"
)

func Unauthorized(router fiber.Router) {

	// A middleware to make sure the spellcast token is valid
	router.Use(func(c *fiber.Ctx) error {

		// Get the token from the header
		token := c.Get("SC-Token", "-")
		if token == "-" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		// Get the build from the token to verify it
		var build database.Build
		if err := database.DBConn.Where("spellcast_token = ? AND status = ?", token, database.BuildStatusStarting).Take(&build).Error; err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		// Add it to the locals
		c.Locals(constants.LocalsForgeBuild, build)

		// Continue to the next handler in case the token is valid
		return c.Next()
	})

	router.Post("/connect", acceptConnection)
}
