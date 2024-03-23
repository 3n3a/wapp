package wapp

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// MODULES

// TODO: add default modules for root, error and layout
func DefaultRootModule() *Module {
	m := NewModule(ModuleConfig{
		Name:         "DefaultRootModule",
		InternalName: "",
		IsRoot:       true,
	})

	// TODO: fix
	m.handler = func(c *fiber.Ctx) error {
		return c.Status(200).
			Render("root", nil, "layout")
	}

	return m
}

func DefaultErrorModule() *ErrorModule {
	m := NewErrorModule(ModuleConfig{
		Name:         "DefaultErrorModule",
		InternalName: "error",
		Method:       HTTPMethodAll,
	})

	m.errorHandler = func(c *fiber.Ctx, err error) error {
		// TODO: if datatype x --> return error in type x

		return c.Status(500).
			SendString(fmt.Sprintf("Error: %#v", err))
		// Render("error", nil, "layout")
	}

	return m
}

// HANDLERS

// Default fiber Handler
func DefaultFiberHandler(c *fiber.Ctx) error {
	return c.SendString("hello world: " + c.Path())
}

// ACTIONS

func getAllFormValues(ac *ActionCtx) (map[string]string, error) {
	var out map[string]string

	out = ac.Queries()

	formValues, err := ac.MultipartForm()
	if err != nil {
		return out, nil
	}

	for k, v := range formValues.Value {
		if len(v) > 0 {
			out[k] = v[0]
		}
	}

	return out, nil
}

// Loads Query and Form Key, Value into ActionCtx Store
//
// Name Field: ac.Store.GetString("name")
func ActionLoadFormValues() *Action {
	a := NewAction(func(ac *ActionCtx) error {
		// TOOD: extrapolate into function

		// END TODO

		values, err := getAllFormValues(ac)
		if err != nil {
			return err
		}

		for key, val := range values {
			ac.Store.SetString(key, val)
		}

		return nil
	})

	return a
}
