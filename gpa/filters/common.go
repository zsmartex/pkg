package filters

import (
	"time"

	"gorm.io/gorm"

	"github.com/zsmartex/pkg/v2/gpa"
)

func WithID(id interface{}) gpa.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where("id=?", id)
	}
}

func WithIDs(ids ...interface{}) gpa.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where("id IN ?", ids)
	}
}

func WithCreatedAtBy(created_at time.Time) gpa.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where("created_by = ?", created_at)
	}
}

func WithUpdatedAtBy(updated_at time.Time) gpa.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where("updated_by = ?", updated_at)
	}
}

func WithCreatedAtAfter(t time.Time) gpa.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where("created_at > ?", t)
	}
}

func WithCreateAtBefore(t time.Time) gpa.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where("created_at < ?", t)
	}
}

func WithUpdatedAtAfter(t time.Time) gpa.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where("created_at > ?", t)
	}
}

func WithUpdateAtBefore(t time.Time) gpa.Filter {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where("created_at < ?", t)
	}
}
