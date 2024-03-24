package configpassing

import (
	"github.com/gofiber/fiber/v2"
)

// New creates a new middleware handler
func New[T any](config Config[T]) fiber.Handler {
	// Return new handler
	return func(c *fiber.Ctx) error {

		// Pass wapp config as Local, if html
		if c.Accepts("text/html") != "" {
			c.Locals("_internal", config.WappConfig)
		}

		// Debug
		// fmt.Printf("%#v\n", c.Locals("_internal"))

		// Return from handler
		return c.Next()
	}
}