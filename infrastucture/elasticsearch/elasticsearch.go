package elasticsearch

import (
	"github.com/olivere/elastic/v7"
	"github.com/zsmartex/pkg/v2/log"
)

type Config struct {
	URL      string
	Username string
	Password string
	Sniff    bool
}

type LoggerError struct {
}

func (LoggerError) Printf(format string, v ...interface{}) {
	log.Errorf(format, v...)
}

type LoggerInfo struct {
}

func (LoggerInfo) Printf(format string, v ...interface{}) {
	log.Infof(format, v...)
}

type LoggerTrace struct {
}

func (LoggerTrace) Printf(format string, v ...interface{}) {
	log.Tracef(format, v...)
}

func New(cfg *Config) (*elastic.Client, error) {
	return elastic.NewClient(
		elastic.SetBasicAuth(cfg.Username, cfg.Password),
		elastic.SetGzip(true),
		elastic.SetURL(cfg.URL),
		elastic.SetSniff(cfg.Sniff),
		elastic.SetTraceLog(LoggerError{}),
		elastic.SetInfoLog(LoggerInfo{}),
		elastic.SetInfoLog(LoggerTrace{}),
	)
}
