package wapp

import (
	"encoding/xml"
	"errors"

	"github.com/gofiber/fiber/v2"
)

// MODULES

// TODO: add default modules for root, error and layout
func DefaultRootModule() *Module {
	m := NewModule(ModuleConfig{
		Name:         "DefaultRootModule",
		InternalName: "",
		IsRoot:       true,
	})

	// TODO: fix
	m.handler = func(c *fiber.Ctx) error {
		return c.Status(200).
			Render("root", nil, "layout")
	}

	return m
}

func DefaultErrorModule() *ErrorModule {
	m := NewErrorModule(ModuleConfig{
		Name:         "DefaultErrorModule",
	})

	m.errorHandler = func(c *fiber.Ctx, err error) error {
		// TODO: if datatype x --> return error in type x

		c.Context().Logger().Printf("Error: %#v", err)
		return c.Status(500).
			Render("views/error", nil)
			// SendString(fmt.Sprintf("Error: %#v", err))
	}

	return m
}

// HANDLERS

// Default fiber Handler
func DefaultFiberHandler(c *fiber.Ctx) error {
	return c.SendString("hello world: " + c.Path())
}

// ACTIONS

func getAllFormValues(ac *ActionCtx) (map[string]string, error) {
	var out map[string]string

	out = ac.Queries()

	formValues, err := ac.MultipartForm()
	if err != nil {
		return out, nil
	}

	for k, v := range formValues.Value {
		if len(v) > 0 {
			out[k] = v[0]
		}
	}

	return out, nil
}

// Loads Query and Form Key, Value into ActionCtx Store
//
// Name Field: ac.Store.GetString("name")
func ActionLoadFormValues() *Action {
	a := NewAction(func(ac *ActionCtx) error {
		values, err := getAllFormValues(ac)
		if err != nil {
			return err
		}

		for key, val := range values {
			ac.Store.SetString(key, val)
		}

		return nil
	})

	return a
}

// func isString(val interface{}) (string, bool) {
// 	if str, ok := val.(string); ok {
// 		return str, true
// 	}
// 	return "", false
// }

func isMap(val interface{}) (Map, bool) {
	if m, ok := val.(Map); ok {
		return m, true
	}

	return nil, false
}

type XMLKeyValue struct {
	XMLName xml.Name    `xml:"KeyValue"`
	Name    string      `xml:"key,attr"`
	Value   interface{} `xml:",chardata"`
}

// Create a struct to represent the KeyValues element
type XMLKeyValues struct {
	XMLName   xml.Name      `xml:"KeyValues"`
	KeyValues []XMLKeyValue `xml:"KeyValue"`
}

func transformMapXML(m Map) (XMLKeyValues, error) {
	var kvs XMLKeyValues
	for key, val := range m {
		kvs.KeyValues = append(
			kvs.KeyValues,
			XMLKeyValue{
				Name:	 key,
				Value:   val,
			},
		)
	}
	return kvs, nil
}

// Renders a given Map with a given DataType
func ActionRenderData(dataType DataType, data interface{}, templateName ...string) *Action {
	a := NewAction(func(ac *ActionCtx) error {
		switch dataType {
		case DataTypeHTML:
			if len(templateName) > 0 {
				return ac.Render(templateName[0], data)
			}
			return errors.New("please input a templateName for HTML Data Type")
		case DataTypeJSON:
			return ac.JSON(data)
		case DataTypeXML:
			if dataMap, ok := isMap(data); ok {
				// transform
				m, err := transformMapXML(dataMap)
				if err != nil {
					return err
				}
				return ac.XMLWithHeader(m)
			} else {
				return ac.XML(data)
			}
		}

		return nil
	})

	return a
}
