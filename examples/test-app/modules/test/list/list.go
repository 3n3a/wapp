package list

import (
	"github.com/3n3a/wapp"
)

func New() wapp.Module {
	listModule := wapp.NewModule(wapp.ModuleConfig{
		Name: "List",
		UIInputTitle: "Filter",
		UIOutputTitle: "Response",
		UIFields: []wapp.UIField{
			{
				Name: "filter",
				Type: wapp.UITypeDropdown,
				Required: true,
				Children: []wapp.UIChild{
					{
						Value: "filter_1",
						DisplayValue: "Filter 1",
					},
					{
						Value: "filter_2",
						DisplayValue: "Filter 2",
					},
				},
			},
			{
				Name: "Get Response",
				Type: wapp.UITypeSubmit,
			},
		},
	})

	listModule.AddAction(
		wapp.NewAction(func(ac *wapp.ActionCtx) error {
			filter := ac.FormValue("filter")

			var list []wapp.Map
			if filter == "filter_1" {
				list = []wapp.Map{
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

			} else {
				list = []wapp.Map{
					{
						"test1": "Hello",
						"test2": "Filter",
						"test3": "2",
					},
				}
			}

			return ac.RenderDataByAcceptHeader(
				list,
				"in_out_table",
			)
		}),
	)

	return listModule
}