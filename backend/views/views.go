package views

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

// Render a templ component through fiber
func Render(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	return Base(component).Render(c.Context(), c.Response().BodyWriter())
}
