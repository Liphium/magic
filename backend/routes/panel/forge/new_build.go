package forge_routes

import (
	"context"

	"github.com/Liphium/magic/backend/database"
	github_utils "github.com/Liphium/magic/backend/util/github"
	"github.com/Liphium/magic/backend/views"
	panel_views "github.com/Liphium/magic/backend/views/panel"
	forge_views "github.com/Liphium/magic/backend/views/panel/forge"
	"github.com/gofiber/fiber/v2"
	"github.com/google/go-github/v69/github"
)

// Route: /a/panel/forge/:id/builds/new
func newBuildPage(c *fiber.Ctx) error {

	// Get all of the base stuff for later
	forge, sidebar, valid := getBaseInfo(c)
	if !valid {
		return c.Redirect("/a/panel/forge", fiber.StatusPermanentRedirect)
	}

	// Get the repository and github client for the Forge
	client, repo, err := github_utils.ClientAndRepoFromForge(forge)
	if err != nil {
		return c.Redirect("/a/panel/forge", fiber.StatusTemporaryRedirect)
	}

	// Get all the branches for the UI
	branches, res, err := client.Repositories.ListBranches(context.Background(), repo.Owner, repo.Name, &github.BranchListOptions{})
	if err != nil {
		return c.Redirect("/a/panel/forge", fiber.StatusTemporaryRedirect)
	}
	if res.Rate.Remaining <= 0 {
		return c.Redirect("/a/panel/forge", fiber.StatusTemporaryRedirect)
	}

	// Render all the branches
	rendered := make([]forge_views.Branch, len(branches))
	for i, branch := range branches {
		rendered[i] = forge_views.Branch{
			Name:   branch.GetName(),
			Commit: branch.Commit.GetSHA(),
		}
	}

	// Render all the branches
	forgePage := forge_views.NewBuildPage(forge, rendered)
	panelPage := panel_views.PanelPage("Build a branch", forgePage)

	return views.RenderHTMX(c, panel_views.Base(sidebar, panelPage), panelPage)
}

// Route: /a/panel/forge/:id/builds/new/:branch
func newBuildRequest(c *fiber.Ctx) error {

	// Get all of the base stuff for later
	forge, _, valid := getBaseInfo(c)
	if !valid {
		return views.RenderError(c, forge_views.BranchPageMessage(true, "This forge doesn't exist."), nil)
	}

	// Check if there is a branch path param
	branchId := c.Params("branch", "")
	if branchId == "" {
		return views.RenderError(c, forge_views.BranchPageMessage(true, "Please specify a branch."), nil)
	}

	// Get the repository and github client for the Forge
	client, repo, err := github_utils.ClientAndRepoFromForge(forge)
	if err != nil {
		return views.RenderError(c, forge_views.BranchPageMessage(true, "Couldn't interact with GitHub."), err)
	}

	// Get the branch from the repository
	branch, res, err := client.Repositories.GetBranch(context.Background(), repo.Owner, repo.Name, branchId, 5)
	if err != nil {
		return views.RenderError(c, forge_views.BranchPageMessage(true, "Couldn't find branch on GitHub."), nil)
	}
	if res.Rate.Remaining <= 0 {
		return views.RenderError(c, forge_views.BranchPageMessage(true, "You hit your GitHub rate limit."), nil)
	}

	// Make space for a new build
	// TODO: Make sure this deletes everything (like artifacts and stuff too)
	if err := database.DBConn.Raw("delete from builds where id in ( select id from builds order by created_at desc limit 1000 offset 15 )").Error; err != nil {
		return views.RenderError(c, forge_views.BranchPageMessage(true, "Something went wrong on the server."), err)
	}

	// Create a new build
	build := database.Build{
		Forge:       forge.ID,
		DisplayName: branch.Commit.Commit.GetMessage(),
		Source:      database.BuildSourceBranch(branch.GetName()),
	}
	if err := database.DBConn.Create(&build).Error; err != nil {
		return views.RenderError(c, forge_views.BranchPageMessage(true, "Something went wrong during creation."), err)
	}

	// Redirect to the page of the build
	c.Set("HX-Redirect", "/a/panel/forge/"+forge.ID.String()+"/builds/"+build.ID.String())

	return views.RenderJust(c, forge_views.BranchPageMessage(false, "Redirecting to the build.."))
}
