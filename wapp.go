package wapp

import (
	"embed"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"
	"slices"
	"strings"
	"time"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
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

	// Multiple Processes
	//
	// Sets the fiber Prefork Option.
	// Make sure you start from a shell (Docker `CMD ./app` or `CMD ["sh", "-c", "/app"]`)
	//
	// Default: false
	MultipleProcesses bool `json:"multiple_processes"`
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
	w.rootModule = DefaultRootModule(w)
	w.errorModule = DefaultErrorModule(w)

	// create html engine
	engine := html.NewFileSystem(
		http.FS(viewsfs),
		".html",
	)

	fiberConfig.Views = engine
	fiberConfig.ViewsLayout = "frontend/views/layout"

	// Error Handling
	fiberConfig.ErrorHandler = w.errorModule.errorHandler

	// TODO: allow custom fiber config
	fiberConfig.ServerHeader = "Wapp"
	w.ffiber = fiber.New(fiberConfig)
	
	// TODO: process all core modules
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
						return false // cached
					}
				}
				return true // not cached
			},
			Expiration: cacheDuration,
			CacheControl: true,
			KeyGenerator: func(c *fiber.Ctx) string {
				return c.OriginalURL()
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

	// embedded fs
	w.ffiber.Use("/public", filesystem.New(filesystem.Config{
		Root: http.FS(staticfs),
		PathPrefix: "frontend/public",
		Browse: false,
		MaxAge: int(cacheDuration.Seconds()), // Same as cache duration for Cache plugin
	}))

	// Static
	w.ffiber.Static("/public", "public")
}

// recursively process all submodules and create tree
func (w *Wapp) processModules(modules []Module) {
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

		// TODO: additional processing for module like menu building etc.
	
		// process submodules
		w.processModules(currModule.submodules)
	})
}

// Start needs to be executed 
// after registering all the modules
func (w *Wapp) Start() {
	// process root
	w.processModules([]Module{w.rootModule})

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
		config:		Config{},
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

	// Init wapp
	wapp.init()

	return wapp
}
