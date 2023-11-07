package questdb_fx

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

var (
	Module = fx.Module("questdb_fx", questdbProviders, questdbInvokes)

	questdbProviders = fx.Provide(New)

	questdbInvokes = fx.Invoke(registerHooks)
)

func registerHooks(lc fx.Lifecycle, db *pgxpool.Pool) {
	lc.Append(fx.StopHook(func() {
		db.Close()
	}))
}
