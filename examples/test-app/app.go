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
	testModule := wapp.NewModule(wapp.ModuleConfig{Name: "TestModule", PathName: "test"})
	// "/test/sub"
	testSubModule := wapp.NewModule(wapp.ModuleConfig{Name: "TestSubModule", PathName: "sub", Method: "POST"})
	testModule.Register(testSubModule)
	w.Register(testModule)
	w.Start()
}
