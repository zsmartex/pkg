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

type Config struct {
	Max        int
	Expiration time.Duration
	Prefix     string // optional
}

type Limiter func(config Config) fiber.Handler

type limiterParams struct {
	fx.In

	RedisClient *redis_fx.Client
}

func New(params limiterParams) Limiter {
	store := &RedisStore{
		params.RedisClient,
	}

	return func(config Config) fiber.Handler {
		return limiter.New(limiter.Config{
			Max:        config.Max,
			Expiration: config.Expiration,
			KeyGenerator: func(c *fiber.Ctx) string {
				ip := c.Locals("remote_ip").(net.IP)

				if config.Prefix != "" {
					return config.Prefix + ":" + ip.String()
				}

				return ip.String()
			},
			LimitReached: func(c *fiber.Ctx) error {
				return ErrRequestLimitExceeded
			},
			Store: store,
		})
	}
}
