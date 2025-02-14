package panel_routes

import (
	"context"

	"github.com/Liphium/magic/backend/database"
	"github.com/Liphium/magic/backend/util/constants"
	"github.com/Liphium/magic/backend/views"
	panel_views "github.com/Liphium/magic/backend/views/panel"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/google/go-github/v69/github"
)

// Route: /a/panel/installations
func installationListPage(c *fiber.Ctx) error {

	// Get the user access token for Github
	acc := c.Locals(constants.LocalsAccountID)
	var credential database.Credential
	if err := database.DBConn.Where("account = ?", acc).Take(&credential).Error; err != nil {
		return panel_views.RenderPanelError(c, "We couldn't get your GitHub credentials. Please try again later.", err)
	}

	// Get all installations from Github
	client := github.NewClient(nil).WithAuthToken(credential.Token)
	installation, res, err := client.Apps.ListUserInstallations(context.Background(), &github.ListOptions{})
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
