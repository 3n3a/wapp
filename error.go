package wapp

import (
	"github.com/gofiber/fiber/v2"
)

type ErrorModule struct {
	*Module

	errorHandler fiber.ErrorHandler
}

func NewErrorModule(moduleConfigs ...ModuleConfig) ErrorModule {
	// Create a new module
	defaultModule := NewModule(moduleConfigs...)

	errorModule := ErrorModule{
		Module: &defaultModule,
	}

	return errorModule
}
