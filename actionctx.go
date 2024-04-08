package wapp

import (
	"encoding/xml"
	"errors"
	"slices"

	"strings"

	"github.com/3n3a/wapp/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/huandu/go-clone"
)

// Wrapper around fiber.Ctx
// extended with wrapper functions,
// aswell as other convenience attributes etc.
type ActionCtx struct {
	*fiber.Ctx

	ModuleConfig ModuleConfig
	WappConfig   Config
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
func (ac *ActionCtx) renderDataByDataType(dataType DataType, data []utils.Map, displayColumns []string, templateName []string) error {
	// fmt.Printf("%s: %#v\n", dataType, data)

	// remove columns not for display
	noInternal := []DataType{DataTypeJSON, DataTypeXML}
	for _, currMap := range data {
		for key, _ := range currMap {
			if len(displayColumns) > 0 && !slices.Contains[[]string, string](displayColumns, key) {
				delete(currMap, key)
			}

			if strings.HasPrefix(key, "_") && slices.Contains[[]DataType, DataType](noInternal, dataType) {
				// internal AND specific datatype
				delete(currMap, key)
			}
		}
	}

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
	return ac.JSON(data)
}

func (ac *ActionCtx) renderHTML(data []utils.Map, templateName []string) error {
	if len(templateName) > 0 {
		templateName_ := templateName[0]

		// prepend viewspath if not already contained
		if !strings.Contains(templateName_, DefaultViewsPath) {
			templateName_ = DefaultViewsPath + templateName_
		}

		c := clone.Clone(ac.WappConfig).(Config) // deep copy
		c.Menu.CurrentPath = ac.Path()

		moduleConfig := clone.Clone(ac.ModuleConfig).(ModuleConfig) // deep copy

		for i, field := range moduleConfig.UIFields {
			moduleConfig.UIFields[i].Default = ac.FormValue(field.Name, field.Default)
		}

		dataMap := utils.Map{
			"values":    data,
			"_internal": c,
			"_module":   moduleConfig,
		}

		// fmt.Printf("%#v\n", dataMap)

		if ac.Get("HX-Boosted", "false") == "true" {
			// only send back "embed part"
			return ac.Status(200).Render(templateName_, dataMap, DefaultViewsPath+"layout-hx")
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
func (ac *ActionCtx) RenderDataByAcceptHeader(data []utils.Map, displayColumns []string, templateName ...string) error {
	// This is the part that devices which DataType gets used
	// specifically the Accepts() function
	dataType := DataType(
		ac.Accepts(
			string(DataTypeHTML),
			string(DataTypeJSON),
			string(DataTypeXML),
		),
	)

	return ac.renderDataByDataType(dataType, data, displayColumns, templateName)
}

// RenderData based on DataType
//
// Set "data.values" field to []utils.Map for "table" template
func (ac *ActionCtx) RenderData(dataType DataType, data []utils.Map, templateName ...string) error {
	return ac.renderDataByDataType(dataType, data, nil, templateName)
}
