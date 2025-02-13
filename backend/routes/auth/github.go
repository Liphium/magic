package auth_routes

import (
	"github.com/Liphium/magic/backend/util"
	"github.com/Liphium/magic/backend/views"
	"github.com/Liphium/magic/backend/views/components"
	"github.com/gofiber/fiber/v2"
)

// Route: /auth/gh/go
func goToGitHubAuth(c *fiber.Ctx) error {

	// Get the state from the query parameter or just start a new auth session
	state := c.Query("state", util.GenerateToken(64))

	// Generate a new session for the auth process
	session, err := githubProvider.BeginAuth(state)
	if err != nil {
		return views.RenderWithBase(c, components.ErrorPage("Something went wrong with GitHub (1). Please try again later."))
	}

	// Get the actual url for the session
	url, err := session.GetAuthURL()
	if err != nil {
		return views.RenderWithBase(c, components.ErrorPage("Something went wrong with GitHub (2). Please try again later."))
	}

	// Redirect the user to the desired url
	return c.Redirect(url, fiber.StatusTemporaryRedirect)
}

// Route: /auth/gh/callback
func githubAuthCallback(c *fiber.Ctx) error {

	// Get the query parameter from GitHub
	code := c.Query("code", "-")

	// Make sure there wasn't an error
	if code == "-" {
		return views.RenderWithBase(c, components.ErrorPage("Seems like authorization with GitHub failed. Please try again."))
	}

	// TODO: Actually handle the stuff

	// For now let's just tell the user that it worked
	return views.RenderWithBase(c, components.ErrorPage("Auth actually worked, nice!"))
}
