package forge_routes

import (
	"github.com/Liphium/magic/backend/database"
	"github.com/Liphium/magic/backend/views"
	forge_views "github.com/Liphium/magic/backend/views/panel/forge"
	"github.com/gofiber/fiber/v2"
)

// Route: /a/panel/forge/:id/builds/:build/cancel
func cancelBuild(c *fiber.Ctx) error {

	// Get all of the base stuff for later
	forge, _, valid := getBaseInfo(c)
	if !valid {
		return views.RenderError(c, forge_views.BuildViewCancelMessage(true, "This forge doesn't exist."), nil)
	}

	// Check if there is a branch path param
	buildId := c.Params("build", "")
	if buildId == "" {
		return views.RenderError(c, forge_views.BuildViewCancelMessage(true, "This build doesn't exist."), nil)
	}

	// Delete the build
	if err := database.DBConn.Where("forge = ? AND id = ?", forge.ID.String(), buildId).Delete(&database.Build{}).Error; err != nil {
		return views.RenderError(c, forge_views.BuildViewCancelMessage(true, "Build couldn't be deleted."), err)
	}

	// Redirect the user to the panel
	c.Set("HX-Redirect", "/a/panel/forge/"+forge.ID.String()+"/builds")

	return views.RenderJust(c, forge_views.BuildViewCancelMessage(false, "Redirecting to builds.."))
}
