package views

import (
	"context"
	"log"
	"os"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

// Custom type for context keys to avoid collisions
type URLContextKey struct{}

// Render the component surrounded by html head and footer and stuff
func RenderWithBase(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	ctx := context.WithValue(c.Context(), URLContextKey{}, templ.SafeURL(c.Path()))
	return Base(component).Render(ctx, c.Response().BodyWriter())
}

// Render just the component without the base
func RenderJust(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	ctx := context.WithValue(c.Context(), URLContextKey{}, templ.SafeURL(c.Path()))
	return component.Render(ctx, c.Response().BodyWriter())
}

// Render just the component without the base (but you can pass an error for logging purposes)
func RenderError(c *fiber.Ctx, component templ.Component, err error) error {

	if err != nil && os.Getenv("MAGIC_TESTING") == "true" {
		log.Println(c.Path()+":", err)
	}

	// Render the component like normal
	return RenderJust(c, component)
}

// Render full in case of no HTMX and just page in case it is an HTMX request
func RenderHTMX(c *fiber.Ctx, full templ.Component, page templ.Component) error {

	if c.Get("HX-Request") == "true" {
		return RenderJust(c, page)
	}

	return RenderWithBase(c, full)
}
