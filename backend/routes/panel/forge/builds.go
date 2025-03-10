package forge_routes

import (
	"github.com/Liphium/magic/backend/database"
	"github.com/Liphium/magic/backend/views"
	panel_views "github.com/Liphium/magic/backend/views/panel"
	forge_views "github.com/Liphium/magic/backend/views/panel/forge"
	"github.com/gofiber/fiber/v2"
)

// Route: /a/panel/forge/:id/builds
func buildListPage(c *fiber.Ctx) error {

	// Get all of the base stuff for later
	forge, sidebar, valid := getBaseInfo(c)
	if !valid {
		return c.Redirect("/a/panel/forge", fiber.StatusPermanentRedirect)
	}

	// Get all the builds for the Forge
	var builds []database.Build
	if err := database.DBConn.Where("forge = ?", forge.ID).Limit(25).Find(&builds).Error; err != nil {
		return panel_views.RenderPanelError(c, "Something went wrong on the server.", err)
	}

	// Render all the builds
	forgePage := forge_views.BuildsListPage(forge, builds)
	panelPage := panel_views.PanelPage("Builds", forgePage)

	return views.RenderHTMX(c, panel_views.Base(sidebar, panelPage), panelPage)
}
