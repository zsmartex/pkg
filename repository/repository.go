package repository

import (
	"context"

	"github.com/cockroachdb/errors"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/zsmartex/pkg/v2/gpa"
)

type Repository[T schema.Tabler] interface {
	DB() *gorm.DB
	TableName() string
	WithTrx(trxHandle *gorm.DB) Repository[T]
	Count(context context.Context, filters ...gpa.Filter) (int, error)
	First(context context.Context, dst interface{}, filters ...gpa.Filter) error
	Last(context context.Context, dst interface{}, filters ...gpa.Filter) error
	Find(context context.Context, dst interface{}, filters ...gpa.Filter) error
	Transaction(handler func(tx *gorm.DB) error) error
	FirstOrCreate(context context.Context, dst interface{}, filters ...gpa.Filter) error
	Create(context context.Context, dst interface{}, filters ...gpa.Filter) error
	Updates(context context.Context, dst interface{}, value interface{}, filters ...gpa.Filter) error
	UpdateColumns(context context.Context, dst interface{}, value interface{}, filters ...gpa.Filter) error
	Delete(context context.Context, dst interface{}, filters ...gpa.Filter) error
	Raw(context context.Context, sql string, values ...interface{}) *gorm.DB
	Exec(context context.Context, sql string, values ...interface{}) *gorm.DB
}

type repository[T schema.Tabler] struct {
	repository gpa.Repository
	entity     T
}

func New[T schema.Tabler](db *gorm.DB, entity T) Repository[T] {
	return repository[T]{
		repository: gpa.New(db, entity),
	}
}

func (r repository[T]) DB() *gorm.DB {
	return r.repository.DB
}

func (r repository[T]) TableName() string {
	return r.entity.TableName()
}

func (r repository[T]) WithTrx(trxHandle *gorm.DB) Repository[T] {
	r.repository = r.repository.WithTrx(trxHandle)

	return r
}

func (r repository[T]) Count(context context.Context, filters ...gpa.Filter) (int, error) {
	count, err := r.repository.Count(context, filters...)
	if err != nil {
		return 0, errors.Wrap(err, "repository count")
	}

	return count, nil
}

func (r repository[T]) First(context context.Context, model interface{}, filters ...gpa.Filter) error {
	err := r.repository.First(context, model, filters...)

	return errors.Wrap(err, "repository first")
}

func (r repository[T]) Last(context context.Context, model interface{}, filters ...gpa.Filter) error {
	err := r.repository.Last(context, model, filters...)

	return errors.Wrap(err, "repository last")
}

func (r repository[T]) Find(context context.Context, models interface{}, filters ...gpa.Filter) error {
	err := r.repository.Find(context, models, filters...)

	return errors.Wrap(err, "repository find")
}

func (r repository[T]) Transaction(handler func(tx *gorm.DB) error) error {
	err := r.repository.DB.Transaction(handler)

	return errors.Wrap(err, "repository transaction")
}

func (r repository[T]) FirstOrCreate(context context.Context, model interface{}, filters ...gpa.Filter) error {
	err := r.repository.FirstOrCreate(context, model, filters...)

	return errors.Wrap(err, "repository first or create")
}

func (r repository[T]) Create(context context.Context, model interface{}, filters ...gpa.Filter) error {
	err := r.repository.Create(context, model, filters...)

	return errors.Wrap(err, "repository create")
}

func (r repository[T]) Updates(context context.Context, model interface{}, value interface{}, filters ...gpa.Filter) error {
	err := r.repository.Updates(context, model, value, filters...)

	return errors.Wrap(err, "repository update")
}

func (r repository[T]) UpdateColumns(context context.Context, model interface{}, value interface{}, filters ...gpa.Filter) error {
	err := r.repository.UpdateColumns(context, model, value, filters...)

	return errors.Wrap(err, "repository update columns")
}

func (r repository[T]) Delete(context context.Context, model interface{}, filters ...gpa.Filter) error {
	err := r.repository.Delete(context, model, filters...)

	return errors.Wrap(err, "repository delete")
}

func (r repository[T]) Raw(context context.Context, sql string, values ...interface{}) *gorm.DB {
	return r.repository.Raw(context, sql, values...)
}

func (r repository[T]) Exec(context context.Context, sql string, values ...interface{}) *gorm.DB {
	return r.repository.Exec(context, sql, values...)
}
