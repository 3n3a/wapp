package wapp

import (
	// "reflect"
	"testing"

	"github.com/stretchr/testify/require"
	// "github.com/stretchr/testify/assert"
)

// TODO: write more unit tests

func Test_Wapp_NewDefaultConfig(t *testing.T) {
	w := New()
	
	require.Equal(t, DefaultPort, w.config.Port, "Port Default")
	require.Equal(t, DefaultAddress, w.config.Address, "Address Default")
	require.Equal(t, DefaultVersion, w.config.Version, "Version Default")
	require.Equal(t, CoreModulesFromString(DefaultCoreModules), w.config.CoreModules, "CoreModules Default")
	
	// Do not actually start
	w.init()
}

func Test_Wapp_CustomConfigBasic(t *testing.T) {
	w := New(Config{
		Port:    4000,
		Address: "0.0.0.0",
		Version: "v1.0.1",
		CoreModules: []CoreModule{
			Logger,
			Recover,
		},
	})
	
	require.Equal(t, uint16(4000), w.config.Port, "Port Custom")
	require.Equal(t, "0.0.0.0", w.config.Address, "Address Custom")
	require.Equal(t, "v1.0.1", w.config.Version, "Version Custom")
	require.Equal(t, CoreModulesFromString("logger,recover"), w.config.CoreModules, "CoreModules Custom")
	
	// Do not actually start
	w.init()
}
