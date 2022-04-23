package repository

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Repository struct {
	*gorm.DB
	schema.Tabler
}

func New(db *gorm.DB, entity schema.Tabler) *Repository {
	return &Repository{db, entity}
}

func (r Repository) Count(ctx context.Context, filters []Filter) (int, error) {
	var result int64
	err := ApplyFilters(r.DB.WithContext(ctx).Table(r.TableName()), filters).Count(&result).Error
	return int(result), err
}

func (r Repository) Find(ctx context.Context, models interface{}, filters []Filter) error {
	return ApplyFilters(r.DB.WithContext(ctx).Table(r.TableName()), filters).Find(models).Error
}

func (r Repository) First(ctx context.Context, model interface{}, filters []Filter) error {
	return ApplyFilters(r.DB.WithContext(ctx).Table(r.TableName()), filters).First(model).Error
}

func (r Repository) Last(ctx context.Context, model interface{}, filters []Filter) error {
	return ApplyFilters(r.DB.WithContext(ctx).Table(r.TableName()), filters).Last(model).Error
}

func (r Repository) FirstOrCreate(ctx context.Context, model interface{}, opts []TransactionOption) error {
	opts = append(opts, WithFilters(
		ChainFilters(func(query *gorm.DB) *gorm.DB {
			return query.Assign(model)
		}),
	))

	return MakeTransactionWithActionNonValue(FirstOrCreate, model, opts)(r.DB.WithContext(ctx).Table(r.TableName()))
}

func (r Repository) Take(ctx context.Context, model interface{}, filters []Filter) error {
	return ApplyFilters(r.DB.WithContext(ctx).Table(r.TableName()), filters).Take(model).Error
}

func (r Repository) Create(ctx context.Context, model interface{}, opts []TransactionOption) error {
	return MakeTransactionWithActionNonValue(Create, model, opts)(r.DB.WithContext(ctx).Table(r.TableName()))
}

func (r Repository) Updates(ctx context.Context, model interface{}, value interface{}, opts []TransactionOption) error {
	return MakeTransactionWithActionValue(Updates, model, value, opts)(r.DB.WithContext(ctx).Table(r.TableName()))
}

func (r Repository) UpdateColumns(ctx context.Context, model interface{}, value interface{}, opts []TransactionOption) error {
	return MakeTransactionWithActionValue(UpdateColumns, model, value, opts)(r.DB.WithContext(ctx).Table(r.TableName()))
}

func (r Repository) Save(ctx context.Context, model interface{}, opts []TransactionOption) error {
	return MakeTransactionWithActionNonValue(Save, model, opts)(r.DB.WithContext(ctx).Table(r.TableName()))
}

func (r Repository) Delete(ctx context.Context, model interface{}, opts []TransactionOption) error {
	return MakeTransactionWithActionNonValue(Delete, model, opts)(r.DB.WithContext(ctx).Table(r.TableName()))
}
