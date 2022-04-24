package filters

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"github.com/zsmartex/pkg/repository"
)

func WithSchema(schema, table schema.Tabler) repository.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Table(fmt.Sprintf("%s.%s", schema, table))
	}
}

func WithPreload(tableName string, cond ...interface{}) repository.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Preload(tableName, cond...)
	}
}

func WithPreloadAll() repository.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Preload(clause.Associations)
	}
}

func WithJoin(modelName string, args ...interface{}) repository.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Joins(modelName, args...)
	}
}

func WithAssign(attrs ...interface{}) repository.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Assign(attrs...)
	}
}

func WithSelect(query string, args ...interface{}) repository.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Select(query, args...)
	}
}

func WithOmit(fields ...string) repository.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Omit(fields...)
	}
}

func WithGroup(name string) repository.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Group(name)
	}
}

func WithPluck(column string, dest interface{}) repository.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Pluck(column, dest)
	}
}
