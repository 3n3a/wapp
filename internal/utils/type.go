package utils

// IsString - tests if "val" is of type "string"
// 
// When is string: returns val cast into string
func IsString(val interface{}) (bool) {
	if _, ok := val.(string); ok {
		return true
	}
	return false
}

// IsMap - tests if "val" is of type "Map"
// 
// When is map: returns val cast into Map
func IsMap(val interface{}) (bool) {
	if _, ok := val.(Map); ok {
		return true
	}

	return false
}