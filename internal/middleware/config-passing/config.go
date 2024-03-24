package configpassing

// "github.com/gofiber/fiber/v2"

// Config defines the config for middleware.
type Config[T any] struct {
	WappConfig *T
}
