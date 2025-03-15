package form_routes

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Liphium/magic/backend/database"
	"github.com/Liphium/magic/backend/util/constants"
	github_utils "github.com/Liphium/magic/backend/util/github"
	"github.com/Liphium/magic/backend/views"
	form_views "github.com/Liphium/magic/backend/views/forms"
	"github.com/gofiber/fiber/v2"
	"github.com/google/go-github/v69/github"
)

// Route: /a/panel/_forms/repository/installations
func repoPickerInstallations(c *fiber.Ctx) error {

	// Get the header for the input id
	inputId := c.Get("M-ID", "-")
	if inputId == "-" {
		return form_views.RenderFormError(c, "Something went wrong again. Please try again later.", nil)
	}

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
			ID:   i.GetID(),
			Name: i.Account.GetLogin(),
			URL:  fmt.Sprintf("/a/panel/_forms/repository/gh/%d", i.GetID()),
		})
	}

	// Render the chips
	return views.RenderJust(c, form_views.InstallationChips(inputId, rendered))
}

// Route: /a/panel/forms/repository/gh/:id
func getGitHubRepositories(c *fiber.Ctx) error {

	// Get the app id
	installationIdStr := c.Params("id", "-")
	if installationIdStr == "-" {
		return form_views.RenderFormError(c, "We couldn't find this repository (1). Please try again later.", nil)
	}
	installationId, err := strconv.ParseInt(installationIdStr, 10, 64)
	if err != nil {
		return form_views.RenderFormError(c, "We couldn't find this repository (2). Please try again later.", nil)
	}

	// Get the user access token for Github
	acc := c.Locals(constants.LocalsAccountID)
	var credential database.Credential
	if err := database.DBConn.Where("account = ?", acc).Take(&credential).Error; err != nil {
		return form_views.RenderFormError(c, "We couldn't get your GitHub credentials. Please try again later.", err)
	}

	// Get all installations from Github
	client, _, err := github_utils.GetInstallationClient(installationId)
	if err != nil {
		return form_views.RenderFormError(c, "Something went wrong with GitHub (1). Maybe check their status?", err)
	}
	repos, res, err := client.Apps.ListRepos(context.Background(), &github.ListOptions{})
	if err != nil {
		return form_views.RenderFormError(c, "Something went wrong with GitHub (2). Maybe check their status?", err)
	}
	if res.Rate.Remaining <= 0 {
		return form_views.RenderFormError(c, "Your rate limit for GitHub has been hit. Please try again in 1 hour.", nil)
	}

	// Parse to output
	var rendered []form_views.Repository
	for _, r := range repos.Repositories {
		rendered = append(rendered, form_views.Repository{
			Name: r.GetName(),
			ID:   r.GetFullName(),
			URL:  r.GetURL(),
		})
	}

	// Render the repositories
	return views.RenderJust(c, form_views.RenderRepositories(rendered))
}
