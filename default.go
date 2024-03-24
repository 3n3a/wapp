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
	DefaultViewsPath         string = "frontend/views/"
	DefaultMultipleProcesses bool   = false
)

// Module Defaults
const (
	DefaultModuleName   = "Module1"
	DefaultModuleMethod = HTTPMethodGet
)

// MODULES
func DefaultRootModule(w *Wapp) Module {
	m := NewModule(ModuleConfig{
		Name:         "DefaultRootModule",
		InternalName: "",
		IsRoot:       true,
		wappConfig:   &w.config,
	})

	// TODO: fix
	m.handler = func(c *fiber.Ctx) error {
		return c.Status(200).
			Render("root", nil, "layout")
	}

	return m
}

func DefaultErrorModule(w *Wapp) ErrorModule {
	m := NewErrorModule(ModuleConfig{
		Name:       "DefaultErrorModule",
		wappConfig: &w.config,
	})

	m.errorHandler = func(c *fiber.Ctx, err error) error {
		// TODO: if datatype x --> return error in type x

		c.Context().Logger().Printf("Error: %#v", err)
		return c.Status(500).
			Render("frontend/views/error", utils.Map{
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
