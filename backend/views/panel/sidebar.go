package panel_views

import (
	"github.com/Liphium/magic/backend/views/components"
	"github.com/a-h/templ"
)

func PanelSidebar() templ.Component {
	return components.Sidebar([]components.SBCategory{
		{
			Name: "Account",
			Links: []components.SBLink{
				{
					Name: "Welcome",
					Link: "/a/panel",
				},
			},
		},
		{
			Name: "Forge & Preview",
			Links: []components.SBLink{
				{
					Name: "Forge",
					Link: "/a/panel/forge",
				},
				{
					Name: "Preview",
					Link: "/a/panel/preview",
				},
				{
					Name: "Environments",
					Link: "/a/panel/environments",
				},
			},
		},
		{
			Name: "Legal documents",
			Links: []components.SBLink{
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
