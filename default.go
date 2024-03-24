package wapp

import (
	"encoding/xml"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// MODULES

// TODO: add default modules for root, error and layout
func DefaultRootModule(w *Wapp) Module {
	m := NewModule(ModuleConfig{
		Name:         "DefaultRootModule",
		InternalName: "",
		IsRoot:       true,
		wappConfig:   &w.config,
	})

	// TODO: fix
	m.handler = func(c *fiber.Ctx) error {
		return c.Status(200).
			Render("root", nil, "layout")
	}

	return m
}

func DefaultErrorModule(w *Wapp) ErrorModule {
	m := NewErrorModule(ModuleConfig{
		Name:       "DefaultErrorModule",
		wappConfig: &w.config,
	})

	m.errorHandler = func(c *fiber.Ctx, err error) error {
		// TODO: if datatype x --> return error in type x

		c.Context().Logger().Printf("Error: %#v", err)
		return c.Status(500).
			Render("frontend/views/error", Map{
				"Message": err.Error(),
			})
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
// func ActionLoadFormValues() *Action {
// 	a := NewAction(func(ac *ActionCtx) error {
// 		values, err := getAllFormValues(ac)
// 		if err != nil {
// 			return err
// 		}

// 		for key, val := range values {
// 			ac.Store.SetString(key, val)
// 		}

// 		return nil
// 	})

// 	return a
// }

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
				Name:  key,
				Value: val,
			},
		)
	}
	return kvs, nil
}

// Renders given Data by the accept header
func ActionRenderDataAccept(data interface{}, templateName ...string) Action {
	a := NewAction(func(ac *ActionCtx) error {
		dataType := DataType(
			ac.Accepts("text/html", "application/json", "text/xml"),
		)

		return dataRenderByType(dataType, templateName, data, ac)
	})

	return a
}

// Renders a given Map with a given DataType
func ActionRenderData(dataType DataType, data interface{}, templateName ...string) Action {
	a := NewAction(func(ac *ActionCtx) error {
		return dataRenderByType(dataType, templateName, data, ac)
	})

	return a
}

func dataRenderByType(dataType DataType, templateName []string, data interface{}, ac *ActionCtx) error {
	switch dataType {
	case DataTypeHTML:
		return renderHTML(templateName, data, ac)
	case DataTypeJSON:
		return renderJSON(ac, data)
	case DataTypeXML:
		return renderXML(data, ac)
	}

	return nil
}

func renderXML(data interface{}, ac *ActionCtx) error {
	if dataMap, ok := isMap(data); ok {

		m, err := transformMapXML(dataMap)
		if err != nil {
			return err
		}
		return ac.XMLWithHeader(m)
	} else {
		return ac.XML(data)
	}
}

func renderJSON(ac *ActionCtx, data interface{}) error {
	if dataMap, ok := isMap(data); ok {
		if dataMap["_internal"] != nil {
			delete(dataMap, "_internal")
		}
		return ac.JSON(dataMap)
	}
	return ac.JSON(data)
}

func renderHTML(templateName []string, data interface{}, ac *ActionCtx) error {
	if len(templateName) > 0 {
		templateName_ := templateName[0]
		if !strings.Contains(templateName_, DefaultViewsPath) {
			templateName_ = DefaultViewsPath + templateName_
		}
		if dataMap, ok := isMap(data); ok {
			dataMap["_internal"] = ac.WappConfig

			// if templateName_ == "table" {
			// 	// TODO: transform key value into table format for displaying ...
			// 	dataMap["KeyValues"] = maps
			// }

			return ac.Status(200).Render(templateName_, dataMap)
		}
		return ac.Status(200).Render(templateName_, data)
	}
	return errors.New("please input a templateName for HTML Data Type")
}
