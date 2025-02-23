package forge_routes

import (
	"github.com/Liphium/magic/backend/database"
	"github.com/Liphium/magic/backend/views"
	forge_views "github.com/Liphium/magic/backend/views/panel/forge"
	"github.com/gofiber/fiber/v2"
)

func buildBranch(c *fiber.Ctx) error {

	forge, _, valid := getBaseInfo(c)
	if !valid {
		return views.RenderJust(c, forge_views.BranchPageMessage(true, "This Forge couldn't be found."))
	}

	if err := database.DBConn.Create(&database.Job{
		Type:   database.JobTypeBuild,
		Target: forge.Repository,
	}).Error; err != nil {
		return views.RenderError(c, forge_views.BranchPageMessage(true, "Couldn't create Job."), err)
	}

	return views.RenderJust(c, forge_views.BranchPageMessage(false, "Building of the branch has been queued."))
}
