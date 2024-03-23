package wapp

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/smirzaei/parallel"
)

//go:embed views/* public/*
var viewsfs embed.FS

// builds on gofiber, tailwindcss
//
// like a smart wrapper around that
//
// which wraps around resty for easy html client
// and XX for easy db
// and also offers custom functions which are generally useful

// Config is a struct holding the wapp configuration
type Config struct {
	
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
	// Default: 3600s - 1 hour
	CacheDuration string `json:"cache_duration"`
}

// Wapp is the main object for interacting with this library
type Wapp struct {
	// Wapp config
	config Config
	// Fiber Instance
	ffiber *fiber.App
	// Root Module
	rootModule *Module // provides layout and container for submodules
	// Error Module
	errorModule *ErrorModule
}

// init executes initial functions for wapp
func (w *Wapp) init() {
	// TODO: what needs to go here: server?
	fiberConfig := fiber.Config{}

	// initialize with default, but allow override
	// TODO: how to override module defaults
	w.rootModule = DefaultRootModule()
	w.errorModule = DefaultErrorModule()

	// create html engine
	engine := html.NewFileSystem(
		http.FS(viewsfs),
		".html",
	)

	engine.AddFunc("getCssAsset", func(name string) (res template.HTML) {
		filepath.Walk("public/assets", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Name() == name {
				res = template.HTML("<link rel=\"stylesheet\" href=\"/" + path + "\">")
			}
			return nil
		})
		return
	})

	engine.AddFunc("getJsAsset", func(name string) (res template.HTML) {
		filepath.Walk("public/assets", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Name() == name {
				res = template.HTML("<script src=\"/" + path + "\"></script>")
			}
			return nil
		})
		return
	})

	engine.AddFunc("getCssInline", func(name string) (res template.HTML) {
		filepath.Walk("public/assets", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Name() == name {
				data, err := ioutil.ReadFile(path)
				if err != nil {
					return err	
				}
				res = template.HTML("<style>" + string(data) + "</style>")
			}
			return nil
		})
		return
	})

	engine.AddFunc("getJsInline", func(name string) (res template.HTML) {
		filepath.Walk("public/assets", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Name() == name {
				data, err := ioutil.ReadFile(path)
				if err != nil {
					return err	
				}
				res = template.HTML("<script>" + string(data) + "</script>")
			}
			return nil
		})
		return
	})

	// TODO: how could i emebed these?
	fiberConfig.Views = engine
	fiberConfig.ViewsLayout = "views/layout"

	// Error Handling
	fiberConfig.ErrorHandler = w.errorModule.errorHandler

	// TODO: allow custom fiber config
	fiberConfig.ServerHeader = "Wapp"
	w.ffiber = fiber.New(fiberConfig)
	
	// TODO: process all core modules
	if slices.Contains(w.config.CoreModules, Cache) {
		cacheDuration, err := time.ParseDuration(w.config.CacheDuration)
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

	// Static
	// TODO: embed and add path to config
	w.ffiber.Static("/public", "public")
}

// recursively process all submodules and create tree
func (w *Wapp) processModules(modules []*Module) {
	parallel.ForEach(modules, func(currModule *Module) {
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
	w.processModules([]*Module{w.rootModule})

	// Start server
	hostPort := net.JoinHostPort(w.config.Address, fmt.Sprint(w.config.Port))
	log.Fatal(
		w.ffiber.Listen(hostPort),
	)
}

// Register adds a configured module to wapp
func (w *Wapp) Register(module *Module) {
	// Because at root level add as submodule to root
	w.rootModule.Register(module)
}

// Default values when not provided
const (
	DefaultPort uint16 = 3000
	DefaultAddress string = "127.0.0.1"
	DefaultVersion string = "v0.0.1"
	DefaultCoreModules string = "cache,recover,logger,compress"
	DefaultCacheInclude string = "/*"
	DefaultCacheDuration string = "1h"
)

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

	// Init wapp
	wapp.init()

	return wapp
}
