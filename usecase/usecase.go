package usecase

import (
	"context"
	"errors"

	"github.com/creasty/defaults"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gookit/validate"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zsmartex/mergo"
	"github.com/zsmartex/pkg/v2/epa"
	"github.com/zsmartex/pkg/v2/gpa"
	"github.com/zsmartex/pkg/v2/gpa/filters"
	"github.com/zsmartex/pkg/v2/repository"
	"github.com/zsmartex/pkg/v2/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	ErrBadConnection = errors.New("driver: bad connection")
)

type CallbackType string

const (
	CallbackTypeBeforeCreate CallbackType = "BeforeCreate"
	CallbackTypeAfterCreate  CallbackType = "AfterCreate"
	CallbackTypeBeforeSave   CallbackType = "BeforeSave"
	CallbackTypeAfterSave    CallbackType = "AfterSave"
)

var callbackReady = false
var callbacks = make(map[CallbackType]map[string]func(*gorm.DB) error, 0)

var _ IUsecase[schema.Tabler] = (*Usecase[schema.Tabler])(nil)

type Usecase[V schema.Tabler] struct {
	Repository           repository.Repository[V]
	ElasticsearchUsecase ElasticsearchUsecase[V]
	QuestDBUsecase       QuestDBUsecase[V]
	Omits                []string
}

type IUsecase[V schema.Tabler] interface {
	Count(ctx context.Context, filters ...gpa.Filter) (count int, err error)
	Last(ctx context.Context, filters ...gpa.Filter) (model *V, err error)
	First(ctx context.Context, filters ...gpa.Filter) (model *V, err error)
	Find(ctx context.Context, filters ...gpa.Filter) (models []*V, err error)
	Transaction(handler func(tx *gorm.DB) error) error
	FirstOrCreate(ctx context.Context, model *V, filters ...gpa.Filter) error
	Create(ctx context.Context, model *V, filters ...gpa.Filter) error
	Updates(ctx context.Context, model *V, value interface{}, filters ...gpa.Filter) error
	UpdateColumns(ctx context.Context, model *V, value interface{}, filters ...gpa.Filter) error
	Delete(ctx context.Context, model *V, filters ...gpa.Filter) error
	Exec(ctx context.Context, sql string, attrs ...interface{}) error
	RawFind(ctx context.Context, dst interface{}, sql string, attrs ...interface{}) error
	RawScan(ctx context.Context, dst interface{}, sql string, attrs ...interface{}) error
	RawFirst(ctx context.Context, dst interface{}, sql string, attrs ...interface{}) error

	Es() ElasticsearchUsecase[V]
	QuestDB() QuestDBUsecase[V]
}

func validateModel(model any) error {
	v := validate.Struct(model)

	if !v.Validate() {
		return v.Errors.OneError()
	}

	return nil
}

func InitCallback(db *gorm.DB) {
	callbacks[CallbackTypeBeforeCreate] = make(map[string]func(*gorm.DB) error, 0)
	callbacks[CallbackTypeAfterCreate] = make(map[string]func(*gorm.DB) error, 0)
	callbacks[CallbackTypeBeforeSave] = make(map[string]func(*gorm.DB) error, 0)
	callbacks[CallbackTypeAfterSave] = make(map[string]func(*gorm.DB) error, 0)

	db.Callback().Create().Before("gorm:create").Register("callback:before_create", func(db *gorm.DB) {
		defaults.Set(db.Statement.Dest)
		tableName := db.Statement.Table

		if callback, ok := callbacks[CallbackTypeBeforeCreate][tableName]; ok {
			db.AddError(callback(db))
		}

		if callback, ok := callbacks[CallbackTypeBeforeSave][tableName]; ok {
			db.AddError(callback(db))
		}
	})

	db.Callback().Create().After("gorm:create").Register("callback:after_create", func(db *gorm.DB) {
		tableName := db.Statement.Table

		if callback, ok := callbacks[CallbackTypeAfterCreate][tableName]; ok {
			db.AddError(callback(db))
		}

		if callback, ok := callbacks[CallbackTypeAfterSave][tableName]; ok {
			db.AddError(callback(db))
		}
	})

	db.Callback().Update().Before("gorm:update").Register("callback:before_update", func(db *gorm.DB) {
		tableName := db.Statement.Table

		if callback, ok := callbacks[CallbackTypeBeforeSave][tableName]; ok {
			db.AddError(callback(db))
		}
	})

	db.Callback().Update().After("gorm:update").Register("callback:after_update", func(db *gorm.DB) {
		tableName := db.Statement.Table

		callback, ok := callbacks[CallbackTypeAfterSave][tableName]
		if ok {
			db.AddError(callback(db))
		}
	})

	db.Callback().Delete().Before("gorm:delete").Register("callback:before_delete", func(db *gorm.DB) {
		tableName := db.Statement.Table

		if callback, ok := callbacks[CallbackTypeBeforeSave][tableName]; ok {
			db.AddError(callback(db))
		}
	})

	db.Callback().Delete().After("gorm:delete").Register("callback:after_delete", func(db *gorm.DB) {
		tableName := db.Statement.Table

		if callback, ok := callbacks[CallbackTypeAfterSave][tableName]; ok {
			db.AddError(callback(db))
		}
	})

	callbackReady = true
}

