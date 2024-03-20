package main

import (
	"github.com/3n3a/wapp"
)

func main() {

	// with config
	w := wapp.New(wapp.Config{
		CoreModules: []wapp.CoreModule{
			wapp.Recover,
			wapp.Logger,
		},
	})
	// "/test" and "/test/"
	testModule := wapp.NewModule(wapp.ModuleConfig{
		Name: "TestModule", 
		PathName: "test",
	})

	// Pre
	testModule.AddPreActions(wapp.ActionLoadFormValues())

	// Main
	testModule.AddMainActions(wapp.NewAction(func(ac *wapp.ActionCtx) error {
		err := ac.Store.SetString("url2", "https://" + "google.com")
		return err
	}))

	// Post
	testModule.AddPostActions(wapp.NewAction(func(ac *wapp.ActionCtx) error {
		url, _ := ac.Store.GetString("url")
		url2, _ := ac.Store.GetString("url2")
		return ac.SendString(url + url2)
	}))

	// "/test/sub"
	testSubModule := wapp.NewModule(wapp.ModuleConfig{
		Name: "TestSubModule", 
		PathName: "sub", 
		Method: "POST",
	})
	testSubModule.AddActions(wapp.ActionTypePost, wapp.NewAction(func(ac *wapp.ActionCtx) error {
		return ac.SendString("Test hello sub")
	}))

	testModule.Register(testSubModule)
	w.Register(testModule)
	w.Start()
}
