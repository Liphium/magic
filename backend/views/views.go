package views

import (
	"context"
	"log"
	"os"

	"github.com/Liphium/magic/backend/util/constants"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

// Custom type for context keys to avoid collisions
type URLContextKey struct{}

// Custom type for context keys to avoid collisions
type PermissionLevelContextKey struct{}

// Render the component surrounded by html head and footer and stuff
func RenderWithBase(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	ctx := buildContext(c)
	return Base(component).Render(ctx, c.Response().BodyWriter())
}

// Render just the component without the base
func RenderJust(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	ctx := buildContext(c)
	return component.Render(ctx, c.Response().BodyWriter())
}

// Build context for templ with all the nessecary values
func buildContext(c *fiber.Ctx) context.Context {
	ctx := context.WithValue(c.Context(), URLContextKey{}, templ.SafeURL(c.Path()))
	if c.Locals(constants.LocalsPermissionLevel) != nil {
		ctx = context.WithValue(ctx, PermissionLevelContextKey{}, c.Locals(constants.LocalsPermissionLevel).(uint))
	}
	return ctx
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
