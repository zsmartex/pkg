//go:generate mockgen -source=repository.go -destination=mock/repository.go -package=mock_questdb

package questdb_fx

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

var _ Repository = (*repository)(nil)

type Repository interface {
	Exec(ctx context.Context, sql string, attrs ...interface{}) error
	Query(ctx context.Context, dst interface{}, sql string, attrs ...interface{}) error
	QueryRow(ctx context.Context, dst interface{}, sql string, attrs ...interface{}) error
}

type repository struct {
	conn *pgxpool.Pool
}

func NewRepository(conn *pgxpool.Pool) repository {
	return repository{
		conn: conn,
	}
}

func (u repository) Exec(ctx context.Context, sql string, attrs ...interface{}) error {
	_, err := u.conn.Exec(ctx, sql, attrs...)
	if err != nil {
		return err
	}

	return nil
}

func (u repository) Query(ctx context.Context, dst interface{}, sql string, attrs ...interface{}) error {
	err := pgxscan.Select(ctx, u.conn, dst, sql, attrs...)

	return errors.Wrap(err, "questdb query")
}

func (u repository) QueryRow(ctx context.Context, dst interface{}, sql string, attrs ...interface{}) error {
	err := pgxscan.Get(ctx, u.conn, dst, sql, attrs...)

	return errors.Wrap(err, "questdb query row")
}
