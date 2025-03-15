package spellcast_forge_routes

import (
	"fmt"

	"github.com/Liphium/magic/backend/database"
	"github.com/Liphium/magic/backend/util/constants"
	github_utils "github.com/Liphium/magic/backend/util/github"
	"github.com/gofiber/fiber/v2"
)

// Route: /api/spellcast/forge/connect
func acceptConnection(c *fiber.Ctx) error {
	build := c.Locals(constants.LocalsForgeBuild).(database.Build)

	// Get the related forge
	var forge database.Forge
	if err := database.DBConn.Where("id = ?", build.Forge).Take(&forge).Error; err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Get the GitHub API client and access token
	_, repo, token, err := github_utils.ClientAndRepoFromForge(forge)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Set the status of the build to started
	build.Status = database.BuildStatusStarted
	if err := database.DBConn.Save(&build).Error; err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Build a clone command
	command := fmt.Sprintf("git clone https://x-access-token:%s@github.com/%s/%s.git .", token, repo.Owner, repo.Name)

	return c.JSON(fiber.Map{
		"success":       true,
		"owner":         repo.Owner,
		"repository":    repo.Name,
		"clone_command": command,
	})
}
