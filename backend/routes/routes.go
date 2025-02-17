package routes

import (
	"fmt"
	"os"

	auth_routes "github.com/Liphium/magic/backend/routes/auth"
	panel_routes "github.com/Liphium/magic/backend/routes/panel"
	wizard_routes "github.com/Liphium/magic/backend/routes/wizard"
	"github.com/Liphium/magic/backend/util/constants"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func InitializeRoutes(router fiber.Router) {

	// Serve the static files for the frontend
	router.Static("/static/", "./static")

	// Add the skip auth endpoint (only when enabled)
	if os.Getenv("MAGIC_SKIP_AUTH") == "true" {
		router.Get("_dev/skip", skipAuth)
	}

	// Add all the routes
	router.Route("/", unauthorizedRouter)
	router.Route("/a", authorizedRouter)
}

func unauthorizedRouter(router fiber.Router) {
	router.Route("/auth", auth_routes.Unauthorized)
	router.Route("/wizard", wizard_routes.Unauthorized)
}

func authorizedRouter(router fiber.Router) {

	// Add an auth middleware that parses thw JWT tokens
	router.Use(func(c *fiber.Ctx) error {

		// Get the cookie
		tokenString := c.Cookies(constants.CookieMagicSession, "-")
		if tokenString == "-" {
			return c.Redirect("/auth/login", fiber.StatusTemporaryRedirect)
		}

		// Validate the actual token from the cookie
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {

			// Make sure it's the correct method
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signing method for jwt: %v", t.Header["alg"])
			}

			return []byte(os.Getenv("MAGIC_JWT_SECRET")), nil
		})
		if err != nil {
			return c.Redirect("/auth/login", fiber.StatusTemporaryRedirect)
		}

		// Make sure the claims are valid and actually there
		var claims jwt.MapClaims
		var ok bool
		if claims, ok = token.Claims.(jwt.MapClaims); !ok {
			return c.Redirect("/auth/login", fiber.StatusTemporaryRedirect)
		}

		// Parse the values into the fiber context
		if acc, ok := claims["acc"]; ok {
			c.Locals(constants.LocalsAccountID, acc)
		} else {
			return c.Redirect("/auth/login", fiber.StatusTemporaryRedirect)
		}
		if name, ok := claims["name"]; ok {
			c.Locals(constants.LocalsAccountName, name)
		} else {
			return c.Redirect("/auth/login", fiber.StatusTemporaryRedirect)
		}
		if pLvl, ok := claims["plvl"]; ok {
			c.Locals(constants.LocalsPermissionLevel, uint(pLvl.(float64)))
		} else {
			return c.Redirect("/auth/login", fiber.StatusTemporaryRedirect)
		}

		// Continue to the handler
		return c.Next()
	})

	router.Route("/panel", panel_routes.Authorized)
}
