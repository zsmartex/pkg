package mongo_fx

import (
	"sync"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBlogger struct {
	mu sync.Mutex
}

func (logger *DBlogger) Info(level int, msg string, keyandvalues ...interface{}) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	if options.LogLevel(level+1) == options.LogLevelDebug {
		logrus.Debugf("message: %s value:%v", msg, keyandvalues)
	} else {
		logrus.Infof("message: %s value:%v", msg, keyandvalues)
	}
}
func (logger *DBlogger) Error(err error, msg string, keyandvalues ...interface{}) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logrus.Errorf("error: %v, message: %s, value:%v", err, msg, keyandvalues)
}
