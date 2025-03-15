package spellcast_routes

import (
	spellcast_forge_routes "github.com/Liphium/magic/backend/routes/api/spellcast/forge"
	"github.com/gofiber/fiber/v2"
)

func Unauthorized(router fiber.Router) {
	router.Route("/forge", spellcast_forge_routes.Unauthorized)
}
