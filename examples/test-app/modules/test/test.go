package test

import (
	"test-app/modules/test/sub"

	"github.com/3n3a/wapp"
)

func New() *wapp.Module {
	// Configure Module
	testModule := wapp.NewModule(wapp.ModuleConfig{
		Name: "Test Module",
	})

	// Actions
	testModule.AddActions(
		wapp.ActionLoadFormValues(),
		wapp.NewAction(func(ac *wapp.ActionCtx) error {
			err := ac.Store.SetString("url2", "https://"+"google.com")
			return err
		}),
		wapp.ActionRenderData(wapp.DataTypeHTML, wapp.Map{
			"url": "sfdlfkjdlfj",
			"url2": "skldjflk213123l",
		}, "frontend/views/root"),
	)

	// Add Submodules
	testModule.Register(sub.New())

	return testModule
}
