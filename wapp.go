package wapp

import (
	"bytes"
	"crypto/sha256"
	"embed"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"
	"slices"
	"strings"
	"time"

	configpassing "github.com/3n3a/wapp/internal/middleware/config-passing"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/template/html/v2"
	"github.com/smirzaei/parallel"
)

//go:embed frontend/views/*
var viewsfs embed.FS

//go:embed frontend/public/* frontend/public/assets/*
var staticfs embed.FS

// builds on gofiber, tailwindcss
//
// like a smart wrapper around that
//
// which wraps around resty for easy html client
// and XX for easy db
// and also offers custom functions which are generally useful

type MenuNode struct {
	Name string `json:"name"`
	FullPath string `json:"full_path"`

	SubNodes []MenuNode `json:"sub_nodes"`
	
	CurrentPath string `json:"current_path"`
}

func MenuSetCurrentPath(m MenuNode, current string) MenuNode {
	m.CurrentPath = current
	return m
}

// Config is a struct holding the wapp configuration
type Config struct {
	// Name is the name of the application
	//
	// Default "Wapp"
	Name string `json:"name"`

	// Port is the port the application will be listening at
	//
	// Default: 3000
	Port uint16 `json:"port"`

	// Address to listen on
	//
	// Default: "127.0.0.1"
	Address string `json:"address"`

	// Version is the version you give your app
	//
	// Default: "v0.0.1"
	Version string `json:"version"`

	// Enabled Core Modules
	//
	// Default: [Cache, Recover, Logger, Compress]
	CoreModules []CoreModule `json:"core_modules"`

	// Included Paths for caching
	//
	// Default: "/*" - all paths
	CacheInclude []string `json:"cache_include"`

	// Cache Duration for above paths
	//
	// Default: "1h" (3600s - 1 hour)
	CacheDuration string `json:"cache_duration"`

	// Cors allow headers
	//
	// Default: []
	CorsAllowHeaders []string `json:"cors_allow_headers"`

	// Cors allow origins
	//
	// Default: ["*"]
	CorsAllowOrigins []string `json:"cors_allow_origins"`

	// Multiple Processes
	//
	// Sets the fiber Prefork Option.
	// Make sure you start from a shell (Docker `CMD ./app` or `CMD ["sh", "-c", "/app"]`)
	//
	// Default: false
	MultipleProcesses bool `json:"multiple_processes"`

	// internal menu tree
	Menu MenuNode `json:"menu"`

	// path - Name Map
	pathNameMap map[string]Module `json:"-"`

	// Debug Mode
	//
	// Activates multiple features regarding debugging
	// Framework
	//
	// Default: false
	DebugMode bool `json:"debug_mode"`

 // TODO: Base ActionCtx with ref to WappConfig and ref to Current Module (added in Module when creating handler)
}

func (c *Config) GetCurrentModule(curr string) Module {
	return c.pathNameMap[curr]
}

// Wapp is the main object for interacting with this library
type Wapp struct {
	// Wapp config
	config Config
	// Fiber Instance
	ffiber *fiber.App
	// Root Module
	rootModule Module // provides layout and container for submodules
	// Error Module
	errorModule ErrorModule
}

// init executes initial functions for wapp
func (w *Wapp) init() {
	fiberConfig := fiber.Config{
		// Print all routes with methods on startup
		EnablePrintRoutes: true,
	}

	// initialize with default, but allow override
	// TODO: how to override module defaults
	w.rootModule = DefaultRootModule()
	w.errorModule = DefaultErrorModule()

	// create html engine
	engine := html.NewFileSystem(
		http.FS(viewsfs),
		".html",
	)

	// Additional Functions for use in template
	engine.Funcmap["hasPrefix"] = strings.HasPrefix
	engine.Funcmap["MenuSetCurrentPath"] = MenuSetCurrentPath

	fiberConfig.Views = engine
	fiberConfig.ViewsLayout = "frontend/views/layout"

	// Error Handling
	fiberConfig.ErrorHandler = w.errorModule.errorHandler

	// TODO: allow custom fiber config
	fiberConfig.ServerHeader = "Wapp"
	fiberConfig.AppName = w.config.Name + " " + w.config.Version
	
	// Encoders
	fiberConfig.XMLEncoder = func(v interface{}) ([]byte, error) {
		prefix := ""
		indent := "    "
		return xml.MarshalIndent(v, prefix, indent)
	}
	
	w.ffiber = fiber.New(fiberConfig)

	cacheDuration, err := time.ParseDuration(w.config.CacheDuration)
	if slices.Contains(w.config.CoreModules, Cache) {
		if err != nil {
			log.Fatal(errors.New("please provide valid cache duration"))
		}
		w.ffiber.Use(cache.New(cache.Config{
			Next: func(c *fiber.Ctx) bool {
				for _, pathMatch := range w.config.CacheInclude {
					match, _ := regexp.MatchString(pathMatch, c.Path())
					if match {
						// fmt.Println("cached: ", match)
						return false // cached
					}
				}
				return true // not cached
			},
			Expiration:   cacheDuration,
			CacheControl: true,
			KeyGenerator: func(c *fiber.Ctx) string {
				// uses hash to prevent big urls / headers from affecting us
				ogUrl := utils.CopyString(c.OriginalURL())
				ogAcpt := utils.CopyString(c.Get("accept"))
				hxBoosted := utils.CopyString(c.Get("hx-boosted", ""))
				hxCurrent := utils.CopyString(c.Get("hx-current-url", ""))
				hxReq := utils.CopyString(c.Get("hx-request", ""))
				bytes := bytes.Join(
					[][]byte{
						[]byte(ogUrl),
						[]byte(ogAcpt),
						[]byte(hxBoosted),
						[]byte(hxCurrent),
						[]byte(hxReq),
					},
					[]byte("."),
				)
				key := sha256.New()
				key.Write(bytes)
				hash := key.Sum(nil)
				hashStr := fmt.Sprintf("%x", hash)
				// fmt.Println("key: ", hashStr)
				return hashStr
			},
		}))
	}
	if slices.Contains(w.config.CoreModules, Recover) {
		w.ffiber.Use(recover.New())
	}
	if slices.Contains(w.config.CoreModules, Logger) {
		w.ffiber.Use(logger.New())
	}
	if slices.Contains(w.config.CoreModules, Compress) {
		w.ffiber.Use(compress.New())
	}
	if slices.Contains(w.config.CoreModules, CORS) {
		w.ffiber.Use(cors.New(cors.Config{
			AllowOrigins: strings.Join(w.config.CorsAllowOrigins, ", "),
			AllowHeaders: strings.Join(w.config.CorsAllowHeaders, ", "),
		}))
	}

	w.ffiber.Use(configpassing.New[Config](configpassing.Config[Config]{
		WappConfig: &w.config,
	}))

	// embedded fs
	w.ffiber.Use("/public", filesystem.New(filesystem.Config{
		Root:       http.FS(staticfs),
		PathPrefix: "frontend/public",
		Browse:     false,
		MaxAge:     int(cacheDuration.Seconds()), // Same as cache duration for Cache plugin
	}))

	// Static
	w.ffiber.Static("/public", "public")
}

