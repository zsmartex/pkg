package mongo_fx

import (
	"context"

	"github.com/pkg/errors"
	"github.com/zsmartex/pkg/v2/mpa"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReadRepository[T Tabler] interface {
	DB() *mongo.Database
	TableName() string
	Count(ctx context.Context, filters ...mpa.Filter) (int, error)
	Find(ctx context.Context, opts *options.FindOptions, filters ...mpa.Filter) (models []*T, err error)
	First(ctx context.Context, filters ...mpa.Filter) (model *T, err error)
	Last(ctx context.Context, filters ...mpa.Filter) (model *T, err error)
}

type readRepository[T Tabler] struct {
	repository mpa.Repository
	entity     T
}

func NewRepository[T Tabler](db *mongo.Database, entity T) ReadRepository[T] {
	return readRepository[T]{
		repository: mpa.New(db, entity),
	}
}

func (r readRepository[T]) DB() *mongo.Database {
	return r.repository.DB
}

func (r readRepository[T]) TableName() string {
	return r.entity.TableName()
}

func (r readRepository[T]) Count(ctx context.Context, filters ...mpa.Filter) (int, error) {
	count, err := r.repository.Count(ctx, filters...)
	if err != nil {
		return 0, errors.Wrap(err, "repository count")
	}

	return count, err
}

func (r readRepository[T]) Find(ctx context.Context, opts *options.FindOptions, filters ...mpa.Filter) (models []*T, err error) {
	err = r.repository.Find(ctx, &models, opts, filters...)
	if err != nil {
		return nil, errors.Wrap(err, "repository find")
	}

	return models, nil
}

func (r readRepository[T]) First(ctx context.Context, filters ...mpa.Filter) (model *T, err error) {
	err = r.repository.First(ctx, &model, filters...)
	if err != nil {
		return nil, errors.Wrap(err, "repository first")
	}

	return model, nil
}

func (r readRepository[T]) Last(ctx context.Context, filters ...mpa.Filter) (model *T, err error) {
	err = r.repository.Last(ctx, &model, filters...)
	if err != nil {
		return nil, errors.Wrap(err, "repository last")
	}

	return model, nil
}
