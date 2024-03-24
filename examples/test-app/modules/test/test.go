package test

import (
	"test-app/modules/test/sub"

	"github.com/3n3a/wapp"
)

func New() wapp.Module {
	// Configure Module
	testModule := wapp.NewModule(wapp.ModuleConfig{
		Name: "Test Module",
	})

	// Actions
	testModule.AddAction(
		wapp.NewAction(func(ac *wapp.ActionCtx) error {
			return ac.SendString("test test test")
		}),
	)

	// Add Submodules
	testModule.Register(sub.New())

	return testModule
}
