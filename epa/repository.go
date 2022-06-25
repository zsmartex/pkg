package epa

import (
	"context"
	"reflect"

	"github.com/olivere/elastic/v7"
)

type Result[T any] struct {
	Values    []T
	TotalHits int64
}

type repository[T any] struct {
	client *elastic.Client
}

type Config struct {
	URL      string
	Username string
	Password string
}

func New[T any](cfg *Config) (*repository[T], error) {
	client, err := elastic.NewClient(elastic.SetBasicAuth(cfg.Username, cfg.Password), elastic.SetURL(cfg.URL), elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}

	return &repository[T]{client}, nil
}

func (r *repository[T]) Create(ctx context.Context, index, id string, body *T) (*elastic.IndexResponse, error) {
	return r.client.
		Index().
		Index(index).
		Id(id).
		BodyJson(body).
		Do(ctx)
}

func (r *repository[T]) Find(ctx context.Context, index string, query Query) (*Result[T], error) {
	search := r.client.
		Search().
		Index(index).
		Query(ApplyFilters(elastic.NewBoolQuery(), query.Filters))

	if query.Page > 0 {
		search = search.From((query.Page - 1) * query.Limit)
	}

	if query.Limit > 0 {
		search = search.Size(query.Limit)
	}

	if query.OrderBy != "" {
		search = search.Sort(query.OrderBy, query.Ordering == OrderingAscending)
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
		Values:    values,
		TotalHits: result.TotalHits(),
	}, nil
}
