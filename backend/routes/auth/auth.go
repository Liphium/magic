package auth_routes

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
)

var githubProvider *github.Provider

func Unauthorized(router fiber.Router) {

	// Get the callback url for github
	callback := os.Getenv("MAGIC_PROTOCOL") + os.Getenv("MAGIC_DOMAIN") + "/auth/gh/callback"
	if os.Getenv("MAGIC_LOCAL") == "true" {
		callback = os.Getenv("MAGIC_PROTOCOL") + os.Getenv("MAGIC_DOMAIN") + ":" + os.Getenv("TEMPL_PROXY_PORT") + "/auth/gh/callback"
	}

	// Create a new github auth provider
	githubProvider = github.New(os.Getenv("MAGIC_GH_CLIENT"), os.Getenv("MAGIC_GH_SECRET"), callback, "email")
	goth.UseProviders(githubProvider)

	router.Get("/login", loginRoute)

	router.Get("/gh/go", goToGitHubAuth)
	router.Get("/gh/callback", githubAuthCallback)
}
