package form_routes

import "github.com/gofiber/fiber/v2"

func Authorized(router fiber.Router) {

	// Repository picker component
	router.Get("/repository/installations", repoPickerInstallations)
	router.Get("/repository/gh/:id", getGitHubRepositories)
}
