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
	Max          int
	Expiration   time.Duration
	KeyGenerator func(*fiber.Ctx) string
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
			Max:          config.Max,
			Expiration:   config.Expiration,
			KeyGenerator: config.KeyGenerator,
			LimitReached: func(c *fiber.Ctx) error {
				return ErrRequestLimitExceeded
			},
			Store:                  store,
			SkipFailedRequests:     true,
			SkipSuccessfulRequests: false,
			LimiterMiddleware:      &FlexLimiter{},
		})
	}
}

type FlexLimiter struct{}

func (l *FlexLimiter) New(cfg limiter.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.Locals("remote_ip").(net.IP)

		if ip.IsPrivate() {
			return c.Next()
		}

		return limiter.FixedWindow{}.New(cfg)(c)
	}
}
