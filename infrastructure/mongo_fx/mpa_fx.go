package mongo_fx

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
)

type Tabler interface {
	TableName() string
}

var (
	Module = fx.Module("mongo_fx", mongoProviders, mongoInvokes)

	mongoProviders = fx.Provide(New)

	mongoInvokes = fx.Invoke(registerHooks)
)

func registerHooks(lc fx.Lifecycle, db *mongo.Database) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if err := db.Client().Disconnect(ctx); err != nil {
				return err
			}

			return nil
		},
	})
}
