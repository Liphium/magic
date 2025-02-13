package auth_routes

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth/providers/github"
)

var githubProvider *github.Provider

func Unauthorized(router fiber.Router) {

	// Create a new github auth provider
	githubProvider = github.New(os.Getenv("MAGIC_GH_CLIENT"), os.Getenv("MAGIC_GH_SECRET"), os.Getenv("MAGIC_GH_CALLBACK"), "email")

	router.Get("/login", loginRoute)

	router.Get("/gh/go")
	router.Get("/gh/callback")
}
