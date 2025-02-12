package views

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

// Render the component surrounded by html head and footer and stuff
func RenderWithBase(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	return Base(component).Render(c.Context(), c.Response().BodyWriter())
}

// Render just the component without the base
func RenderJust(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	return component.Render(c.Context(), c.Response().BodyWriter())
}

// Render full in case of no HTMX and just page in case it is an HTMX request
func RenderHTMX(c *fiber.Ctx, full templ.Component, page templ.Component) error {
	if c.Get("HX-Request") == "true" {
		return RenderJust(c, page)
	}

	return RenderWithBase(c, full)
}
