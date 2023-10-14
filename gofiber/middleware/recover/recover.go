package recover

import (
	"runtime"

	"github.com/gofiber/fiber/v2"
	"github.com/zsmartex/pkg/v2"
	"github.com/zsmartex/pkg/v2/log"
)

func defaultStackTraceHandler(_ *fiber.Ctx, e interface{}) {
	buf := make([]byte, 2048)
	buf = buf[:runtime.Stack(buf, false)]
	log.Errorf("Panic: %v\n%s\n", e, string(buf))
}

// New creates a new middleware handler
func New(config ...Config) fiber.Handler {
	// Set default config
	cfg := configDefault(config...)

	// Return new handler
	return func(c *fiber.Ctx) (err error) {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Catch panics
		defer func() {
			if r := recover(); r != nil {
				if cfg.EnableStackTrace {
					cfg.StackTraceHandler(c, r)
				}

				err = pkg.ErrServerInternal
			}
		}()

		// Return err if existed, else move to next handler
		return c.Next()
	}
}
