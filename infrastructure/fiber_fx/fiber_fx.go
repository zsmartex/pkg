package fiber_fx

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/helmet/v2"
	"github.com/zsmartex/pkg/v2/config"
	"github.com/zsmartex/pkg/v2/infrastructure/fiber_fx/middleware/error_handler"
	"github.com/zsmartex/pkg/v2/infrastructure/fiber_fx/middleware/logger"
	"github.com/zsmartex/pkg/v2/infrastructure/fiber_fx/middleware/recover"
	"go.uber.org/fx"
)

var (
	Module = fx.Module(
		"fiber_fx",
		fiberProviders,
		fiberInvokes,
	)

	fiberProviders = fx.Provide(
		New,
	)

	fiberInvokes = fx.Options(fx.Invoke(registerHooks))
)

type fiberParams struct {
	fx.In

	ApplicationName string `name:"application_name"`
}

func New(params fiberParams, lc fx.Lifecycle) *fiber.App {
	fiberApp := fiber.New(fiber.Config{
		BodyLimit:                10 * 1024 * 1024, // this is the default limit of 10MB
		EnableTrustedProxyCheck:  true,
		ProxyHeader:              "X-Forwarded-For",
		TrustedProxies:           []string{},
		AppName:                  params.ApplicationName,
		ErrorHandler:             error_handler.ErrorHandler,
		EnableSplittingOnParsers: true,
	})

	fiberApp.Use(compress.New())
	fiberApp.Use(helmet.New())
	fiberApp.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))
	fiberApp.Use(requestid.New())
	fiberApp.Use(logger.New())
	fiberApp.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	return fiberApp
}

type hookParams struct {
	fx.In

	FiberApp *fiber.App
	Config   config.HTTP `name:"http_server"`
}

func registerHooks(lc fx.Lifecycle, params hookParams) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				err := params.FiberApp.Listen(params.Config.Address())
				if err != nil {
					log.Fatal(err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return params.FiberApp.ShutdownWithContext(ctx)
		},
	})
}
