package utils

import "reflect"

// default map type
// shamelessly stolen from fiber.Map
type Map map[string]interface{}

func StructToMap(obj interface{}) Map {
    objType := reflect.TypeOf(obj)
    objValue := reflect.ValueOf(obj)

    if objType.Kind() != reflect.Struct {
        return nil
    }

    data := make(Map)

    for i := 0; i < objType.NumField(); i++ {
        field := objType.Field(i)
        value := objValue.Field(i).Interface()
        fieldName := field.Name
        data[fieldName] = value
    }

    return data
}

func StructArrToMaps(obj interface{}) []Map {
    sliceValue := reflect.ValueOf(obj)

    if sliceValue.Kind() != reflect.Slice {
        return nil
    }

    maps := make([]Map, 0)

    for i := 0; i < sliceValue.Len(); i++ {
        structInstance := sliceValue.Index(i).Interface()
        mappedStruct := StructToMap(structInstance)
        maps = append(maps, mappedStruct)
    }

    return maps
}