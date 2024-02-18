package limiter

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/storage/redis/v3"
	"github.com/zsmartex/pkg/v2/config"
)

func New(config config.Redis) func(*fiber.Ctx) error {
	store := redis.New(redis.Config{
		Host:     config.Host,
		Port:     config.Port,
		Password: config.Password,
		Database: 0,
	})

	return limiter.New(limiter.Config{
		Max:        10,
		Expiration: 1 * time.Minute,
		Storage:    store,
	})
}
