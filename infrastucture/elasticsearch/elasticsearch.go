package elasticsearch

import "github.com/olivere/elastic/v7"

type Config struct {
	URL      string
	Username string
	Password string
	Sniff    bool
}

func New(cfg *Config) (*elastic.Client, error) {
	return elastic.NewClient(elastic.SetBasicAuth(cfg.Username, cfg.Password), elastic.SetGzip(true), elastic.SetURL(cfg.URL), elastic.SetSniff(cfg.Sniff))
}
