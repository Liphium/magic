package routes

import (
	"log"
	"os"

	auth_routes "github.com/Liphium/magic/backend/routes/auth"
	panel_routes "github.com/Liphium/magic/backend/routes/panel"
	"github.com/Liphium/magic/backend/util/constants"
	jwtware "github.com/gofiber/contrib/jwt"
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
}

func authorizedRouter(router fiber.Router) {

	// Add an auth middleware that parses thw JWT tokens
	router.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.HS512,
			Key:    []byte(os.Getenv("MAGIC_JWT_SECRET")),
		},

		// A success handler for passing down the arguments in the jwt token using context
		SuccessHandler: func(c *fiber.Ctx) error {

			// Get user claims
			user, valid := c.Locals("user").(*jwt.Token)
			if !valid {
				return c.SendStatus(fiber.StatusInternalServerError)
			}
			claims := user.Claims.(jwt.MapClaims)

			// Parse the values into the fiber context
			c.Locals(constants.LocalsAccountID, claims["acc"])
			c.Locals(constants.LocalsPermissionLevel, claims["plvl"])
			return c.Next()
		},

		// Error handler (log errors in case they happen)
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Println(err.Error())
			// Return error message
			return c.SendStatus(401)
		},
	}))

	router.Route("/panel", panel_routes.Authorized)
}
