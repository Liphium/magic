package panel_routes

import (
	"github.com/Liphium/magic/backend/views"
	"github.com/Liphium/magic/backend/views/components"
	panel_views "github.com/Liphium/magic/backend/views/panel"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

func Authorized(router fiber.Router) {
	router.Get("/", baseRoute)
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

	return views.Render(c, panel_views.Base(panelBaseSidebar("/a/panel"), "Welcome, Unbreathable!", welcome))
}

func panelBaseSidebar(selected string) templ.Component {
	return components.Sidebar([]components.SBCategory{
		{
			Name: "Account",
			Links: []components.SBLink{
				{
					Name:     "Welcome",
					Link:     "/a/panel",
					Selected: selected == "/a/panel",
				},
			},
		},
		{
			Name: "Magic",
			Links: []components.SBLink{
				{
					Name:     "Contact & support",
					Link:     "/a/panel/contact",
					Selected: selected == "/a/panel/contact",
				},
				{
					Name: "Terms of service",
					Link: "https://liphium.com/legal/terms",
				},
				{
					Name: "Privacy policy",
					Link: "https://liphium.com/legal/terms",
				},
			},
		},
	})
}
