package test

import (
	"test-app/modules/test/sub"

	"github.com/3n3a/wapp"
)

func New() *wapp.Module {
	// Configure Module
	testModule := wapp.NewModule(wapp.ModuleConfig{
		Name:         "Test Module",
	})

	// Actions
	testModule.AddActions(
		wapp.ActionLoadFormValues(),
		wapp.NewAction(func(ac *wapp.ActionCtx) error {
			err := ac.Store.SetString("url2", "https://"+"google.com")
			return err
		}),
		wapp.NewAction(func(ac *wapp.ActionCtx) error {
			url, _ := ac.Store.GetString("url")
			url2, _ := ac.Store.GetString("url2")
			return ac.SendString(url + url2)
		}),
	)

	// Add Submodules
	testModule.Register(sub.New())

	return testModule
}
