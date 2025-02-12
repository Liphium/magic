package panel_routes

import (
	forge_routes "github.com/Liphium/magic/backend/routes/panel/forge"
	preview_routes "github.com/Liphium/magic/backend/routes/panel/preview"
	"github.com/Liphium/magic/backend/views"
	panel_views "github.com/Liphium/magic/backend/views/panel"
	"github.com/gofiber/fiber/v2"
)

func Authorized(router fiber.Router) {
	router.Get("/", baseRoute)

	router.Route("/forge", forge_routes.Authorized)
	router.Route("/preview", preview_routes.Authorized)
}

// Route: /a/panel
func baseRoute(c *fiber.Ctx) error {

	welcome := panel_views.WelcomePage([]panel_views.RecentlyViewed{
		{
			Label:       "Liphium",
			Description: "Last viewed on 04/21/2025",
			URL:         "/a/panel/forge/...",
		},
	})
	panelPage := panel_views.PanelPage("Welcome, Unbreathable!", welcome)
	sidebar := panel_views.PanelSidebar()

	return views.RenderHTMX(c, panel_views.Base(sidebar, panelPage), panelPage)
}
