package mongo_fx

import (
	"context"

	"github.com/pkg/errors"
	"github.com/zsmartex/pkg/v2/mpa"
	"go.mongodb.org/mongo-driver/mongo"
)

type WriteRepository[T Tabler] interface {
	Transaction(ctx context.Context, handler func(sessionContext mongo.SessionContext) error) error
	FirstOrCreate(ctx context.Context, model *T, create *T, filters ...mpa.Filter) error
	Create(ctx context.Context, model *T, filters ...mpa.Filter) error
	Updates(ctx context.Context, model *T, value map[string]interface{}, filters ...mpa.Filter) error
	Delete(ctx context.Context, filters ...mpa.Filter) error
}

type writeRepository[T Tabler] struct {
	repository mpa.Repository
}

func NewWriteRepository[T Tabler](db *mongo.Database, entity T) WriteRepository[T] {
	return writeRepository[T]{
		repository: mpa.New(db, entity),
	}
}

func (r writeRepository[T]) Transaction(ctx context.Context, handler func(sessionContext mongo.SessionContext) error) error {
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

func (r writeRepository[T]) FirstOrCreate(ctx context.Context, model *T, create *T, filters ...mpa.Filter) error {
	err := r.repository.FirstOrCreate(ctx, model, &create, filters...)
	if err != nil {
		return errors.Wrap(err, "repository first or create")
	}

	return nil
}

func (r writeRepository[T]) Create(ctx context.Context, model *T, filters ...mpa.Filter) error {
	err := r.repository.Create(ctx, &model, filters...)
	if err != nil {
		return errors.Wrap(err, "repository create")
	}

	return nil
}

func (r writeRepository[T]) Updates(ctx context.Context, model *T, value map[string]interface{}, filters ...mpa.Filter) error {
	err := r.repository.Updates(ctx, &model, value, filters...)
	if err != nil {
		return errors.Wrap(err, "repository update")
	}

	return nil
}

func (r writeRepository[T]) Delete(ctx context.Context, filters ...mpa.Filter) error {
	err := r.repository.Delete(ctx, filters...)
	if err != nil {
		return errors.Wrap(err, "repository delete")
	}

	return nil
}
