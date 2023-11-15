//go:generate mockgen -source=repository.go -destination=mock/repository.go -package=mock_questdb

package questdb_fx

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"gorm.io/gorm/schema"
)

type Repository[T schema.Tabler] interface {
	Exec(ctx context.Context, sql string, attrs ...interface{}) error
	Query(ctx context.Context, dst interface{}, sql string, attrs ...interface{}) error
	QueryRow(ctx context.Context, dst interface{}, sql string, attrs ...interface{}) error
}

type repository[T schema.Tabler] struct {
	conn *pgxpool.Pool
}

func NewRepository[T schema.Tabler](conn *pgxpool.Pool) Repository[T] {
	return repository[T]{
		conn: conn,
	}
}

func (u repository[T]) Exec(ctx context.Context, sql string, attrs ...interface{}) error {
	_, err := u.conn.Exec(ctx, sql, attrs...)
	if err != nil {
		return err
	}

	return nil
}

func (u repository[T]) Query(ctx context.Context, dst interface{}, sql string, attrs ...interface{}) error {
	err := pgxscan.Select(ctx, u.conn, dst, sql, attrs...)

	return errors.Wrap(err, "questdb query")
}

func (u repository[T]) QueryRow(ctx context.Context, dst interface{}, sql string, attrs ...interface{}) error {
	err := pgxscan.Get(ctx, u.conn, dst, sql, attrs...)

	return errors.Wrap(err, "questdb query row")
}
