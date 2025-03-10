package admin_panel_routes

import (
	"fmt"
	"time"

	"github.com/Liphium/magic/backend/util"
	"github.com/Liphium/magic/backend/util/constants"
	"github.com/Liphium/magic/backend/views"
	"github.com/Liphium/magic/backend/views/components"
	panel_views "github.com/Liphium/magic/backend/views/panel"
	admin_panel_views "github.com/Liphium/magic/backend/views/panel/admin"
	"github.com/gofiber/fiber/v2"
)

// Route: /a/panel/admin/demo
func demoRoute(c *fiber.Ctx) error {

	// Render all the wizards
	wizardPage := admin_panel_views.DemoPage()
	panelPage := panel_views.PanelPage("Demo & testing", wizardPage)
	sidebar := panel_views.PanelSidebar(c.Locals(constants.LocalsPermissionLevel).(uint))

	return views.RenderHTMX(c, panel_views.Base(sidebar, panelPage), panelPage)
}

// Route: /a/panel/admin/demo/progress
func demoProgress(c *fiber.Ctx) error {

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
			}

			if err := w.SendEvent("progress", fmt.Sprintf("Loading.. (%d/100)", n)); err != nil {
				break
			}
			n++

			time.Sleep(100 * time.Millisecond)
		}
	})
}
