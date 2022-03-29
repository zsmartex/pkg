package filters

import (
	"gorm.io/gorm"

	"github.com/zsmartex/pkg/repository"
)

func WithLimit(limit int) repository.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Limit(limit)
	}
}

func WithOrder(order string) repository.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Order(order)
	}
}

func WithOffset(offset int) repository.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Offset(offset)
	}
}

func WithPageable(page, size int, order string) repository.Filter {
	return func(query *gorm.DB) *gorm.DB {
		query = query.Limit(size)
		if page > 0 {
			offset := (page - 1) * size
			query = query.Offset(offset)
		}
		if order != "" {
			query = query.Order(order)
		}
		return query
	}
}
