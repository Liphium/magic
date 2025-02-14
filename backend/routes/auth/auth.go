package auth_routes

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
)

// I'd like to say this here because it's the easiest place to say it.
// Thanks to Shareed2k this whole thing was way easier than implementing all of it by
// myself without any reference: https://github.com/Shareed2k/goth_fiber.
//
// Thanks for your amazing work and because I still wanted to mention you, this is here.
// If you are not credited on the credits page on the Liphium Magic website when it's finally
// done, please let me know and we shall change that.

var githubProvider *github.Provider

func Unauthorized(router fiber.Router) {

	// Get the callback url for github
	callback := os.Getenv("MAGIC_PROTOCOL") + os.Getenv("MAGIC_DOMAIN") + "/auth/gh/callback"
	if os.Getenv("MAGIC_LOCAL") == "true" {
		callback = os.Getenv("MAGIC_PROTOCOL") + os.Getenv("MAGIC_DOMAIN") + ":" + os.Getenv("TEMPL_PROXY_PORT") + "/auth/gh/callback"
	}

	// Create a new github auth provider
	githubProvider = github.New(os.Getenv("MAGIC_GH_CLIENT"), os.Getenv("MAGIC_GH_SECRET"), callback, "user:email")
	goth.UseProviders(githubProvider)

	router.Get("/login", loginRoute)

	router.Get("/gh/go", goToGitHubAuth)
	router.Get("/gh/callback", githubAuthCallback)
}

// Implementation of goth.Params copied from goth_fiber (look at notice above)
type Params struct {
	ctx *fiber.Ctx
}

func (p *Params) Get(key string) string {
	return p.ctx.Query(key)
}
