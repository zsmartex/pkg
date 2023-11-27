package mongo_fx

import (
	"context"
	"fmt"

	"github.com/zsmartex/pkg/v2/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
)

type mongoParams struct {
	fx.In

	Config config.DB
}

func New(params mongoParams) (*mongo.Database, error) {
	sink := &DBlogger{}

	loggerOptions := options.
		Logger().
		SetSink(sink).
		SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)

	options := options.Client().
		ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%d",
			params.Config.User,
			params.Config.Pass,
			params.Config.Host,
			params.Config.Port,
		)).
		SetLoggerOptions(loggerOptions)

	client, err := mongo.Connect(context.TODO(), options)
	if err != nil {
		return nil, err
	}

	return client.Database(params.Config.Name), nil
}
