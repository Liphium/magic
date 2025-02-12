package forge_routes

import (
	"github.com/Liphium/magic/backend/database"
	"github.com/Liphium/magic/backend/views"
	panel_views "github.com/Liphium/magic/backend/views/panel"
	"github.com/gofiber/fiber/v2"
)

func Authorized(router fiber.Router) {
	router.Get("/", baseRoute)
}

// Route: /a/panel/forge
func baseRoute(c *fiber.Ctx) error {

	forgePage := panel_views.ForgePage([]database.Forge{
		{
			Label:      "Liphium Chat",
			Repository: "https://github.com/Liphium/chat_interface",
		},
	})
	panelPage := panel_views.PanelPage("Magic Forge", forgePage)
	sidebar := panel_views.PanelSidebar()

	return views.RenderHTMX(c, panel_views.Base(sidebar, panelPage), panelPage)
}
