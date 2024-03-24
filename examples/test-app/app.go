package main

import (
	"test-app/modules/test"

	"github.com/3n3a/wapp"
)

func main() {
	// with config
	w := wapp.New(wapp.Config{
		Name: "Test Wapp",
		CoreModules: []wapp.CoreModule{
			wapp.Recover,
			wapp.Logger,
			wapp.CORS,
		},
	})
	
	// Register Lowest Level Modules (not Submodules)
	w.Register(test.New())

	// Start
	w.Start()
}
