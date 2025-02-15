package forge_routes

import (
	"github.com/Liphium/magic/backend/database"
	"github.com/Liphium/magic/backend/util/constants"
	"github.com/Liphium/magic/backend/views"
	panel_views "github.com/Liphium/magic/backend/views/panel"
	forge_views "github.com/Liphium/magic/backend/views/panel/forge"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Authorized(router fiber.Router) {
	router.Get("/", baseRoute)

	// Process for creating a new Forge
	router.Get("/new", newForgeStart)

	// All views when a forge id is present
	router.Get("/:id", targetsPage)
}

// Route: /a/panel/forge
func baseRoute(c *fiber.Ctx) error {

	// Get all forges from the database
	var forges []database.Forge
	if err := database.DBConn.Where("account = ?", c.Locals(constants.LocalsAccountID)).Order("last_viewed DESC").Find(&forges).Error; err != nil {
		return panel_views.RenderPanelError(c, "Something went wrong with the database. Please try again later.", err)
	}

	// Render all the forges
	forgePage := forge_views.ForgeListPage(forges)
	panelPage := panel_views.PanelPage("Magic Forge", forgePage)
	sidebar := panel_views.PanelSidebar()

	return views.RenderHTMX(c, panel_views.Base(sidebar, panelPage), panelPage)
}

// Get all of the base information from the route
func getBaseInfo(c *fiber.Ctx) (database.Forge, templ.Component, bool) {

	// Try parsing the forge id retrieved from the request
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return database.Forge{}, nil, false
	}

	// TODO: Get the Forge from the database
	forge := database.Forge{
		ID:         id,
		Label:      "Liphium station",
		Repository: "https://github.com/Liphium/station",
	}

	// Create the sidebar just to save some repeated code
	sidebar := forge_views.ForgeSidebar(forge)

	return forge, sidebar, true
}
