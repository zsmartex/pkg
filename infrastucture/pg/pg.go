package pg

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	pgxDecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
	"github.com/zsmartex/pkg/v3/log"
)

type QueryTracer struct {
	pgx.QueryTracer
}

type traceQueryData struct {
	startTime time.Time
	sql       string
	args      []any
}

func (q *QueryTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	return context.WithValue(ctx, 1, &traceQueryData{
		startTime: time.Now(),
		sql:       data.SQL,
		args:      data.Args,
	})
}

func (q *QueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	queryData := ctx.Value(1).(*traceQueryData)

	endTime := time.Now()
	interval := endTime.Sub(queryData.startTime)

	sql := queryData.sql
	space := regexp.MustCompile(`\s+`)
	sql = space.ReplaceAllString(sql, " ")
	sql = strings.TrimSpace(sql)
	for i, v := range queryData.args {
		sql = strings.Replace(sql, fmt.Sprintf("$%d", i+1), fmt.Sprint(v), 1)
	}

	if data.Err != nil {
		log.Error(data.Err)
		log.Errorf("%s [%s]", sql, interval)
		return
	}

	log.Tracef("%s [%s]", sql, interval)
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

	pgxConfig.ConnConfig.Tracer = &QueryTracer{}

	pgxConnPool, err := pgxpool.NewWithConfig(context.TODO(), pgxConfig)
	if err != nil {
		return nil, err
	}

	return pgxConnPool, nil
}
