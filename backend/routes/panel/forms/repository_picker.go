package form_routes

import (
	"context"

	"github.com/Liphium/magic/backend/database"
	"github.com/Liphium/magic/backend/util/constants"
	"github.com/Liphium/magic/backend/views"
	form_views "github.com/Liphium/magic/backend/views/forms"
	"github.com/gofiber/fiber/v2"
	"github.com/google/go-github/v69/github"
)

// Route: /a/panel/_forms/repository/installations
func repoPickerInstallations(c *fiber.Ctx) error {

	// Get the user access token for Github
	acc := c.Locals(constants.LocalsAccountID)
	var credential database.Credential
	if err := database.DBConn.Where("account = ?", acc).Take(&credential).Error; err != nil {
		return form_views.RenderFormError(c, "We couldn't get your GitHub credentials. Please try again later.", err)
	}

	// Get all installations from Github
	client := github.NewClient(nil).WithAuthToken(credential.Token)
	installations, res, err := client.Apps.ListUserInstallations(context.Background(), &github.ListOptions{})
	if err != nil {
		return form_views.RenderFormError(c, "Something went wrong with GitHub. Maybe check their status?", err)
	}
	if res.Rate.Remaining <= 0 {
		return form_views.RenderFormError(c, "Your rate limit for GitHub has been hit. Please try again in 1 hour.", nil)
	}

	// Convert to installations
	var rendered []form_views.Installation
	for _, i := range installations {
		rendered = append(rendered, form_views.Installation{
			Provider: "GitHub",
			Name:     i.Account.GetLogin(),
		})
	}

	// Render the chips
	return views.RenderJust(c, form_views.InstallationChips(rendered))
}
