package forge_routes

import (
	"github.com/Liphium/magic/backend/database"
	"github.com/Liphium/magic/backend/views"
	form_views "github.com/Liphium/magic/backend/views/forms"
	"github.com/gofiber/fiber/v2"
)

// Route: /a/panel/forge/:id/delete
func deleteForge(c *fiber.Ctx) error {

	// Get the forge that should be deleted
	forge, _, valid := getBaseInfo(c)
	if !valid {
		return views.RenderJust(c, form_views.FormSubmitError("This Forge couldn't be found."))
	}

	// TODO: Proper deletion of all things the Forge has
	if err := database.DBConn.Delete(&forge).Error; err != nil {
		return views.RenderJust(c, form_views.FormSubmitError("This Forge couldn't be deleted."))
	}

	// Redirect the user to the panel
	c.Set("HX-Redirect", "/a/panel/forge")

	return views.RenderJust(c, form_views.FormSubmitSuccess("The Forge has been deleted."))

}
