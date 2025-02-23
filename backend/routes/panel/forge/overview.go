package forge_routes

import (
	"context"

	github_utils "github.com/Liphium/magic/backend/util/github"
	"github.com/Liphium/magic/backend/views"
	panel_views "github.com/Liphium/magic/backend/views/panel"
	forge_views "github.com/Liphium/magic/backend/views/panel/forge"
	"github.com/gofiber/fiber/v2"
	"github.com/google/go-github/v69/github"
)

// Route: /a/panel/forge/:id
func targetsPage(c *fiber.Ctx) error {

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
			Commit: branch.GetCommit().GetSHA(),
		}
	}

	// Create the page and the sidebar based on the Forge
	page := panel_views.PanelPage(forge.Label, forge_views.TargetsPage(forge, rendered))

	return views.RenderHTMX(c, panel_views.Base(sidebar, page), page)
}