// recursively process all submodules and create tree
// executed at startup :)
func (w *Wapp) processModules(modules []Module) []MenuNode {
	menuNodes := []MenuNode{}

	parallel.ForEach(modules, func(currModule Module) {
		currModule.OnBeforeProcess()

		// ...processing
		if currModule.config.Method == HTTPMethodAll {
			w.ffiber.All(
				currModule.GetFullPath(),
				currModule.handler,
			)
		} else {
			// add handler
			w.ffiber.Add(
				string(currModule.config.Method),
				currModule.GetFullPath(),
				currModule.handler,
			)
		}

		if currModule.config.IsRoot {
			// set name to app name
			currModule.config.Name = w.config.Name
		}

		// menu
		currNode := MenuNode{
			Name: currModule.config.Name,
			FullPath: currModule.GetFullPath(),
		}

		// add to map for easy retrieval
		w.config.pathNameMap[currNode.FullPath] = currModule

		// process submodules
		subNodes := w.processModules(currModule.submodules)
		currNode.SubNodes = subNodes
		
		menuNodes = append(menuNodes, currNode)
	})

	return menuNodes
}

// Start needs to be executed
// after registering all the modules
func (w *Wapp) Start() {
	// process root
	w.config.pathNameMap = make(map[string]Module)

	rootNodes := w.processModules([]Module{w.rootModule})
	w.config.Menu =  rootNodes[0]
	w.config.Menu.SubNodes = append([]MenuNode{
		{
			Name: "Home",
			FullPath: "/",
		},
	}, w.config.Menu.SubNodes...)

	// Start server
	hostPort := net.JoinHostPort(w.config.Address, fmt.Sprint(w.config.Port))
	log.Fatal(
		w.ffiber.Listen(hostPort),
	)
}

// Register adds a configured module to wapp
func (w *Wapp) Register(module ...Module) {
	// Because at root level add as submodule to root
	w.rootModule.Register(module...)
}

// New creates a new wapp named instance
func New(config ...Config) *Wapp {
	// Create a new Wapp
	wapp := &Wapp{
		// Create Config
		config: Config{},
	}

	// Override config if provided
	if len(config) > 0 {
		wapp.config = config[0]
	}

	// Initial default / override values (if invalid config)
	if wapp.config.Name == "" {
		wapp.config.Name = DefaultName
	}
	if wapp.config.Port == 0 {
		wapp.config.Port = DefaultPort
	}
	if wapp.config.Address == "" {
		wapp.config.Address = DefaultAddress
	}
	if wapp.config.Version == "" {
		wapp.config.Version = DefaultVersion
	}
	if len(wapp.config.CoreModules) == 0 {
		// TODO: config: ability to override default configs for above + wrapper for easier config (but still able to go completly custom)
		wapp.config.CoreModules = CoreModulesFromString(DefaultCoreModules)
	}
	if len(wapp.config.CacheInclude) == 0 {
		wapp.config.CacheInclude = strings.Split(DefaultCacheInclude, ",")
	}
	if wapp.config.CacheDuration == "" {
		wapp.config.CacheDuration = DefaultCacheDuration
	}
	if wapp.config.MultipleProcesses == false {
		wapp.config.MultipleProcesses = DefaultMultipleProcesses
	}
	if wapp.config.CorsAllowHeaders == nil {
		wapp.config.CorsAllowHeaders = []string{DefaultCorsAllowHeaders}
	}
	if wapp.config.CorsAllowOrigins == nil {
		wapp.config.CorsAllowOrigins = []string{DefaultCorsAllowOrigins}
	}
	if wapp.config.DebugMode == false {
		wapp.config.DebugMode = DefaultDebugMode
	}	

	// Init wapp
	wapp.init()

	return wapp
}
