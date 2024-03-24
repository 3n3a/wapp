package configpassing

import (
	"github.com/gofiber/fiber/v2"
)

// New creates a new middleware handler
func New(config ...Config) fiber.Handler {
	currentConfig := config[0]

	// do not start if no wappconfig
	if currentConfig.WappConfig == nil {
		panic("cannot be used without providing wappConfig")
	}

	// Return new handler
	return func(c *fiber.Ctx) error {
		// Pass wapp config as Local
		c.Locals("_internal", currentConfig.WappConfig)

		// Return from handler
		return c.Next()
	}
}