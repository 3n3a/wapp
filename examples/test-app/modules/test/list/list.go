package list

import (
	"github.com/3n3a/wapp"
)

func New() wapp.Module {
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