package wapp

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type DataType string

const (
	DataTypeHTML DataType = "HTML"
	DataTypeJSON DataType = "JSON"
	DataTypeXML  DataType = "XML"
)

type HTTPMethod string

const (
	HTTPMethodAll     HTTPMethod = "ALL"
	HTTPMethodGet     HTTPMethod = "GET"
	HTTPMethodPost    HTTPMethod = "POST"
	HTTPMethodPut     HTTPMethod = "PUT"
	HTTPMethodGDelete HTTPMethod = "DELETE"
	HTTPMethodHead    HTTPMethod = "HEAD"
	HTTPMethodConnect HTTPMethod = "CONNECT"
	HTTPMethodOptions HTTPMethod = "OPTIONS"
	HTTPMethodTrace   HTTPMethod = "TRACE"
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
}

// Module is the basic building block
// in this framework
//
// # To create a new Module
//
// wapp.NewModule(...ModuleConfig) *Module
type Module struct {
	// PreActions are Actions executed
	// before main data action and is thought to
	// be used for initial input transformations
	//
	// Default: []
	PreActions []*Action

	// Actions are the main actions
	// in here are data read/write operations
	// but also calculations / transformations
	//
	// Default: []
	Actions []*Action

	// PostActions are actions which
	// are executed after main actions
	// meant for transformations of output
	// be that rendering or changing structure
	//
	// Default: []
	PostActions []*Action

	// Module config (aka Metadata about and around module)
	config ModuleConfig
	// Submodules list
	submodules []*Module
	// Handler is the content that will be returned
	// TODO: remove this and replace with actions wrapper framework
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

type ErrorModule struct {
	*Module

	errorHandler fiber.ErrorHandler
}

// do intial module stuff
func (m *Module) init() {
	// create initial pathname array
	// TODO: make sure runs on each update --> value could theoreticall change
	if m.config.PathName != "" {
		m.config.fullPath = append(m.config.fullPath, m.config.PathName)
	}
}

// Register adds configured Submodule
func (m *Module) Register(module *Module) {
	module.config.fullPath = append(m.config.fullPath, module.config.PathName)
	m.submodules = append(m.submodules, module)
}

func (m *Module) AddPreActions(action ...*Action) {
	m.AddActions(ActionTypePre, action...)
}

func (m *Module) AddMainActions(action ...*Action) {
	m.AddActions(ActionTypeMain, action...)
}

func (m *Module) AddPostActions(action ...*Action) {
	m.AddActions(ActionTypePost, action...)
}

// Adds one or many actions to the specified array
func (m *Module) AddActions(actionType ActionType, action ...*Action) {
	for _, currentAction := range action {
		if actionType == ActionTypePre {
			m.PreActions = append(m.PreActions, currentAction)
		}
		if actionType == ActionTypeMain {
			m.Actions = append(m.Actions, currentAction)
		}
		if actionType == ActionTypePost {
			m.Actions = append(m.Actions, currentAction)
		}
	}
}

func (m *Module) buildHandler() {
	// bundle together the actions and create one function AKA the handler
	m.handler = func(c *fiber.Ctx) error {
		logger := log.New(os.Stdout, "MODULE", log.Lshortfile)
		actionCtx := &ActionCtx{
			Ctx:   c,
			Store: NewKV(),
		}

		// pre actions
		for _, a := range m.PreActions {
			err := a.f(actionCtx) // call func in action
			if err != nil {
				
				logger.Fatalf("Error: %#v\n", err)
				return err
			}
		}

		// main actions
		for _, a := range m.Actions {
			err := a.f(actionCtx) // call func in action
			if err != nil {
				logger.Fatalf("Error: %#v\n", err)
				return err
			}
		}

		// post actions
		for _, a := range m.PostActions {
			err := a.f(actionCtx) // call func in action
			if err != nil {
				logger.Fatalf("Error: %#v\n", err)
				return err
			}
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

const (
	DefaultModuleName   = "Module1"
	DefaultModuleMethod = HTTPMethodGet
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
	if strings.ContainsAny(mod.config.PathName, "/") {
		// TODO: should i error out like this?
		log.Fatal(errors.New("\"/\" (slashes) are not allowed in path name"))
	}

	// run init
	mod.init()

	return mod
}

func NewErrorModule(moduleConfigs ...ModuleConfig) *ErrorModule {
	// Create a new module
	defaultModule := NewModule(moduleConfigs...)

	errorModule := &ErrorModule{
		Module: defaultModule,
	}

	return errorModule
}
