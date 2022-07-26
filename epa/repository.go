package epa

import (
	"context"
	"reflect"

	"github.com/olivere/elastic/v7"
)

type Result[T any] struct {
	Values       []T
	TotalHits    int64
	Aggregations elastic.Aggregations
}

type Schema interface {
	IndexName() string
}

type Repository[T any] struct {
	*elastic.Client
	Schema
}

func New[T Schema](client *elastic.Client, entity T) Repository[T] {
	return Repository[T]{
		client,
		entity,
	}
}

func (r Repository[T]) CheckHealth(ctx context.Context) bool {
	_, err := r.ClusterHealth().Do(ctx)

	return err == nil
}

func (r Repository[T]) Create(ctx context.Context, id string, body *T) (*elastic.IndexResponse, error) {
	return r.
		Index().
		Index(r.IndexName()).
		Id(id).
		BodyJson(body).
		Do(ctx)
}

func (r Repository[T]) Find(ctx context.Context, query Query) (*Result[T], error) {
	search := r.
		Search().
		Index(r.IndexName()).
		Query(ApplyFilters(elastic.NewBoolQuery(), query.Filters))

	if query.Page > 0 {
		search = search.From((query.Page - 1) * query.Limit)
		search = search.TrackTotalHits(true)
	}

	search = search.Size(query.Limit)

	if query.OrderBy != "" {
		search = search.Sort(query.OrderBy, query.Ordering == OrderingAscending)
	}

	if query.Addons != nil {
		search = query.Addons(search)
	}

	result, err := search.Do(ctx)
	if err != nil {
		return nil, err
	}

	var value T
	values := make([]T, 0)

	for _, v := range result.Each(reflect.TypeOf(value)) {
		values = append(values, v.(T))
	}

	return &Result[T]{
		Values:       values,
		TotalHits:    result.TotalHits(),
		Aggregations: result.Aggregations,
	}, nil
}
