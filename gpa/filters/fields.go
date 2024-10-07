package filters

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/zsmartex/pkg/v2/gpa"
)

func WithFieldEqual(field string, value interface{}) gpa.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(fmt.Sprintf("%s = ?", field), value)
	}
}

func WithFieldNotEqual(field string, value interface{}) gpa.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(fmt.Sprintf("%s != ?", field), value)
	}
}

func WithFieldGreaterThan(field string, value interface{}) gpa.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(fmt.Sprintf("%s > ?", field), value)
	}
}

func WithFieldLessThan(field string, value interface{}) gpa.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(fmt.Sprintf("%s < ?", field), value)
	}
}

func WithFieldGreaterThanOrEqualTo(field string, value interface{}) gpa.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(fmt.Sprintf("%s >= ?", field), value)
	}
}

func WithFieldLessThanOrEqualTo(field string, value interface{}) gpa.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(fmt.Sprintf("%s <= ?", field), value)
	}
}

func WithFieldIn(field string, value interface{}) gpa.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(fmt.Sprintf("%s IN ?", field), value)
	}
}

func WithFieldNotIn(field string, value interface{}) gpa.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(fmt.Sprintf("%s NOT IN ?", field), value)
	}
}

func WithFieldLike(field string, value interface{}) gpa.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(fmt.Sprintf("%s LIKE ?", field), value)
	}
}

func WithNotLIKE(field string, value interface{}) gpa.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(fmt.Sprintf("%s NOT LIKE ?", field), value)
	}
}

func WithFieldIsNull(field string) gpa.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(fmt.Sprintf("%s IS NULL", field))
	}
}

func WithFieldIsNotNull(field string) gpa.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(fmt.Sprintf("%s IS NOT NULL", field))
	}
}
