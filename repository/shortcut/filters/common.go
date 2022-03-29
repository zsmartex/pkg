package filters

import (
	"time"

	"gorm.io/gorm"

	"github.com/zsmartex/pkg/repository"
)

func WithID(id interface{}) repository.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where("id=?", id)
	}
}

func WithIDs(ids ...interface{}) repository.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where("id IN ?", ids)
	}
}

func WithCreatedBy(user interface{}) repository.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where("created_by = ?", user)
	}
}

func WithUpdatedBy(user interface{}) repository.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where("updated_by = ?", user)
	}
}

func WithCreatedAtAfter(t time.Time) repository.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where("created_at > ?", t)
	}
}

func WithCreateAtBefore(t time.Time) repository.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where("created_at < ?", t)
	}
}

func WithUpdatedAtAfter(t time.Time) repository.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where("created_at > ?", t)
	}
}

func WithUpdateAtBefore(t time.Time) repository.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where("created_at < ?", t)
	}
}
