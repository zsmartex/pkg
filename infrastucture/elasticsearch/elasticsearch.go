package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v8"
)

type Config struct {
	URL      []string
	Username string
	Password string
}

func New(cfg *Config) (*elasticsearch.Client, error) {
	return elasticsearch.NewClient(elasticsearch.Config{
		Addresses: cfg.URL,
		Username:  cfg.Username,
		Password:  cfg.Password,
	})
}
