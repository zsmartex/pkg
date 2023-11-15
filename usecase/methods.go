package usecase

import (
	"context"

	"github.com/zsmartex/pkg/v2/gpa"
	"github.com/zsmartex/pkg/v2/infrastructure/elasticsearch_fx"
	"github.com/zsmartex/pkg/v2/infrastructure/gorm_fx"
	"github.com/zsmartex/pkg/v2/infrastructure/questdb_fx"
	"gorm.io/gorm"
)

func (u Usecase[V]) AddCallback(kind gorm_fx.CallbackType, callback func(db *gorm.DB, value *V) error) {
	u.DatabaseRepo.AddCallback(kind, callback)
}

func (u Usecase[V]) Repository() gorm_fx.Repository[V] {
	return u.DatabaseRepo
}

func (u Usecase[V]) Count(ctx context.Context, filters ...gpa.Filter) (count int, err error) {
	return u.DatabaseRepo.Count(ctx, filters...)
}

func (u Usecase[V]) Last(ctx context.Context, filters ...gpa.Filter) (model *V, err error) {
	return u.DatabaseRepo.Last(ctx, filters...)
}

func (u Usecase[V]) First(ctx context.Context, filters ...gpa.Filter) (model *V, err error) {
	return u.DatabaseRepo.First(ctx, filters...)
}

func (u Usecase[V]) Find(ctx context.Context, filters ...gpa.Filter) (models []*V, err error) {
	return u.DatabaseRepo.Find(ctx, filters...)
}

func (u Usecase[V]) Transaction(handler func(tx *gorm.DB) error) error {
	return u.DatabaseRepo.Transaction(handler)
}

func (u Usecase[V]) FirstOrCreate(ctx context.Context, model *V, filters ...gpa.Filter) error {
	return u.DatabaseRepo.FirstOrCreate(ctx, model, filters...)
}

func (u Usecase[V]) Create(ctx context.Context, model *V, filters ...gpa.Filter) error {
	return u.DatabaseRepo.Create(ctx, model, filters...)
}

func (u Usecase[V]) Updates(ctx context.Context, model *V, value interface{}, filters ...gpa.Filter) error {
	return u.DatabaseRepo.Updates(ctx, model, value, filters...)
}

func (u Usecase[V]) UpdateColumns(ctx context.Context, model *V, value interface{}, filters ...gpa.Filter) error {
	return u.DatabaseRepo.UpdateColumns(ctx, model, value, filters...)
}

func (u Usecase[V]) Delete(ctx context.Context, model *V, filters ...gpa.Filter) error {
	return u.DatabaseRepo.Delete(ctx, model, filters...)
}

func (u Usecase[V]) Exec(ctx context.Context, sql string, attrs ...interface{}) error {
	return u.DatabaseRepo.Exec(ctx, sql, attrs...)
}

func (u Usecase[V]) RawFind(ctx context.Context, dst interface{}, sql string, attrs ...interface{}) error {
	return u.DatabaseRepo.RawFind(ctx, dst, sql, attrs...)
}

func (u Usecase[V]) RawScan(ctx context.Context, dst interface{}, sql string, attrs ...interface{}) error {
	return u.DatabaseRepo.RawScan(ctx, dst, sql, attrs...)
}

func (u Usecase[V]) RawFirst(ctx context.Context, dst interface{}, sql string, attrs ...interface{}) error {
	return u.DatabaseRepo.RawFirst(ctx, dst, sql, attrs...)
}

func (u Usecase[V]) Es() elasticsearch_fx.Repository[V] {
	return u.ElasticsearchRepo
}

func (u Usecase[V]) QuestDB() questdb_fx.Repository[V] {
	return u.QuestDBRepo
}
