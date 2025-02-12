package panel_routes

import (
	"github.com/Liphium/magic/backend/database"
	project_routes "github.com/Liphium/magic/backend/routes/panel/projects"
	"github.com/Liphium/magic/backend/views"
	"github.com/Liphium/magic/backend/views/components"
	panel_views "github.com/Liphium/magic/backend/views/panel"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

func Authorized(router fiber.Router) {
	router.Get("/", baseRoute)

	// Route to all the other endpoints in the panel
	router.Route("/projects", project_routes.Authorized)
}

// Route: /a/panel
func baseRoute(c *fiber.Ctx) error {
	return views.Render(c, panel_views.Base(panelBaseSidebar("/a/panel"), "Your Projects", panel_views.ProjectsPage([]database.Project{})))
}

func panelBaseSidebar(selected string) templ.Component {
	return components.Sidebar([]components.SBCategory{
		{
			Name: "Account",
			Links: []components.SBLink{
				{
					Name:     "Projects",
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
