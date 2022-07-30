package repository

import (
	"context"

	"github.com/zsmartex/pkg/v2/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/zsmartex/pkg/v2/gpa"

	"github.com/zsmartex/pkg/v2/log"
)

type Repository[T schema.Tabler] interface {
	Count(context context.Context, filters ...gpa.Filter) (int, error)
	First(context context.Context, dst interface{}, filters ...gpa.Filter) error
	Last(context context.Context, dst interface{}, filters ...gpa.Filter) error
	Find(context context.Context, dst interface{}, filters ...gpa.Filter) error
	WithTrx(trxHandle *gorm.DB) Repository[T]
	Transaction(handler func(tx *gorm.DB) error) error
	FirstOrCreate(context context.Context, dst interface{}, filters ...gpa.Filter) error
	Create(context context.Context, dst interface{}) error
	Updates(context context.Context, dst interface{}, value interface{}, filters ...gpa.Filter) error
	UpdateColumns(context context.Context, dst interface{}, value interface{}, filters ...gpa.Filter) error
	Delete(context context.Context, dst interface{}, filters ...gpa.Filter) error
	Raw(context context.Context, sql string, values ...interface{}) (tx *gorm.DB)
}

type repository[T schema.Tabler] struct {
	repository gpa.Repository
}

func New[T schema.Tabler](db *gorm.DB, entity T) Repository[T] {
	return repository[T]{
		repository: gpa.New(db, entity),
	}
}

func (r repository[T]) WithTrx(trxHandle *gorm.DB) Repository[T] {
	r.repository = r.repository.WithTrx(trxHandle)

	return r
}

func (r repository[T]) Count(context context.Context, filters ...gpa.Filter) (int, error) {
	return r.repository.Count(context, filters...)
}

func (r repository[T]) First(context context.Context, model interface{}, filters ...gpa.Filter) (err error) {
	return r.repository.First(context, model, filters...)
}

func (r repository[T]) Last(context context.Context, model interface{}, filters ...gpa.Filter) (err error) {
	return r.repository.Last(context, model, filters...)
}

func (r repository[T]) Find(context context.Context, models interface{}, filters ...gpa.Filter) error {
	return r.repository.Find(context, models, filters...)
}

func (r repository[T]) Transaction(handler func(tx *gorm.DB) error) (err error) {
	tx := r.repository.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			if err := tx.Rollback().Error; err != nil {
				log.Errorf("failed to rollback transaction: %v", err)
			}

			utils.StackTraceHandler(r)

			err = r.(error)
		}
	}()

	if err := handler(tx); err != nil {
		if err := tx.Rollback().Error; err != nil {
			log.Errorf("failed to rollback transaction: %v", err)
		}

		return err
	}

	if err := tx.Commit().Error; err != nil {
		log.Errorf("failed to commit transaction: %v", err)
		return err
	}

	return nil
}

func (r repository[T]) FirstOrCreate(context context.Context, model interface{}, filters ...gpa.Filter) error {
	return r.repository.FirstOrCreate(context, model, filters...)
}

func (r repository[T]) Create(context context.Context, model interface{}) error {
	return r.repository.Create(context, model)
}

func (r repository[T]) Updates(context context.Context, model interface{}, value interface{}, filters ...gpa.Filter) error {
	return r.repository.Updates(context, model, value, filters...)
}

func (r repository[T]) UpdateColumns(context context.Context, model interface{}, value interface{}, filters ...gpa.Filter) error {
	return r.repository.UpdateColumns(context, model, value, filters...)
}

func (r repository[T]) Delete(context context.Context, model interface{}, filters ...gpa.Filter) error {
	return r.repository.Delete(context, model, filters...)
}

func (r repository[T]) Raw(context context.Context, sql string, values ...interface{}) (tx *gorm.DB) {
	return r.repository.Raw(context, sql, values...)
}
