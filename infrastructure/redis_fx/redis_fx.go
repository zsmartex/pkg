package redis_fx

import (
	"go.uber.org/fx"
)

var (
	Module = fx.Module("redis_fx", redisProviders, redisInvokes)

	redisProviders = fx.Provide(New)

	redisInvokes = fx.Invoke(registerHooks)
)

func registerHooks(lc fx.Lifecycle, redisClient *Client) {
	lc.Append(fx.StopHook(redisClient.Close))
}
