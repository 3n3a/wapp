package wapp

import (
	"github.com/3n3a/wapp/internal/utils"
	"github.com/gofiber/fiber/v2"
)

// CONFIG

// Wapp Defaults
const (
	DefaultName              string = "Wapp"
	DefaultPort              uint16 = 3000
	DefaultAddress           string = "127.0.0.1"
	DefaultVersion           string = "v0.0.1"
	DefaultCoreModules       string = "cache,recover,logger,compress"
	DefaultCacheInclude      string = "/*"
	DefaultCacheDuration     string = "1h"
	DefaultCorsAllowOrigins  string = "*"
	DefaultCorsAllowHeaders  string = ""
	DefaultViewsPath         string = "frontend/views/"
	DefaultMultipleProcesses bool   = false

)

// Module Defaults
const (
	DefaultModuleName   = "Module1"
	DefaultModuleMethod = HTTPMethodGet
)

// MODULES
func DefaultRootModule() Module {
	m := NewModule(ModuleConfig{
		Name:         "DefaultRootModule",
		InternalName: "",
		IsRoot:       true,
	})

	m.handler = func(c *fiber.Ctx) error {
		return c.Status(200).
			Render(DefaultViewsPath + "root", nil)
	}

	return m
}

func DefaultErrorModule() ErrorModule {
	m := NewErrorModule(ModuleConfig{
		Name:       "DefaultErrorModule",
	})

	m.errorHandler = func(c *fiber.Ctx, err error) error {
		// TODO: if datatype x --> return error in type x

		c.Context().Logger().Printf("Error: %#v", err)
		return c.Status(500).
			Render(DefaultViewsPath + "error", utils.Map{
				"Message": err.Error(),
			})
	}

	return m
}

// HANDLERS

// Default fiber Handler
// TODO: should be default Action...
func DefaultFiberHandler(c *fiber.Ctx) error {
	return c.SendString("hello world: " + c.Path())
}
