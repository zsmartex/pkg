package usecase

import (
	"context"
	"errors"

	"github.com/creasty/defaults"
	"github.com/gookit/validate"
	"github.com/imdario/mergo"
	"github.com/zsmartex/pkg/v2/epa"
	"github.com/zsmartex/pkg/v2/gpa"
	"github.com/zsmartex/pkg/v2/gpa/filters"
	"github.com/zsmartex/pkg/v2/repository"
	"github.com/zsmartex/pkg/v2/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type CallbackType string

const (
	CallbackTypeBeforeCreate CallbackType = "BeforeCreate"
	CallbackTypeAfterCreate  CallbackType = "AfterCreate"
	CallbackTypeBeforeSave   CallbackType = "BeforeSave"
	CallbackTypeAfterSave    CallbackType = "AfterSave"
)

var callbacks = make(map[CallbackType]map[string]func(*gorm.DB) error, 0)

type Usecase[V schema.Tabler] struct {
	repository           repository.Repository[V]
	elasticsearchUsecase ElasticsearchUsecase[V]
	questDBUsecase       QuestDBUsecase[V]
	omits                []string
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
}

func (u Usecase[V]) AddCallback(kind CallbackType, callback func(db *gorm.DB, value *V) error) {
	if callbacks[kind][u.repository.TableName()] != nil {
		return
	}

	callbacks[kind][u.repository.TableName()] = func(db *gorm.DB) error {
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

		mergo.Merge(&destCopy, model)

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

func (u Usecase[V]) Count(context context.Context, filters ...gpa.Filter) int {
	if count, err := u.repository.Count(context, filters...); err != nil {
		panic(err)
	} else {
		return count
	}
}

func (u Usecase[V]) Last(context context.Context, filters ...gpa.Filter) (model *V, err error) {
	if err := u.repository.Last(context, &model, filters...); errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	} else if err != nil {
		panic(err)
	}

	return
}

func (u Usecase[V]) First(context context.Context, filters ...gpa.Filter) (model *V, err error) {
	if err := u.repository.First(context, &model, filters...); errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	} else if err != nil {
		panic(err)
	}

	return
}

func (u Usecase[V]) Find(context context.Context, filters ...gpa.Filter) (models []*V) {
	if err := u.repository.Find(context, &models, filters...); err != nil {
		panic(err)
	}

	return
}

func (u Usecase[V]) Transaction(handler func(tx *gorm.DB) error) error {
	return u.repository.Transaction(handler)
}

func (u Usecase[V]) FirstOrCreate(context context.Context, model *V, filters ...gpa.Filter) {
	if err := u.repository.FirstOrCreate(context, model, filters...); err != nil {
		panic(err)
	}
}

func (u Usecase[V]) Create(context context.Context, model *V, fs ...gpa.Filter) {
	fs = append(fs, filters.WithOmit(u.omits...))

	if err := u.repository.Create(context, model, fs...); err != nil {
		panic(err)
	}
}

func (u Usecase[V]) Updates(context context.Context, model *V, value interface{}, fs ...gpa.Filter) {
	fs = append(fs, filters.WithOmit(u.omits...))

	if err := u.repository.Updates(context, model, value, fs...); err != nil {
		panic(err)
	}
}

func (u Usecase[V]) UpdateColumns(context context.Context, model *V, value interface{}, fs ...gpa.Filter) {
	fs = append(fs, filters.WithOmit(u.omits...))

	if err := u.repository.UpdateColumns(context, model, value, fs...); err != nil {
		panic(err)
	}
}

func (u Usecase[V]) Delete(context context.Context, model *V, fs ...gpa.Filter) {
	fs = append(fs, filters.WithOmit(u.omits...))

	if err := u.repository.Delete(context, model, fs...); err != nil {
		panic(err)
	}
}

func (u Usecase[V]) Es() ElasticsearchUsecase[V] {
	return u.elasticsearchUsecase
}

func (u Usecase[V]) QuestDB() QuestDBUsecase[V] {
	return u.questDBUsecase
}

type ElasticsearchUsecase[T schema.Tabler] struct {
	repository epa.Repository[T]
}

func (u ElasticsearchUsecase[T]) Find(context context.Context, query epa.Query) *epa.Result[T] {
	result, err := u.repository.Find(context, query)
	if err != nil {
		panic(err)
	}

	return result
}

type QuestDBUsecase[V schema.Tabler] struct {
	repository repository.Repository[V]
}

func (u QuestDBUsecase[V]) Count(context context.Context, filters ...gpa.Filter) int {
	if count, err := u.repository.Count(context, filters...); err != nil {
		panic(err)
	} else {
		return count
	}
}

func (u QuestDBUsecase[V]) First(context context.Context, filters ...gpa.Filter) (model *V, err error) {
	if err := u.repository.First(context, &model, filters...); errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	} else if err != nil {
		panic(err)
	}

	return
}

func (u QuestDBUsecase[V]) Find(context context.Context, filters ...gpa.Filter) (models []*V) {
	if err := u.repository.Find(context, &models, filters...); err != nil {
		panic(err)
	}

	return
}

func (u QuestDBUsecase[V]) RawFind(context context.Context, dst interface{}, sql string, attrs ...interface{}) error {
	return u.repository.Raw(context, sql, attrs...).Find(dst).Error
}

func (u QuestDBUsecase[V]) RawScan(context context.Context, dst interface{}, sql string, attrs ...interface{}) error {
	return u.repository.Raw(context, sql, attrs...).Scan(dst).Error
}

func (u QuestDBUsecase[V]) RawFirst(context context.Context, dst interface{}, sql string, attrs ...interface{}) error {
	return u.repository.Raw(context, sql, attrs...).First(dst).Error
}
