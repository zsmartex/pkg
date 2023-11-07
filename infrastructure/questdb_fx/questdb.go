package questdb_fx

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zsmartex/pkg/v2/config"
	"github.com/zsmartex/pkg/v2/infrastructure/pg"
	"go.uber.org/fx"
)

type questDBParams struct {
	fx.In

	Config config.QuestDB
}

func New(params questDBParams) (*pgxpool.Pool, error) {
	pool, err := pg.New(
		params.Config.Host,
		params.Config.Port,
		params.Config.User,
		params.Config.Pass,
		params.Config.Name,
	)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
