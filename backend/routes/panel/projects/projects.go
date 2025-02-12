package project_routes

import (
	"github.com/Liphium/magic/backend/views"
	"github.com/Liphium/magic/backend/views/components"
	panel_views "github.com/Liphium/magic/backend/views/panel"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

func Authorized(router fiber.Router) {
	router.Get("/:id", projectBase)
}

func projectBase(c *fiber.Ctx) error {
	return views.Render(c, panel_views.Base(projectSidebar(), c.Params("id"), components.LinkButtonPrimary("hi", "hi")))
}

func projectSidebar() templ.Component {
	return components.Sidebar([]components.SBCategory{
		{
			Name: "Project",
			Links: []components.SBLink{
				{
					Name: "Overview",
					Link: "/a/panel/projects/...",
				},
			},
		},
		{
			Name: "Forge",
			Links: []components.SBLink{
				{
					Name: "Overview",
					Link: "/a/panel/projects/.../forges",
				},
				{
					Name: "Builds",
					Link: "/a/panel/projects/.../builds",
				},
				{
					Name: "Assets",
					Link: "/a/panel/projects/.../assets",
				},
			},
		},
		{
			Name: "Preview",
			Links: []components.SBLink{
				{
					Name: "Environments",
					Link: "/a/panel/projects/.../environments",
				},
				{
					Name: "Configurations",
					Link: "/a/panel/projects/.../settings",
				},
			},
		},
	})
}
