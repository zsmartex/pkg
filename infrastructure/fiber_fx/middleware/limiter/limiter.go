package limiter

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/zsmartex/pkg/v2/infrastructure/redis_fx"
)

func New(redisClient *redis_fx.Client) func(*fiber.Ctx) error {
	store := &RedisStore{
		redisClient,
	}

	return limiter.New(limiter.Config{
		Max:        10,
		Expiration: 1 * time.Minute,
		Storage:    store,
	})
}
