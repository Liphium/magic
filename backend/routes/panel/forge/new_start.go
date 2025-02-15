package forge_routes

import (
	"github.com/Liphium/magic/backend/views"
	panel_views "github.com/Liphium/magic/backend/views/panel"
	forge_views "github.com/Liphium/magic/backend/views/panel/forge"
	"github.com/gofiber/fiber/v2"
)

// Route: /a/panel/forge/new
func newForgeStart(c *fiber.Ctx) error {

	// Render all the forges
	stepPage := forge_views.NewForgeStep1()
	panelPage := panel_views.PanelPage("1. Select a repository", stepPage)
	sidebar := forge_views.NewForgeSidebar()

	return views.RenderHTMX(c, panel_views.Base(sidebar, panelPage), panelPage)

}
