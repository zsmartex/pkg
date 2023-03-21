package filters

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"github.com/zsmartex/pkg/v3/gpa"
)

func WithSchema(schema, table schema.Tabler) gpa.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Table(fmt.Sprintf("%s.%s", schema, table))
	}
}

func WithPreload(tableName string, cond ...interface{}) gpa.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Preload(tableName, cond...)
	}
}

func WithPreloadAll() gpa.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Preload(clause.Associations)
	}
}

func WithJoin(modelName string, args ...interface{}) gpa.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Joins(modelName, args...)
	}
}

func WithAssign(attrs ...interface{}) gpa.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Assign(attrs...)
	}
}

func WithSelect(first_column interface{}, columns ...interface{}) gpa.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Select(first_column, columns...)
	}
}

func WithOmit(fields ...string) gpa.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Omit(fields...)
	}
}

func WithGroup(name string) gpa.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Group(name)
	}
}

func WithLock() gpa.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Clauses(clause.Locking{
			Strength: "UPDATE",
		})
	}
}
