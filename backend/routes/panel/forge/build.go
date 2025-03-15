package forge_routes

import (
	"bufio"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Liphium/magic/backend/database"
	"github.com/Liphium/magic/backend/util"
	"github.com/Liphium/magic/backend/views"
	form_views "github.com/Liphium/magic/backend/views/forms"
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

	// Get the base info for later
	forge, _, valid := getBaseInfo(c)
	if !valid {
		c.Set("HX-Redirect", "/a/panel/forge/")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// Get the build
	buildId := c.Params("build", "")
	if buildId == "" {
		c.Set("HX-Redirect", "/a/panel/forge/")
		return c.SendStatus(fiber.StatusBadRequest)
	}
	var build database.Build
	if err := database.DBConn.Where("forge = ? AND id = ?", forge.ID.String(), buildId).Take(&build).Error; err != nil {
		c.Set("HX-Redirect", "/a/panel/forge/")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// Start sending random events
	return util.StartEvents(c, func(w *util.EventWriter) {
		logsStarted := false
		for {
			// Get the current status of the build
			if err := database.DBConn.Where("forge = ? AND id = ?", forge.ID.String(), buildId).Take(&build).Error; err != nil {
				str, err := views.RenderToString(form_views.FormSubmitError("Something went wrong during the build process."))
				if err != nil {
					log.Println("error during rendering:", err)
					break
				}

				w.SendEvent("end", str)
				break
			}

			endSending := false
			switch build.Status {
			case database.BuildStatusStarting:
				w.SendEvent("progress", "Starting server..")
			case database.BuildStatusStarted:
				w.SendEvent("progress", "Building image..")

				// Start the log stream
				if !logsStarted {
					logsStarted = true
					go func() {
						client := &http.Client{}
						req, err := http.NewRequest("POST", "http://localhost:9000/logs", nil)
						if err != nil {
							log.Println("Error creating log stream request:", err)
							return
						}

						resp, err := client.Do(req)
						if err != nil {
							log.Println("Error connecting to log stream:", err)
							return
						}
						defer resp.Body.Close()

						scanner := bufio.NewScanner(resp.Body)
						for scanner.Scan() {
							line := scanner.Text()
							if strings.HasPrefix(line, "data: ") {
								logData := strings.TrimPrefix(line, "data: ")
								w.SendEvent("log", "<p class=\"text-middle-text\">"+logData+"</p>")
							}
						}

						if err := scanner.Err(); err != nil {
							log.Println("Error reading log stream:", err)
						}
					}()
				}

			case database.BuildStatusError:

				// Send an error in case the build failed
				endSending = true
				str, err := views.RenderToString(form_views.FormSubmitError("Something went wrong during the build process."))
				if err != nil {
					log.Println("error during rendering:", err)
					break
				}

				w.SendEvent("end", str)
			case database.BuildStatusFinished:

				// Send a success message in case the build succeeded
				endSending = true
				str, err := views.RenderToString(form_views.FormSubmitSuccess("Build finished successfully."))
				if err != nil {
					log.Println("error during rendering:", err)
					break
				}

				w.SendEvent("end", str)
			}

			if endSending {
				break
			}

			time.Sleep(700 * time.Millisecond)
		}
	})
}
