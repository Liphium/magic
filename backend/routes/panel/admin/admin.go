package admin_panel_routes

import (
	"github.com/Liphium/magic/backend/util/constants"
	"github.com/gofiber/fiber/v2"
)

func Authorized(router fiber.Router) {

	// Make sure to not let non admins through
	router.Use(func(c *fiber.Ctx) error {

		// Check to make sure the user has admin permissions
		if c.Locals(constants.LocalsPermissionLevel).(uint) < constants.PermissionAdmin {
			return c.Redirect("/a/panel", fiber.StatusTemporaryRedirect)
		}

		// Go to next handler
		return c.Next()
	})

	router.Get("/wizards", wizardListPage)
	router.Get("/demo", demoRoute)
	router.Get("/demo/progress", demoProgress)
}
