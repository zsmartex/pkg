package limiter

import (
	"net"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"go.uber.org/fx"

	"github.com/zsmartex/pkg/v2"
	"github.com/zsmartex/pkg/v2/infrastructure/redis_fx"
)

var (
	ErrRequestLimitExceeded = pkg.NewError(429, "request.limit_exceeded", "request limit exceeded")
)

type Limiter func(*fiber.Ctx) error

type limiterParams struct {
	fx.In

	RedisClient *redis_fx.Client
}

func New(params limiterParams) Limiter {
	store := &RedisStore{
		params.RedisClient,
	}

	return limiter.New(limiter.Config{
		Max:        60,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			ip := c.Locals("remote_ip").(net.IP)
			return ip.String()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return ErrRequestLimitExceeded
		},
		Storage: store,
	})
}
