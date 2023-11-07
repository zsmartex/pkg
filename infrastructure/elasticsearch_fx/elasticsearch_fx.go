package elasticsearch_fx

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/zsmartex/pkg/v2/config"
	"go.uber.org/fx"
)

var (
	Module = fx.Module("elasticsearch_fx", elasticsearchProviders)

	elasticsearchProviders = fx.Provide(New)
)

type elasticsearchParams struct {
	fx.In

	Config config.Elasticsearch
}

func New(params elasticsearchParams) (*elasticsearch.Client, error) {
	return elasticsearch.NewClient(elasticsearch.Config{
		Addresses: params.Config.URL,
		Username:  params.Config.Username,
		Password:  params.Config.Password,
	})
}
