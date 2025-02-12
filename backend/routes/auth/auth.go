package auth_routes

import "github.com/gofiber/fiber/v2"

func Unauthorized(router fiber.Router) {
	router.Get("/login", loginRoute)
}
