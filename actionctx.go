package wapp

import (
	"encoding/xml"
	"errors"
	// "fmt"
	"strings"

	"github.com/3n3a/wapp/internal/utils"
	"github.com/gofiber/fiber/v2"
)

// Wrapper around fiber.Ctx
// extended with wrapper functions,
// aswell as other convenience attributes etc.
type ActionCtx struct {
	*fiber.Ctx
}

// Render XML with the xml.Header as a first line
func (ac *ActionCtx) XMLWithHeader(data interface{}) error {
	raw, err := ac.App().Config().XMLEncoder(data)
	if err != nil {
		return err
	}

	// With Header
	raw = []byte(xml.Header + string(raw))

	ac.Context().Response.SetBodyRaw(raw)
	ac.Context().Response.Header.SetContentType(fiber.MIMEApplicationXML)
	return nil
}

// internal function that renders given data by a type
func (ac *ActionCtx) renderDataByDataType(dataType DataType, data []utils.Map, templateName []string) error {
	// fmt.Printf("%s: %#v\n", dataType, data)

	switch dataType {
	case DataTypeHTML:
		return ac.renderHTML(data, templateName)
	case DataTypeJSON:
		return ac.renderJSON(data)
	case DataTypeXML:
		return ac.renderXML(data)
	}

	return nil
}

func (ac *ActionCtx) renderXML(data []utils.Map) error {
	list := utils.XMLList{}
	for _, curr := range data {
		m, err := curr.ToXML()
		if err != nil {
			return err
		}

		list.Items = append(list.Items, m)

	}
	return ac.XMLWithHeader(list)
}

func (ac *ActionCtx) renderJSON(data []utils.Map) error {
	if dataMap, ok := utils.IsMap(data); ok {
		if dataMap["_internal"] != nil {
			delete(dataMap, "_internal")
		}
		return ac.JSON(dataMap)
	}
	return ac.JSON(data)
}

func (ac *ActionCtx) renderHTML(data []utils.Map, templateName []string) error {
	if len(templateName) > 0 {
		templateName_ := templateName[0]

		// prepend viewspath if not already contained
		if !strings.Contains(templateName_, DefaultViewsPath) {
			templateName_ = DefaultViewsPath + templateName_
		}

		c := *ac.Locals("_internal").(*Config)
		c.Menu.CurrentPath = ac.Path()
		currModule := c.GetCurrentModule(c.Menu.CurrentPath)

		for i, field := range currModule.config.UIFields {
			currModule.config.UIFields[i].Default = ac.FormValue(field.Name, field.Default)
		}

		dataMap := utils.Map{
			"values":    data,
			"_internal": c,
			"_module": currModule.GetConfig(),
		}

		if ac.Get("HX-Boosted", "false") == "true" {
			// only send back "embed part"
			return ac.Status(200).Render(templateName_, dataMap, DefaultViewsPath + "layout-hx")
		} else {
			return ac.Status(200).Render(templateName_, dataMap)
		}

	}
	return errors.New("please input a templateName for HTML Data Type")
}

// Chooses output method based on Accept header values
//
// Set "data.values" field to []utils.Map for "table" template
//
// Wrapper around ActionCtx.renderDataByDataType().
// See under enum.go for possible DataType's.
func (ac *ActionCtx) RenderDataByAcceptHeader(data []utils.Map, templateName ...string) error {
	// This is the part that devices which DataType gets used
	// specifically the Accepts() function
	dataType := DataType(
		ac.Accepts(
			string(DataTypeHTML),
			string(DataTypeJSON),
			string(DataTypeXML),
		),
	)

	return ac.renderDataByDataType(dataType, data, templateName)
}

// RenderData based on DataType
//
// Set "data.values" field to []utils.Map for "table" template
func (ac *ActionCtx) RenderData(dataType DataType, data []utils.Map, templateName ...string) error {
	return ac.renderDataByDataType(dataType, data, templateName)
}
