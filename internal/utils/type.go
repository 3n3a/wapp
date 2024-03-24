package utils

// IsString - tests if "val" is of type "string"
// 
// When is string: returns val cast into string
func IsString(val interface{}) (string, bool) {
	if str, ok := val.(string); ok {
		return str, true
	}
	return "", false
}

// IsMap - tests if "val" is of type "Map"
// 
// When is map: returns val cast into Map
func IsMap(val interface{}) (Map, bool) {
	if m, ok := val.(Map); ok {
		return m, true
	}

	return nil, false
}