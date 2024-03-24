package wapp

import (
	"errors"
	"strings"

	"github.com/3n3a/wapp/internal/utils"
	"github.com/gofiber/fiber/v2"
)

// CONFIG

// Wapp Defaults
const (
	DefaultName string = "Wapp"
	DefaultPort uint16 = 3000
	DefaultAddress string = "127.0.0.1"
	DefaultVersion string = "v0.0.1"
	DefaultCoreModules string = "cache,recover,logger,compress"
	DefaultCacheInclude string = "/*"
	DefaultCacheDuration string = "1h"
	DefaultViewsPath string = "frontend/views/"
)

// Module Defaults
const (
	DefaultModuleName   = "Module1"
	DefaultModuleMethod = HTTPMethodGet
)

// MODULES
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
			Render("frontend/views/error", utils.Map{
				"Message": err.Error(),
			})
	}

	return m
}

// HANDLERS

// Default fiber Handler
// TODO: should be default Action...
func DefaultFiberHandler(c *fiber.Ctx) error {
	return c.SendString("hello world: " + c.Path())
}

// ACTIONS





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
	if dataMap, ok := utils.IsMap(data); ok {

		m, err := dataMap.ToXML()
		if err != nil {
			return err
		}
		return ac.XMLWithHeader(m)
	} else {
		return ac.XML(data)
	}
}

func renderJSON(ac *ActionCtx, data interface{}) error {
	if dataMap, ok := utils.IsMap(data); ok {
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
		if dataMap, ok := utils.IsMap(data); ok {
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
