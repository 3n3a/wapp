package wapp

import (
	"errors"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"github.com/smirzaei/parallel"
)

// ModuleConfig is the configuration of a module
type ModuleConfig struct {
	// Name is the debug name of this module
	//
	// Default: "Module1"
	Name string `json:"name"`

	// InternalName is the technical name of this module
	//
	// Default: Slugified version of "Name" 
	InternalName string `json:"internal_name"`

	// IsRoot defines if this is the root module
	//
	// Default: false
	IsRoot bool `json:"is_root"`

	// Method is the HTTP request method
	//
	// Default: HTTPMethodGet
	Method HTTPMethod `json:"method"`

	// Datatype is the type of date being returned
	//
	// Default: "HTML"
	DataType string `json:"data_type"`

	// FullPath is the full path for registering as handler
	//
	// Calculated
	fullPath []string `json:"-"`

	// TODO: how can i add input fields??
	// as a user i want to add list of fields
	// which will result in form submit with get

	UIInputTitle string `json:"ui_input_title"`

	UIOutputTitle string `json:"ui_output_title"`

	UIFields []UIField `json:"ui_fields"`
}

type UIFieldType string
const (
	UITypeDropdown UIFieldType = "dropdown"
	UITypeText UIFieldType = "text"
	UITypeEmail UIFieldType = "email"
	UITypePhone UIFieldType = "phone"
	UITypeNumber UIFieldType = "number"

	UITypeSubmit UIFieldType = "submit"
)

type UIChild struct {
	Value string `json:"value"`
	DisplayValue string `json:"display_value"`
}

type UIField struct {
	Name string `json:"name"`
	Type UIFieldType `json:"type"`
	Default string `json:"default"`
	Required bool `json:"required"`
	Children []UIChild `json:"children"`
}

// Module is the basic building block
// in this framework
//
// # To create a new Module
//
// wapp.NewModule(...ModuleConfig) *Module
type Module struct {
	// Action is a thin wrapper around fiber.Handler
	//
	// Default: DefaultFiberHandler
	Action Action

	// Module config (aka Metadata about and around module)
	config ModuleConfig
	// Submodules list
	submodules []Module
	// Internal Function for handling a route
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
	// TODO: should be executed each time InternalName updated
	if m.config.InternalName != "" {
		m.config.fullPath = append(m.config.fullPath, m.config.InternalName)
	}
}

// Get a config value
func (m *Module) GetConfig() ModuleConfig {
	return m.config
}

// Register adds configured Submodule
func (m *Module) Register(module ...Module) {
	parallel.ForEach(module, func(el Module) {
		// create full path
		el.config.fullPath = append(m.config.fullPath, el.config.InternalName)
		// add to list of submodules
		m.submodules = append(m.submodules, el)
	})
}

// Adds an Action to a Module
// without one it is without use
func (m *Module) AddAction(action Action) {
	m.Action = action
}

func (m *Module) buildHandler() {
	// bundle together the actions and create one function AKA the handler
	m.handler = func(c *fiber.Ctx) error {
		logger := c.Context().Logger()

		// build action ctx
		// contains extended functions
		// and attributes
  // TODO: include ref to current module,
  // access in action
		actionCtx := &ActionCtx{
			Ctx:   c,
		}

		// main actions
		err := m.Action.f(actionCtx) // call func in action

		// here to print error in action
		if err != nil {
			logger.Printf("error in action: %#v\n", err)
			return err
		}

		return nil
	}
}

// OnBeforeProcess Executed when module is processed
func (m *Module) OnBeforeProcess() {
	// Handle all generation cases
	// based on values that can be
	// configured before bein made into fiber app

	// create handler function
	m.buildHandler()
}

func (m *Module) GetFullPath() string {
	return "/" + strings.Join(m.config.fullPath, "/")
}

func NewModule(moduleConfigs ...ModuleConfig) Module {
	// Create a new module
	mod := Module{
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
	if mod.config.InternalName == "" && !mod.config.IsRoot {
		// Default is slugified version of Name
		mod.config.InternalName = slug.Make(mod.config.Name)
	}
	if strings.ContainsAny(mod.config.InternalName, "/") {
		// TODO: should i error out like this?
		log.Fatal(errors.New("\"/\" (slashes) are not allowed in path name"))
	}
	if mod.config.UIInputTitle == "" {
		mod.config.UIInputTitle = DefaultModuleUIInputTitle
	}
	if mod.config.UIOutputTitle == "" {
		mod.config.UIOutputTitle = DefaultModuleUIOutputTitle
	}
	if mod.config.UIFields == nil || len(mod.config.UIFields) == 0 {
		mod.config.UIFields = make([]UIField, 0)
	}

	// run init
	mod.init()

	return mod
}
