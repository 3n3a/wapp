package wapp

import (
	"encoding/xml"

	"github.com/gofiber/fiber/v2"
)

type Map = fiber.Map

type ActionCtx struct {
	*fiber.Ctx

	Store *KV
}

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

type ActionFunc = func(*ActionCtx) error

// Main Action Container
type Action struct {
	// function that is executed when you know
	f ActionFunc
}

// NewAction creates and initializes a new Action
func NewAction(f ActionFunc) *Action {
	action := &Action{}

	action.f = f

	return action
}