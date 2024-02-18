package limiter

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/zsmartex/pkg/v2"
	"github.com/zsmartex/pkg/v2/infrastructure/redis_fx"
	"go.uber.org/fx"
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
		Max:        20,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			ip := string(c.Locals("remote_ip").([]byte))
			return ip
		},
		LimitReached: func(c *fiber.Ctx) error {
			return ErrRequestLimitExceeded
		},
		Storage: store,
	})
}
