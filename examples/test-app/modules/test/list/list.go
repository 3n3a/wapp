package list

import (
	"github.com/3n3a/wapp"
)

func New() wapp.Module {
	// TODO: how can i add input fields??
	// as a user i want to add list of fields
	// which will result in form submit with get
	listModule := wapp.NewModule(wapp.ModuleConfig{
		Name: "List",
	})

	listModule.AddAction(
		wapp.NewAction(func(ac *wapp.ActionCtx) error {
			list := []wapp.Map{
				{
					"test1": "test1",
					"test2": "test2",
					"test3": "test3",
				},
				{
					"test1": "test1",
					"test2": "test2",
					"test3": "test3",
				},
				{
					"test1": "test1",
					"test2": "test2",
					"test3": "test3",
					"test4": "test4",
				},
			}

			return ac.RenderDataByAcceptHeader(
				list,
				"table",
			)
		}),
	)

	return listModule
}