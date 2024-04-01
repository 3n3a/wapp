package wapp

import (
	"github.com/3n3a/wapp/internal/utils"
)

// exported stuff for external use

// export map
type Map = utils.Map

// export struct to map func
func StructToMap(obj interface{}) utils.Map {
	return utils.StructToMap(obj)
}

// export struct array to map array func
func StructArrToMaps(obj interface{}) []utils.Map {
	return utils.StructArrToMaps(obj)
}