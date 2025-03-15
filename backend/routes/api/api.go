package api_routes

import (
	spellcast_routes "github.com/Liphium/magic/backend/routes/api/spellcast"
	wizard_routes "github.com/Liphium/magic/backend/routes/api/wizard"
	"github.com/gofiber/fiber/v2"
)

func Unauthorized(router fiber.Router) {
	router.Route("/wizard", wizard_routes.Unauthorized)
	router.Route("/spellcast", spellcast_routes.Unauthorized)
}
