package wapp

import (
	"encoding/xml"

	"github.com/gofiber/fiber/v2"
)

// Wrapper around fiber.Ctx
// extended with wrapper functions, 
// aswell as other convenience attributes etc.
type ActionCtx struct {
	*fiber.Ctx

	WappConfig *Config
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

func (ac *ActionCtx) 