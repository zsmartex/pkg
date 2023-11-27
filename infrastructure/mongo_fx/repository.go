package mongo_fx

import (
	"context"

	"github.com/pkg/errors"
	"github.com/zsmartex/pkg/v2/mpa"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Tabler interface {
	TableName() string
}

type Repository[T Tabler] interface {
	DB() *mongo.Database
	TableName() string
	Transaction(ctx context.Context, handler func(sessionContext mongo.SessionContext) error) error
	Count(ctx context.Context, filters ...mpa.Filter) (int, error)
	Find(ctx context.Context, filters []mpa.Filter, opts ...*options.FindOptions) (models []*T, err error)
	First(ctx context.Context, filters ...mpa.Filter) (model *T, err error)
	Last(ctx context.Context, filters ...mpa.Filter) (model *T, err error)
	FirstOrCreate(ctx context.Context, model *T, create *T, filters ...mpa.Filter) error
	Create(ctx context.Context, model *T, filters ...mpa.Filter) error
	Updates(ctx context.Context, model *T, value map[string]interface{}, filters ...mpa.Filter) error
	Delete(ctx context.Context, filters ...mpa.Filter) error
}

type repository[T Tabler] struct {
	repository mpa.Repository
	entity     T
}

func NewRepository[T Tabler](db *mongo.Database, entity T) Repository[T] {
	return repository[T]{
		repository: mpa.New(db, entity),
	}
}

func (r repository[T]) DB() *mongo.Database {
	return r.repository.DB
}

func (r repository[T]) TableName() string {
	return r.entity.TableName()
}

func (r repository[T]) Transaction(ctx context.Context, handler func(sessionContext mongo.SessionContext) error) error {
	session, err := r.repository.DB.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	if err := mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}

		if err := handler(sessionContext); err != nil {
			return err
		}

		if err := session.CommitTransaction(sessionContext); err != nil {
			return err
		}

		return nil
	}); err != nil {
		if abortErr := session.AbortTransaction(context.Background()); abortErr != nil {
			return abortErr
		}

		return err
	}

	return nil
}

func (r repository[T]) Count(ctx context.Context, filters ...mpa.Filter) (int, error) {
	count, err := r.repository.Count(ctx, filters...)
	if err != nil {
		return 0, errors.Wrap(err, "repository count")
	}

	return count, err
}

func (r repository[T]) Find(ctx context.Context, filters []mpa.Filter, opts ...*options.FindOptions) (models []*T, err error) {
	err = r.repository.Find(ctx, &models, filters, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "repository find")
	}

	return models, nil
}

func (r repository[T]) First(ctx context.Context, filters ...mpa.Filter) (model *T, err error) {
	err = r.repository.First(ctx, &model, filters...)
	if err != nil {
		return nil, errors.Wrap(err, "repository first")
	}

	return model, nil
}

func (r repository[T]) Last(ctx context.Context, filters ...mpa.Filter) (model *T, err error) {
	err = r.repository.Last(ctx, &model, filters...)
	if err != nil {
		return nil, errors.Wrap(err, "repository last")
	}

	return model, nil
}

func (r repository[T]) FirstOrCreate(ctx context.Context, model *T, create *T, filters ...mpa.Filter) error {
	err := r.repository.FirstOrCreate(ctx, &model, &create, filters...)
	if err != nil {
		return errors.Wrap(err, "repository first or create")
	}

	return nil
}

func (r repository[T]) Create(ctx context.Context, model *T, filters ...mpa.Filter) error {
	err := r.repository.Create(ctx, &model, filters...)
	if err != nil {
		return errors.Wrap(err, "repository create")
	}

	return nil
}

func (r repository[T]) Updates(ctx context.Context, model *T, value map[string]interface{}, filters ...mpa.Filter) error {
	err := r.repository.Updates(ctx, &model, value, filters...)
	if err != nil {
		return errors.Wrap(err, "repository update")
	}

	return nil
}

func (r repository[T]) Delete(ctx context.Context, filters ...mpa.Filter) error {
	err := r.repository.Delete(ctx, filters...)
	if err != nil {
		return errors.Wrap(err, "repository delete")
	}

	return nil
}
