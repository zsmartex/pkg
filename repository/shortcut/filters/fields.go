package filters

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/zsmartex/pkg/repository"
)

func WithFieldEqual(field string, value interface{}) repository.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(fmt.Sprintf("%s = ?", field), value)
	}
}

func WithFieldNotEqual(field string, value interface{}) repository.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(fmt.Sprintf("%s <> ?", field), value)
	}
}

func WithFieldGreaterThan(field string, value interface{}) repository.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(fmt.Sprintf("%s > ?", field), value)
	}
}

func WithFieldLessThan(field string, value interface{}) repository.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(fmt.Sprintf("%s < ?", field), value)
	}
}

func WithFieldGreaterThanOrEqualTo(field string, value interface{}) repository.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(fmt.Sprintf("%s >= ?", field), value)
	}
}

func WithFieldLessThanOrEqualTo(field string, value interface{}) repository.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(fmt.Sprintf("%s <= ?", field), value)
	}
}

func WithFieldIn(field string, value interface{}) repository.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(fmt.Sprintf("%s IN ?", field), value)
	}
}

func WithFieldNotIn(field string, value interface{}) repository.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(fmt.Sprintf("%s NOT IN ?", field), value)
	}
}

func WithFieldLike(field string, value interface{}) repository.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(fmt.Sprintf("%s LIKE ?", field), value)
	}
}

func WithNotLIKE(field string, value interface{}) repository.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(fmt.Sprintf("%s LIKE ?", field), value)
	}
}

func WithFieldIsNull(field string) repository.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(fmt.Sprintf("%s IS NULL", field))
	}
}
