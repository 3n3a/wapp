package wapp

import (
	"errors"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ModuleConfig is the configuration of a module
type ModuleConfig struct {
	// Name is the debug name of this module
	//
	// Default: "Module1"
	Name string `json:"name"`

	// PathName is the part of the path for this module, without slashes
	//
	// Required
	PathName string `json:"path_name"`

	// IsRoot defines if this is the root module
	//
	// Default: false
	IsRoot bool `json:"is_root"`

	// Method is the HTTP request method
	//
	// Default: "GET"
	Method string `json:"method"`

	// FullPath is the full path for registering as handler
	//
	// Calculated
	fullPath []string `json:"-"`
}

// Module is the basic building block
// in this framework
//
// To create a new Module
// 
// wapp.NewModule(...ModuleConfig) *Module
type Module struct {
	// Module config (aka Metadata about and around module)
	config ModuleConfig
	// Submodules list
	submodules []*Module
	// Handler is the content that will be returned
	handler fiber.Handler
	// TODO: add module contents
	// TODO: each module has option for a menu (would be rendered on main page)
	// TODO: each module has a route
	// TODO: each module has ability to override default css
	// TODO: each module has options for input (url params, form values, json body, xml body)
	// TODO: each module has options for input validation (url, ip-address, text, html-safe whatever)
	// TODO: each module has options for data transformation, before & after retrieve (be able to provide function)
	// TODO: each module has a data retrieve function (db, http: rest, json, xml, html)
	// TODO: each module has options for output (html-page, html-part, json, xml)
}

// do intial module stuff
func (m *Module) init() {
	// create initial pathname array
	if m.config.PathName != "" {
		m.config.fullPath = append(m.config.fullPath, m.config.PathName)
	}
}

// Register allows adding a configured Submodule
func (m *Module) Register(module *Module) {
	module.config.fullPath = append(m.config.fullPath, module.config.PathName)
	m.submodules = append(m.submodules, module)
}

func (m *Module) GetFullPath() string {
	return "/" + strings.Join(m.config.fullPath, "/")
}

const (
	DefaultModuleName = "Module1"
	DefaultModuleMethod = "GET"
)

func NewModule(moduleConfigs ...ModuleConfig) *Module {
	// Create a new module
	mod := &Module{
		// Create Module Config
		config: ModuleConfig{},
	}

	// Override config if provided
	if len(moduleConfigs) > 0 {
		mod.config = moduleConfigs[0]
	}

	// Default values
	if mod.config.Name == "" {
		mod.config.Name = DefaultModuleName
	}
	if mod.config.Method == "" {
		mod.config.Method = DefaultModuleMethod
	}
	if mod.config.PathName == "" && !mod.config.IsRoot {
		// TODO: should i error out like this?
		log.Fatal(errors.New("missing pathname in module" + mod.config.Name))
	}
	// TODO: switch to configurable "actions" 
	// where handler is wrapped for user
	if mod.handler == nil {
		mod.handler = DefaultFiberHandler
	}
	if strings.ContainsAny(mod.config.PathName, "/") {
		// TODO: should i error out like this?
		log.Fatal(errors.New("\"/\" (slashes) are not allowed in path name"))
	}

	// run init
	mod.init()

	return mod
}