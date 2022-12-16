package pg

import (
	"context"
	"fmt"
	"os"

	pgxDecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
	"github.com/zsmartex/pkg/v2/log"
)

type Logger struct {
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
	logger := log.Logger
	if data != nil {
		logger = logger.WithContext(ctx).WithFields(data)
	} else {
		logger = logger.WithContext(ctx)
	}

	logger.Info(data)

	switch level {
	case tracelog.LogLevelTrace:
		logger.Trace(msg)
	case tracelog.LogLevelDebug:
		logger.Debug(msg)
	case tracelog.LogLevelInfo:
		logger.Info(msg)
	case tracelog.LogLevelWarn:
		logger.Warn(msg)
	case tracelog.LogLevelError:
		logger.Error(msg)
	}
}

func New(
	host string,
	port int,
	user string,
	password string,
	dbname string,
) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user, password, host, port, dbname)

	pgxConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	pgxConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())
		pgxDecimal.Register(conn.TypeMap())

		return nil
	}

	var logLevel tracelog.LogLevel
	switch os.Getenv("LOG_LEVEL") {
	case "WARN":
		logLevel = tracelog.LogLevelWarn
	case "INFO":
		logLevel = tracelog.LogLevelInfo
	case "DEBUG":
		logLevel = tracelog.LogLevelDebug
	case "ERROR":
		logLevel = tracelog.LogLevelError
	case "FATAL":
		logLevel = tracelog.LogLevelError
	case "PANIC":
		logLevel = tracelog.LogLevelError
	case "TRACE":
		logLevel = tracelog.LogLevelTrace
	}

	pgxConfig.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   NewLogger(),
		LogLevel: logLevel,
	}

	pgxConnPool, err := pgxpool.NewWithConfig(context.TODO(), pgxConfig)
	if err != nil {
		return nil, err
	}

	return pgxConnPool, nil
}
