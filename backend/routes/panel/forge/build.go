package forge_routes

import (
	"fmt"
	"time"

	"github.com/Liphium/magic/backend/database"
	"github.com/Liphium/magic/backend/util"
	"github.com/Liphium/magic/backend/views"
	"github.com/Liphium/magic/backend/views/components"
	panel_views "github.com/Liphium/magic/backend/views/panel"
	forge_views "github.com/Liphium/magic/backend/views/panel/forge"
	"github.com/gofiber/fiber/v2"
)

// Route: /a/panel/forge/:id/builds/:build
func showBuildPage(c *fiber.Ctx) error {

	// Get the base info for later
	forge, sidebar, valid := getBaseInfo(c)
	if !valid {
		return c.Redirect("/a/panel/forge", fiber.StatusPermanentRedirect)
	}

	// Get the build
	buildId := c.Params("build", "")
	if buildId == "" {
		return c.Redirect("/a/panel/forge/"+forge.ID.String()+"/builds", fiber.StatusPermanentRedirect)
	}
	var build database.Build
	if err := database.DBConn.Where("forge = ? AND id = ?", forge.ID.String(), buildId).Take(&build).Error; err != nil {
		return c.Redirect("/a/panel/forge/"+forge.ID.String()+"/builds", fiber.StatusPermanentRedirect)
	}

	// Render the build page
	page := panel_views.PanelPage("Build details", forge_views.BuildViewPage(forge, build))
	return views.RenderHTMX(c, panel_views.Base(sidebar, page), page)
}

// Route: /a/panel/forge/builds/:id/builds/:build/progress
func pullBuildProgress(c *fiber.Ctx) error {

	// Start sending random events
	return util.StartEvents(c, func(w *util.EventWriter) {
		n := 0
		for {
			if n == 100 {

				// Render an error
				rendered, err := views.RenderToString(components.ErrorText("Something went wrong on the server."))
				if err != nil {
					break
				}

				if err := w.SendEvent("end", rendered); err != nil {
					break
				}

				break
			}

			if err := w.SendEvent("progress", fmt.Sprintf("Starting server.. (%d/100)", n)); err != nil {
				break
			}
			n++

			time.Sleep(100 * time.Millisecond)
		}
	})
}