func (u Usecase[V]) AddCallback(kind CallbackType, callback func(db *gorm.DB, value *V) error) {
	if !callbackReady {
		InitCallback(u.Repository.DB())
	}

	if callbacks[kind][u.Repository.TableName()] != nil {
		return
	}

	callbacks[kind][u.Repository.TableName()] = func(db *gorm.DB) error {
		model, ok := db.Statement.Model.(*V)

		if model == nil {
			dest, ok := db.Statement.Dest.(*V)
			if !ok {
				return nil
			}

			if err := callback(db, dest); err != nil {
				return err
			}

			if err := validateModel(dest); err != nil {
				panic(err)
			}

			return nil
		}

		if !ok {
			return nil
		}

		var dest *V
		if _, ok := db.Statement.Dest.(*V); ok {
			dest = db.Statement.Dest.(*V)
		} else if _, ok := db.Statement.Dest.(V); ok {
			val := db.Statement.Dest.(V)
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

func (u Usecase[V]) Count(context context.Context, filters ...gpa.Filter) (count int, err error) {
	return u.Repository.Count(context, filters...)
}

func (u Usecase[V]) Last(context context.Context, filters ...gpa.Filter) (model *V, err error) {
	if err := u.Repository.Last(context, &model, filters...); err != nil {
		return nil, err
	}

	return
}

func (u Usecase[V]) First(context context.Context, filters ...gpa.Filter) (model *V, err error) {
	if err := u.Repository.First(context, &model, filters...); err != nil {
		return nil, err
	}

	return
}

func (u Usecase[V]) Find(context context.Context, filters ...gpa.Filter) (models []*V, err error) {
	if err := u.Repository.Find(context, &models, filters...); err != nil {
		return nil, err
	}

	return
}

func (u Usecase[V]) Transaction(handler func(tx *gorm.DB) error) error {
	return u.Repository.Transaction(handler)
}

func (u Usecase[V]) FirstOrCreate(context context.Context, model *V, filters ...gpa.Filter) error {
	return u.Repository.FirstOrCreate(context, model, filters...)
}

func (u Usecase[V]) Create(context context.Context, model *V, fs ...gpa.Filter) error {
	fs = append(fs, filters.WithOmit(u.Omits...))

	return u.Repository.Create(context, model, fs...)
}

func (u Usecase[V]) Updates(context context.Context, model *V, value interface{}, fs ...gpa.Filter) error {
	fs = append(fs, filters.WithOmit(u.Omits...))

	return u.Repository.Updates(context, model, value, fs...)
}

func (u Usecase[V]) UpdateColumns(context context.Context, model *V, value interface{}, fs ...gpa.Filter) error {
	fs = append(fs, filters.WithOmit(u.Omits...))

	return u.Repository.UpdateColumns(context, model, value, fs...)
}

func (u Usecase[V]) Delete(context context.Context, model *V, fs ...gpa.Filter) error {
	fs = append(fs, filters.WithOmit(u.Omits...))

	return u.Repository.Delete(context, model, fs...)
}

func (u Usecase[V]) Exec(context context.Context, sql string, attrs ...interface{}) error {
	return u.Repository.Exec(context, sql, attrs...).Error
}

func (u Usecase[V]) RawFind(context context.Context, dst interface{}, sql string, attrs ...interface{}) error {
	return u.Repository.Raw(context, sql, attrs...).Find(dst).Error
}

func (u Usecase[V]) RawScan(context context.Context, dst interface{}, sql string, attrs ...interface{}) error {
	return u.Repository.Raw(context, sql, attrs...).Scan(dst).Error
}

func (u Usecase[V]) RawFirst(context context.Context, dst interface{}, sql string, attrs ...interface{}) error {
	return u.Repository.Raw(context, sql, attrs...).First(dst).Error
}

func (u Usecase[V]) Es() ElasticsearchUsecase[V] {
	return u.ElasticsearchUsecase
}

func (u Usecase[V]) QuestDB() QuestDBUsecase[V] {
	return u.QuestDBUsecase
}

type ElasticsearchUsecase[T schema.Tabler] struct {
	Repository epa.Repository[T]
}

func (u ElasticsearchUsecase[T]) Find(context context.Context, query epa.Query) (*epa.Result[T], error) {
	return u.Repository.Find(context, query)
}

type QuestDBUsecase[V schema.Tabler] struct {
	Conn *pgxpool.Pool
}

func (u QuestDBUsecase[V]) Exec(ctx context.Context, sql string, attrs ...interface{}) error {
	_, err := u.Conn.Exec(ctx, sql, attrs...)
	if err != nil {
		return err
	}

	return nil
}

func (u QuestDBUsecase[V]) Query(ctx context.Context, dst interface{}, sql string, attrs ...interface{}) error {
	return pgxscan.Select(ctx, u.Conn, dst, sql, attrs...)
}

func (u QuestDBUsecase[V]) QueryRow(ctx context.Context, dst interface{}, sql string, attrs ...interface{}) error {
	return pgxscan.Get(ctx, u.Conn, dst, sql, attrs...)
}
