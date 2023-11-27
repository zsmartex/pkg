package mpa

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Tabler interface {
	TableName() string
}

type Repository struct {
	DB     *mongo.Database
	tabler Tabler
}

func New(db *mongo.Database, entity Tabler) Repository {
	return Repository{db, entity}
}

func (r Repository) Count(ctx context.Context, filters ...Filter) (int, error) {
	result, err := r.DB.Collection(r.tabler.TableName()).CountDocuments(ctx, ApplyFilters(filters...))
	return int(result), err
}

func (r Repository) Find(ctx context.Context, models interface{}, filters []Filter, opts ...*options.FindOptions) error {
	cursor, err := r.DB.Collection(r.tabler.TableName()).Find(ctx, ApplyFilters(filters...), opts...)
	if err != nil {
		return err
	}

	if err := cursor.All(ctx, models); err != nil {
		return err
	}

	return nil
}

func (r Repository) First(ctx context.Context, model interface{}, filters ...Filter) error {
	if err := r.DB.Collection(r.tabler.TableName()).FindOne(ctx, ApplyFilters(filters...)).Decode(model); err != nil {
		return err
	}

	return nil
}

func (r Repository) Last(ctx context.Context, model interface{}, filters ...Filter) error {
	if err := r.DB.Collection(r.tabler.TableName()).FindOne(
		ctx,
		ApplyFilters(filters...),
		options.FindOne().SetSort(bson.D{{Key: "_id", Value: -1}}),
	).Decode(model); err != nil {
		return err
	}

	return nil
}

func (r Repository) FirstOrCreate(ctx context.Context, model interface{}, create interface{}, filters ...Filter) error {
	err := r.DB.Collection(r.tabler.TableName()).FindOne(ctx, ApplyFilters(filters...)).Decode(&model)
	if err == nil {
		return nil
	}

	if _, err := r.DB.Collection(r.tabler.TableName()).InsertOne(ctx, create); err != nil {
		return err
	}

	model = create

	return nil
}

func (r Repository) Create(ctx context.Context, model interface{}, filters ...Filter) error {
	if _, err := r.DB.Collection(r.tabler.TableName()).InsertOne(ctx, model); err != nil {
		return err
	}

	return nil
}

func (r Repository) Updates(ctx context.Context, model interface{}, value map[string]interface{}, filters ...Filter) error {
	result, err := r.DB.Collection(r.tabler.TableName()).UpdateOne(ctx, model, map[string]interface{}{"$set": value})
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r Repository) Delete(ctx context.Context, filters ...Filter) error {
	_, err := r.DB.Collection(r.tabler.TableName()).DeleteOne(ctx, ApplyFilters(filters...))
	if err != nil {
		return err
	}

	return nil
}
