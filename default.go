package wapp

import "github.com/gofiber/fiber/v2"

// TODO: add default modules for root, error and layout
func DefaultRootModule() *Module {
	m := NewModule(ModuleConfig{
		Name: "DefaultRootModule",
		PathName: "",
		IsRoot: true,
	})

	// TODO: fix
	m.handler = func(c *fiber.Ctx) error {
		return c.Status(200).
			Render("root", nil, "layout")
	}

	return m
}

func DefaultErrorModule() *Module {
	m := NewModule(ModuleConfig{
		Name: "DefaultErrorModule",
		PathName: "error",
		Method: "ALL",
	})

	// TODO: fix rendering
	m.handler = func(c *fiber.Ctx) error {
		return c.Status(404).
			Render("error", nil, "layout")
	}

	return m
}

// Default fiber Handler
func DefaultFiberHandler(c *fiber.Ctx) error {
	return c.SendString("hello world: " + c.Path())
}