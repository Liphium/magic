package auth_routes

import (
	"os"
	"time"

	"github.com/Liphium/magic/backend/database"
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

	// Store the session in the cookies
	c.Cookie(&fiber.Cookie{
		Name:     githubSessionCookie,
		Value:    session.Marshal(),
		Expires:  time.Now().Add(time.Hour * 24),
		Secure:   os.Getenv("MAGIC_LOCAL") != "true",
		SameSite: "lax",
	})

	// Redirect the user to the desired url
	return c.Redirect(url, fiber.StatusTemporaryRedirect)
}

// Route: /auth/gh/callback
func githubAuthCallback(c *fiber.Ctx) error {

	// Get the query and state parameter from GitHub
	code := c.Query("code", "-")
	state := c.Query("state", "-")

	// Make sure there wasn't an error
	if code == "-" || state == "-" {
		return views.RenderWithBase(c, components.ErrorPage("Seems like authorization with GitHub failed (1). Please try again."))
	}

	// Get the cookie and compare the state token
	cookie := c.Cookies(githubSessionCookie, "-")
	if cookie == "-" {
		return c.Redirect("/auth/gh/go", fiber.StatusTemporaryRedirect)
	}

	// Validate and authorize the session
	session, err := githubProvider.UnmarshalSession(cookie)
	if err != nil {
		return views.RenderWithBase(c, components.ErrorPage("Seems like authorization with GitHub failed (2). Please try again."))
	}
	if _, err = session.Authorize(githubProvider, &Params{
		ctx: c,
	}); err != nil {
		return views.RenderWithBase(c, components.ErrorPage("Seems like authorization with GitHub failed (3). Please try again."))
	}

	// Fetch the user to retrieve everything we need
	user, err := githubProvider.FetchUser(session)
	if err != nil {
		return views.RenderWithBase(c, components.ErrorPage("Seems like authorization with GitHub failed (4). Please try again."))
	}

	// Add to database or retrieve the account
	var account database.Account
	var credential database.Credential
	if err := database.DBConn.Where(&database.Credential{
		Type:   database.CredentialTypeGitHub,
		Secret: user.UserID,
	}).Take(&credential).Error; err != nil {

		// Make sure there is no account with this email already
		if err := database.DBConn.Where(&database.Account{
			Email: user.Email,
		}).Take(&database.Account{}).Error; err == nil {
			return views.RenderWithBase(c, components.ErrorPage("There already is an account with the E-Mail we got from GitHub. If you already have a Magic account, you can add GitHub to your old account as a connection instead."))
		}

		// Create an account for the user when the credential doesn't exist
		account = database.Account{
			Username: user.NickName,
			Email:    user.Email,
			Rank:     1,
		}
		if err := database.DBConn.Create(&account).Error; err != nil {
			return views.RenderWithBase(c, components.ErrorPage("Something went wrong on our side (1). Please try again."))
		}

		// Create a credential to link the account to the github account
		credential = database.Credential{
			Type:    database.CredentialTypeGitHub,
			Secret:  user.UserID,
			Account: account.ID,
		}
		if err := database.DBConn.Create(&credential).Error; err != nil {
			return views.RenderWithBase(c, components.ErrorPage("Something went wrong on our side (2). Please try again."))
		}
	} else {

		// Update and get the account in case it does exist
		if err := database.DBConn.Where(&database.Account{ID: credential.Account}).Take(&account); err != nil {
			return views.RenderWithBase(c, components.ErrorPage("Something went wrong on our side (3). Please try again."))
		}

		// Set all the new information
		account.Username = user.NickName
		account.Email = user.Email
		if err := database.DBConn.Save(&account).Error; err != nil {
			return views.RenderWithBase(c, components.ErrorPage("Something went wrong on our side (4). Please try again."))
		}
	}

	// TODO: New JWT token

	// For now let's just tell the user that it worked
	return views.RenderWithBase(c, components.ErrorPage("Auth actually worked, nice!"))
}
