package panel_routes

import (
	"fmt"
	"slices"

	"github.com/Liphium/magic/backend/database"
	forge_routes "github.com/Liphium/magic/backend/routes/panel/forge"
	preview_routes "github.com/Liphium/magic/backend/routes/panel/preview"
	"github.com/Liphium/magic/backend/util/constants"
	"github.com/Liphium/magic/backend/views"
	"github.com/Liphium/magic/backend/views/components"
	panel_views "github.com/Liphium/magic/backend/views/panel"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

func Authorized(router fiber.Router) {
	router.Get("/", baseRoute)
	router.Get("/installations", installationListPage)

	router.Route("/forge", forge_routes.Authorized)
	router.Route("/preview", preview_routes.Authorized)
}

// Route: /a/panel
func baseRoute(c *fiber.Ctx) error {

	// Get recently viewed things from the database
	var forges []database.Forge
	var previews []database.Preview
	if err := database.DBConn.Order("last_viewed DESC").Limit(10).Find(&forges).Error; err != nil {
		return views.RenderWithBase(c, components.ErrorPage(""))
	}
	if err := database.DBConn.Order("last_viewed DESC").Limit(10).Find(&previews).Error; err != nil {
		return views.RenderWithBase(c, components.ErrorPage(""))
	}

	// Merge all of the thingies together
	var recentViews []panel_views.RecentlyViewed
	for _, f := range forges {
		recentViews = append(recentViews, panel_views.RecentlyViewed{
			Label:       f.Label,
			Description: f.Repository,
			URL:         templ.SafeURL("/a/panel/forge/" + f.ID.String()),
			Time:        f.LastViewed,
		})
	}
	for _, p := range previews {
		recentViews = append(recentViews, panel_views.RecentlyViewed{
			Label:       p.Label,
			Description: p.Repository,
			URL:         templ.SafeURL("/a/panel/preview/" + p.Forge.String()),
			Time:        p.LastViewed,
		})
	}

	// Sort the recently viewed
	slices.SortFunc(recentViews, func(rv1, rv2 panel_views.RecentlyViewed) int {
		return rv1.Time.Compare(rv2.Time)
	})

	// Generate the html required for this
	welcome := panel_views.WelcomePage(recentViews)
	panelPage := panel_views.PanelPage(fmt.Sprintf("Welcome, %s!", c.Locals(constants.LocalsAccountName)), welcome)
	sidebar := panel_views.PanelSidebar()

	return views.RenderHTMX(c, panel_views.Base(sidebar, panelPage), panelPage)
}
