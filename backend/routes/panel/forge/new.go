package forge_routes

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/Liphium/magic/backend/database"
	"github.com/Liphium/magic/backend/util"
	github_utils "github.com/Liphium/magic/backend/util/github"
	"github.com/Liphium/magic/backend/views"
	form_views "github.com/Liphium/magic/backend/views/forms"
	panel_views "github.com/Liphium/magic/backend/views/panel"
	forge_views "github.com/Liphium/magic/backend/views/panel/forge"
	"github.com/gofiber/fiber/v2"
)

// Route: /a/panel/forge/new (get for the page)
func newForgePage(c *fiber.Ctx) error {

	// Render all the forges
	stepPage := forge_views.NewForgeStep1()
	panelPage := panel_views.PanelPageBase(stepPage)
	sidebar := panel_views.PanelSidebar()

	return views.RenderHTMX(c, panel_views.Base(sidebar, panelPage), panelPage)
}

// Route: /a/panel/forge/new (post for the submission)
func createNewForge(c *fiber.Ctx) error {

	// Get and parse the repository id
	repo := c.FormValue("repository", "-")
	if repo == "-" || repo == "0" || repo == "" {
		return views.RenderError(c, form_views.FormSubmitError("Please select a repository for the Forge!"), nil)
	}
	args := strings.Split(repo, "/")
	if len(args) != 2 {
		return views.RenderError(c, form_views.FormSubmitError("Please select a valid GitHub repository for the Forge!"), nil)
	}
	owner := args[0]
	repoSlug := args[1]

	// Get the actual installation id
	installation := c.FormValue("repository_ins", "-")
	if repo == "-" || repo == "0" || repo == "" {
		return views.RenderError(c, form_views.FormSubmitError("Something went wrong on the server (1)."), nil)
	}
	installationId, err := strconv.ParseInt(installation, 10, 64)
	if err != nil {
		return views.RenderError(c, form_views.FormSubmitError("Something went wrong on the server (2)."), err)
	}

	// Get the GitHub client for the installation
	client, err := github_utils.GetInstallationClient(installationId)
	if err != nil {
		return views.RenderError(c, form_views.FormSubmitError("Something went wrong on the server (3)."), err)
	}

	// Check if we have access to this repository
	repository, res, err := client.Repositories.Get(context.Background(), owner, repoSlug)
	if err != nil {
		return views.RenderError(c, form_views.FormSubmitError("Something went wrong with the GitHub API!"), err)
	}
	if res.Rate.Remaining <= 0 {
		return views.RenderError(c, form_views.FormSubmitError("Seems like we ran out of requests to GitHub. Please try again later."), nil)
	}

	// If we add other providers we could get this from the repository probably
	provider := database.ProviderTypeGitHub

	// Check the name provided
	name := c.FormValue("name", "-")
	if name == "-" || name == "" || len(name) <= 3 {
		return views.RenderError(c, form_views.FormSubmitError("Please choose a name that is longer than 3 characters!"), err)
	}
	if len(name) >= 20 {
		return views.RenderError(c, form_views.FormSubmitError("Please choose a name that is shorter than 20 characters!"), err)
	}

	// Create the forge
	forge := database.Forge{
		Account:        util.AccountUUID(c),
		Provider:       provider,
		Installation:   installation,
		Repository:     repository.GetFullName(),
		RepositoryName: repository.GetFullName(),
		Label:          name,
		LastViewed:     time.Now(),
	}
	if err := database.DBConn.Create(&forge).Error; err != nil {
		return views.RenderError(c, form_views.FormSubmitError("Seems like there was a server error. Please try again later."), err)
	}

	// Redirect to Forge
	c.Set("HX-Redirect", "/a/panel/forge")

	return views.RenderJust(c, form_views.FormSubmitSuccess("Everything worked!"))
}
