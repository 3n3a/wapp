package sub

import (
	"github.com/3n3a/wapp"
)

func New() *wapp.Module {
	testSubModule := wapp.NewModule(wapp.ModuleConfig{
		Name:         "TestSubModule",
		InternalName: "sub",
		Method:       "POST",
	})
	testSubModule.AddActions(
		wapp.NewAction(func(ac *wapp.ActionCtx) error {
			return ac.SendString("Test hello sub")
		}),
	)

	return testSubModule
}
