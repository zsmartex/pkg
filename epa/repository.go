package epa

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/zsmartex/pkg/v2/epa/aggregation"
	"github.com/zsmartex/pkg/v2/epa/query"
	"github.com/zsmartex/pkg/v2/queries"
)

type Result[T any] struct {
	Values       []T
	TotalHits    int
	Aggregations aggregation.Aggregations
}

type Schema interface {
	IndexName() string
}

type Repository[T any] struct {
	*elasticsearch.Client
	Schema
}

func New[T Schema](client *elasticsearch.Client, entity T) Repository[T] {
	return Repository[T]{
		client,
		entity,
	}
}

func (r Repository[T]) CheckHealth(ctx context.Context) bool {
	_, err := r.Info(r.Info.WithContext(ctx))
	return err == nil
}

type Response struct {
	Hits struct {
		Hits []struct {
			Index  string          `json:"_index"`
			ID     string          `json:"_id"`
			Score  float64         `json:"_score"`
			Source json.RawMessage `json:"_source"`
		} `json:"hits"`
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
	} `json:"hits"`

	Aggregations json.RawMessage `json:"aggregations"`
}

type ErrorResponse struct {
	Error struct {
		RootCause []struct {
			Type         string `json:"type"`
			Reason       string `json:"reason"`
			ResourceType string `json:"resource.type"`
			ResourceID   string `json:"resource.id"`
			IndexUUID    string `json:"index_uuid"`
			Index        string `json:"index"`
		} `json:"root_cause"`
		Type         string `json:"type"`
		Reason       string `json:"reason"`
		ResourceType string `json:"resource.type"`
		ResourceID   string `json:"resource.id"`
		IndexUUID    string `json:"index_uuid"`
		Index        string `json:"index"`
	} `json:"error"`
	Status int `json:"status"`
}

func (r Repository[T]) Find(ctx context.Context, q Query) (*Result[T], error) {
	searchRequest := make([]func(*esapi.SearchRequest), 0)
	searchRequest = append(searchRequest,
		r.Search.WithContext(ctx),
		r.Client.Search.WithIndex(r.IndexName()),
	)

	if q.Page > 0 {
		searchRequest = append(
			searchRequest,
			r.Client.Search.WithFrom((q.Page-1)*q.Limit),
			r.Client.Search.WithTrackTotalHits(true),
		)
	}

	searchRequest = append(searchRequest, r.Client.Search.WithSize(q.Limit))

	queryMap := map[string]interface{}{}

	if len(q.Filters) > 0 {
		q, err := ApplyFilters(query.NewBoolQuery(), q.Filters).Source()
		if err != nil {
			return nil, err
		}
		queryMap["query"] = q
	}

	if len(q.Aggregations) > 0 {
		aggs := make(map[string]interface{})

		for name, aggregation := range q.Aggregations {
			var err error
			aggs[name], err = aggregation.Source()
			if err != nil {
				return nil, err
			}
		}

		queryMap["aggs"] = aggs
	}

	if q.OrderBy != "" {
		ordering := q.Ordering
		if len(ordering) == 0 {
			ordering = queries.OrderingAsc
		}

		queryMap["sort"] = []interface{}{
			map[string]interface{}{
				q.OrderBy: ordering,
			},
		}
	}

	data, err := json.Marshal(queryMap)
	if err != nil {
		return nil, err
	}

	searchRequest = append(searchRequest, r.Client.Search.WithBody(bytes.NewReader(data)))

	res, err := r.Client.Search(searchRequest...)
	if err != nil {
		return nil, err
	}

	resBuf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 404 {
		return &Result[T]{
			TotalHits: 0,
			Values:    []T{},
		}, nil
	}

	if res.IsError() {
		var errResponse ErrorResponse
		if err := json.Unmarshal(resBuf, &errResponse); err != nil {
			return nil, err
		}

		return nil, errors.New(errResponse.Error.RootCause[0].Reason)
	}

	var response Response
	if err := json.Unmarshal(resBuf, &response); err != nil {
		return nil, err
	}

	values := make([]T, 0)

	for _, hit := range response.Hits.Hits {
		var value T
		if err := json.Unmarshal(hit.Source, &value); err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	aggregations := make(aggregation.Aggregations)
	if err := json.Unmarshal(response.Aggregations, &aggregations); err != nil {
		return nil, err
	}

	return &Result[T]{
		TotalHits:    response.Hits.Total.Value,
		Values:       values,
		Aggregations: aggregations,
	}, nil
}
