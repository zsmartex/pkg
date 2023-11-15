package gorm_fx

import (
	"context"

	"github.com/cockroachdb/errors"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/zsmartex/mergo"
	"github.com/zsmartex/pkg/v2/gpa"
	"github.com/zsmartex/pkg/v2/utils"
)

type Repository[T schema.Tabler] interface {
	DB() *gorm.DB
	TableName() string
	AddCallback(kind CallbackType, callback func(db *gorm.DB, value *T) error)
	WithTrx(trxHandle *gorm.DB) Repository[T]
	Transaction(handler func(tx *gorm.DB) error) error
	Count(ctx context.Context, filters ...gpa.Filter) (count int, err error)
	Last(ctx context.Context, filters ...gpa.Filter) (model *T, err error)
	First(ctx context.Context, filters ...gpa.Filter) (model *T, err error)
	Find(ctx context.Context, filters ...gpa.Filter) (models []*T, err error)
	FirstOrCreate(ctx context.Context, model *T, filters ...gpa.Filter) error
	Create(ctx context.Context, model *T, filters ...gpa.Filter) error
	Updates(ctx context.Context, model *T, value interface{}, filters ...gpa.Filter) error
	UpdateColumns(ctx context.Context, model *T, value interface{}, filters ...gpa.Filter) error
	Delete(ctx context.Context, model *T, filters ...gpa.Filter) error
	Exec(ctx context.Context, sql string, attrs ...interface{}) error
	RawFind(ctx context.Context, dst interface{}, sql string, attrs ...interface{}) error
	RawScan(ctx context.Context, dst interface{}, sql string, attrs ...interface{}) error
	RawFirst(ctx context.Context, dst interface{}, sql string, attrs ...interface{}) error
}

type repository[T schema.Tabler] struct {
	repository gpa.Repository
	entity     T
}

func NewRepository[T schema.Tabler](db *gorm.DB, entity T) Repository[T] {
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

func (r repository[T]) Transaction(handler func(tx *gorm.DB) error) error {
	err := r.repository.DB.Transaction(handler)

	return errors.Wrap(err, "repository transaction")
}

func (r repository[T]) AddCallback(kind CallbackType, callback func(db *gorm.DB, value *T) error) {
	if callbacks[kind][r.TableName()] != nil {
		return
	}

	callbacks[kind][r.TableName()] = func(db *gorm.DB) error {
		model, ok := db.Statement.Model.(*T)

		if model == nil {
			dest, ok := db.Statement.Dest.(*T)
			if !ok {
				return nil
			}

			if err := callback(db, dest); err != nil {
				return err
			}

			if err := validateModel(*dest); err != nil {
				return err
			}

			return nil
		}

		if !ok {
			return nil
		}

		var dest *T
		if _, ok := db.Statement.Dest.(*T); ok {
			dest = db.Statement.Dest.(*T)
		} else if _, ok := db.Statement.Dest.(T); ok {
			val := db.Statement.Dest.(T)
			dest = &val
		} else {
			return nil
		}

		destCopy := *dest

		mergo.Merge(&destCopy, model, mergo.WithOverwriteOnlyEmptyValue)

		if err := callback(db, &destCopy); err != nil {
			return err
		}

		if err := validateModel(destCopy); err != nil {
			return err
		}

		if err := utils.CompareDiff(dest, model, destCopy); err != nil {
			return err
		}

		db.Statement.Dest = dest

		return nil
	}
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

func (r repository[T]) First(context context.Context, filters ...gpa.Filter) (model *T, err error) {
	err = r.repository.First(context, &model, filters...)
	if err != nil {
		return nil, errors.Wrap(err, "repository first")
	}

	return model, nil
}

func (r repository[T]) Last(context context.Context, filters ...gpa.Filter) (model *T, err error) {
	err = r.repository.Last(context, &model, filters...)
	if err != nil {
		return nil, errors.Wrap(err, "repository last")
	}

	return model, nil
}

func (r repository[T]) Find(context context.Context, filters ...gpa.Filter) (models []*T, err error) {
	err = r.repository.Find(context, &models, filters...)
	if err != nil {
		return nil, errors.Wrap(err, "repository find")
	}

	return models, nil
}

func (r repository[T]) FirstOrCreate(context context.Context, model *T, filters ...gpa.Filter) error {
	err := r.repository.FirstOrCreate(context, model, filters...)

	return errors.Wrap(err, "repository first or create")
}

func (r repository[T]) Create(context context.Context, model *T, filters ...gpa.Filter) error {
	err := r.repository.Create(context, model, filters...)

	return errors.Wrap(err, "repository create")
}

func (r repository[T]) Updates(context context.Context, model *T, value interface{}, filters ...gpa.Filter) error {
	err := r.repository.Updates(context, model, value, filters...)

	return errors.Wrap(err, "repository update")
}

func (r repository[T]) UpdateColumns(context context.Context, model *T, value interface{}, filters ...gpa.Filter) error {
	err := r.repository.UpdateColumns(context, model, value, filters...)

	return errors.Wrap(err, "repository update columns")
}

func (r repository[T]) Delete(context context.Context, model *T, filters ...gpa.Filter) error {
	err := r.repository.Delete(context, model, filters...)

	return errors.Wrap(err, "repository delete")
}

func (r repository[T]) Raw(context context.Context, sql string, values ...interface{}) *gorm.DB {
	return r.repository.Raw(context, sql, values...)
}

func (u repository[T]) RawFind(context context.Context, dst interface{}, sql string, attrs ...interface{}) error {
	err := u.repository.Raw(context, sql, attrs...).Find(dst).Error

	return errors.Wrap(err, "usecase raw find")
}

func (u repository[T]) RawScan(context context.Context, dst interface{}, sql string, attrs ...interface{}) error {
	err := u.repository.Raw(context, sql, attrs...).Scan(dst).Error

	return errors.Wrap(err, "usecase raw scan")
}

func (u repository[T]) RawFirst(context context.Context, dst interface{}, sql string, attrs ...interface{}) error {
	err := u.repository.Raw(context, sql, attrs...).First(dst).Error

	return errors.Wrap(err, "usecase raw first")
}

func (r repository[T]) Exec(context context.Context, sql string, values ...interface{}) error {
	return r.repository.Exec(context, sql, values...).Error
}
