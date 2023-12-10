package gorm_fx

import (
	"github.com/creasty/defaults"
	"github.com/gookit/validate"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/zsmartex/pkg/v2"
)

var (
	Module = fx.Module("gorm_fx", gormProviders, gormInvokes)

	gormProviders = fx.Provide(New)

	gormInvokes = fx.Invoke(registerHooks)
)

func registerHooks(db *gorm.DB) {
	InitCallback(db)
}

type CallbackType string

const (
	CallbackTypeBeforeCreate CallbackType = "BeforeCreate"
	CallbackTypeAfterCreate  CallbackType = "AfterCreate"
	CallbackTypeBeforeSave   CallbackType = "BeforeSave"
	CallbackTypeAfterSave    CallbackType = "AfterSave"
)

var callbacks = map[CallbackType]map[string]func(*gorm.DB) error{
	CallbackTypeBeforeCreate: make(map[string]func(*gorm.DB) error),
	CallbackTypeAfterCreate:  make(map[string]func(*gorm.DB) error),
	CallbackTypeBeforeSave:   make(map[string]func(*gorm.DB) error),
	CallbackTypeAfterSave:    make(map[string]func(*gorm.DB) error),
}

func validateModel[T schema.Tabler](model T) error {
	v := validate.Struct(model)

	if !v.Validate() {
		return pkg.NewError(422, v.Errors.One(), "model validate failed")
	}

	return nil
}

func InitCallback(db *gorm.DB) {
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
