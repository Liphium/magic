package wizard_routes

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func wizardInit(c *fiber.Ctx) error {
	var req struct {
		Token string `json:"token"`
		Host  string `json:"host"`
	}

	if err := c.BodyParser(&req); err != nil {
		return err
	}

	_, err := jwt.Parse(req.Token, func(t *jwt.Token) (interface{}, error) {

		// Make sure it's the correct method
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method for jwt: %v", t.Header["alg"])
		}

		return []byte(os.Getenv("MAGIC_JWT_SECRET")), nil
	})
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}
