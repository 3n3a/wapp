package test

import (
	"fmt"
	"test-app/modules/test/sub"

	"github.com/3n3a/wapp"
	"github.com/muonsoft/validation/validate"
)

func New() wapp.Module {
	// Configure Module
	testModule := wapp.NewModule(wapp.ModuleConfig{
		Name: "Test Module",
	})

	// Action
	testModule.AddAction(
		wapp.NewAction(func(ac *wapp.ActionCtx) error {
			// input
			inputUrl := ac.FormValue("url", "")

			// transform / data
			//// validate url
			err := validate.URL(inputUrl)
			urlValid := err == nil

			fmt.Printf("route: %#v\n", ac.Locals("_internal"))

			// output / render
			return ac.RenderDataByAcceptHeader(
				wapp.Map{
					"url": inputUrl,
					"valid": urlValid,
				},
				"table",
			)
		}),
	)

	// Add Submodules
	testModule.Register(sub.New())

	return testModule
}


