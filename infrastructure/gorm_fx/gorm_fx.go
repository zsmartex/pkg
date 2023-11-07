package gorm_fx

import (
	"github.com/zsmartex/pkg/v2/usecase"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var (
	Module = fx.Module("gorm_fx", gormProviders, gormInvokes)

	gormProviders = fx.Provide(New)

	gormInvokes = fx.Invoke(registerHooks)
)

func registerHooks(lc fx.Lifecycle, db *gorm.DB) {
	lc.Append(fx.StartHook(func() {
		usecase.InitCallback(db)
	}))
}
