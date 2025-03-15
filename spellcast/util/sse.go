package util

import (
	"bufio"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

// Wrapper around bufio.Writer for easier sending of server sent events.
type EventWriter struct {
	writer *bufio.Writer
}

// Send a server sent event.
func (w *EventWriter) SendEvent(event string, data string) error {

	// Convert to a server sent events message
	var msg string
	if event == "" {
		msg = fmt.Sprintf("data: %s\n\n", data)
	} else {
		msg = fmt.Sprintf("event: %s\ndata: %s\n\n", event, data)
	}

	// Send to the browser
	if _, err := fmt.Fprint(w.writer, msg); err != nil {
		return err
	}

	return w.writer.Flush()
}

// A helper function to make server sent events easier
func StartEvents(c *fiber.Ctx, handler func(w *EventWriter)) error {

	// Set the appropriate headers
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	// Start sending events
	c.Status(fiber.StatusOK).Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {

		// Create a new event writer
		writer := &EventWriter{
			writer: w,
		}

		// Call the handler for the server sent events
		handler(writer)
	}))

	return nil
}
