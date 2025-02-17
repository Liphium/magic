package wizard_routes

import (
	"github.com/Liphium/magic/backend/database"
	"github.com/Liphium/magic/backend/util/constants"
	"github.com/gofiber/fiber/v2"
)

func Unauthorized(router fiber.Router) {

	// All the routes that wizard uses to interact with Magic
	router.Route("/api", wizardAPI)
}

func wizardAPI(router fiber.Router) {
	// A middleware to make sure the wizard token is valid
	router.Use(func(c *fiber.Ctx) error {

		// Get the token from the header
		token := c.Get("W-Token", "-")
		if token == "-" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		// Validate the token using the database
		var wizard database.Wizard
		if err := database.DBConn.Where("token = ?", token).Take(&wizard).Error; err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// Put the wizard into locals for the endpoints
		c.Locals(constants.LocalsWizard, wizard)

		// Continue to the next handler in case the token is valid
		return c.Next()
	})
}
