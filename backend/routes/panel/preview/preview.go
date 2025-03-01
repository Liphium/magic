package preview_routes

import (
	"github.com/Liphium/magic/backend/database"
	"github.com/Liphium/magic/backend/util/constants"
	"github.com/Liphium/magic/backend/views"
	panel_views "github.com/Liphium/magic/backend/views/panel"
	"github.com/gofiber/fiber/v2"
)

func Authorized(router fiber.Router) {
	router.Get("/", baseRoute)
}

// Route: /a/panel/preview
func baseRoute(c *fiber.Ctx) error {

	previewPage := panel_views.PreviewPage([]database.Preview{})
	panelPage := panel_views.PanelPage("Magic Preview", previewPage)
	sidebar := panel_views.PanelSidebar(c.Locals(constants.LocalsPermissionLevel).(uint))

	return views.RenderHTMX(c, panel_views.Base(sidebar, panelPage), panelPage)
}
