package admin_panel_routes

import (
	"github.com/Liphium/magic/backend/database"
	"github.com/Liphium/magic/backend/util"
	"github.com/Liphium/magic/backend/util/constants"
	"github.com/Liphium/magic/backend/views"
	panel_views "github.com/Liphium/magic/backend/views/panel"
	admin_panel_views "github.com/Liphium/magic/backend/views/panel/admin"
	"github.com/gofiber/fiber/v2"
)

// Route: /a/panel/admin/wizards
func wizardListPage(c *fiber.Ctx) error {

	// Get all wizards from the database
	var wizards []database.Wizard
	if err := database.DBConn.Order("updated_at DESC").Find(&wizards).Error; err != nil {
		return panel_views.RenderPanelError(c, "Something went wrong with the database. Please try again later.", err)
	}

	// Get a new wizard creation token
	token, err := util.WizardCreationToken()
	if err != nil {
		return panel_views.RenderPanelError(c, "Something went wrong during token creation. Please try again later.", err)
	}

	// Render all the wizards
	wizardPage := admin_panel_views.WizardListPage(token, wizards)
	panelPage := panel_views.PanelPage("Wizards", wizardPage)
	sidebar := panel_views.PanelSidebar(c.Locals(constants.LocalsPermissionLevel).(uint))

	return views.RenderHTMX(c, panel_views.Base(sidebar, panelPage), panelPage)
}
