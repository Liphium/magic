package panel_routes

import (
	"context"

	github_utils "github.com/Liphium/magic/backend/util/github"
	"github.com/Liphium/magic/backend/views"
	panel_views "github.com/Liphium/magic/backend/views/panel"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/google/go-github/v69/github"
)

// Route: /a/panel/installations
func installationListPage(c *fiber.Ctx) error {

	// Get the user access token for Github
	client, err := github_utils.GetUserFromContext(c)
	if err != nil {
		return panel_views.RenderPanelError(c, "We couldn't get your GitHub credentials. Please try again later.", err)
	}

	// Get all the installations they have made
	installation, res, err := client.Apps.ListUserInstallations(context.Background(), &github.ListOptions{
		PerPage: 100,
	})
	if err != nil {
		return panel_views.RenderPanelError(c, "Something went wrong with GitHub. Maybe check their status?", err)
	}
	if res.Rate.Remaining <= 0 {
		return panel_views.RenderPanelError(c, "Your rate limit for GitHub has been hit. Please try again in 1 hour.", nil)
	}

	// Render all the installations
	rendered := []panel_views.RenderedInstallation{}
	for _, i := range installation {
		rendered = append(rendered, panel_views.RenderedInstallation{
			Name:     i.GetAccount().GetLogin(),
			Provider: "GitHub",
			URL:      templ.SafeURL(i.GetHTMLURL()),
		})
	}

	// Render all the content
	installationPage := panel_views.InstallationPage(rendered)
	sidebar := panel_views.PanelSidebar()
	page := panel_views.PanelPage("Installations", installationPage)

	return views.RenderHTMX(c, panel_views.Base(sidebar, page), page)
}
