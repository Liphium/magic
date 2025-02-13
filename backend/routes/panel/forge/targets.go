package forge_routes

import (
	"github.com/Liphium/magic/backend/database"
	"github.com/Liphium/magic/backend/views"
	panel_views "github.com/Liphium/magic/backend/views/panel"
	forge_views "github.com/Liphium/magic/backend/views/panel/forge"
	"github.com/gofiber/fiber/v2"
)

// Route: /a/panel/forge/:id
func targetsPage(c *fiber.Ctx) error {

	// Get all of the base stuff for later
	forge, sidebar, err := getBaseInfo(c)
	if err != nil {
		return c.Redirect("/a/panel", fiber.StatusPermanentRedirect)
	}

	// Create the page and the sidebar based on the Forge
	page := panel_views.PanelPage(forge.Label, forge_views.TargetsPage([]database.Target{
		{
			Type:  "Pull Request Target",
			Value: "Builds the app & deploys to Preview",
		},
	}))

	return views.RenderHTMX(c, panel_views.Base(sidebar, page), page)
}
