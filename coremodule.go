package wapp

import "strings"

// Enum values for available core modules
//
// These can be added to config to active for example
// the caching or so
type CoreModule string

const (
	Cache CoreModule = "cache"
	Recover CoreModule = "recover"
	Logger CoreModule = "logger"
	Compress CoreModule = "compress"
)

// CoreModulesFromString - converts comma-separated list of core-modules
// back into list of CoreModules
func CoreModulesFromString(list string) []CoreModule {
	coreModuleList := make([]CoreModule, 0)
	coreModuleValues := strings.Split(list, ",")

	for _, value := range coreModuleValues {
		coreModuleList = append(coreModuleList, CoreModule(value))
	}

	return coreModuleList
}