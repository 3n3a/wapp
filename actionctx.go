package wapp

import (
	"encoding/xml"
	"errors"
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
func (ac *ActionCtx) renderDataByDataType(dataType DataType, data interface{}, templateName []string) error {
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

func (ac *ActionCtx) renderXML(data interface{}) error {
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

func (ac *ActionCtx) renderJSON(data interface{}) error {
	if dataMap, ok := utils.IsMap(data); ok {
		if dataMap["_internal"] != nil {
			delete(dataMap, "_internal")
		}
		return ac.JSON(dataMap)
	}
	return ac.JSON(data)
}

func (ac *ActionCtx) renderHTML(data interface{}, templateName []string) error {
	if len(templateName) > 0 {
		templateName_ := templateName[0]

		// prepend viewspath if not already contained
		if !strings.Contains(templateName_, DefaultViewsPath) {
			templateName_ = DefaultViewsPath + templateName_
		}

		if dataMap, ok := utils.IsMap(data); ok {
			dataMap["_internal"] = ac.Locals("_internal")

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

// Chooses output method based on Accept header values
//
// Wrapper around ActionCtx.renderDataByDataType()
//
// See under enum.go for possible DataType's
func (ac *ActionCtx) RenderDataByAcceptHeader(data interface{}, templateName ...string) error {
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

func (ac *ActionCtx) RenderData(dataType DataType, data interface{}, templateName ...string) error {
	return ac.renderDataByDataType(dataType, data, templateName)
}
