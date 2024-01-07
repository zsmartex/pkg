//go:generate mockgen -source=repository.go -destination=mock/repository.go -package=mock_elasticsearch

package elasticsearch_fx

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"

	"github.com/zsmartex/pkg/v2/epa"
	"github.com/zsmartex/pkg/v2/epa/aggregation"
	"github.com/zsmartex/pkg/v2/epa/query"
	"github.com/zsmartex/pkg/v2/queries"
)

type Result[T any] struct {
	Values       []*T
	TotalHits    int
	Aggregations aggregation.Aggregations
}

type Schema interface {
	IndexName() string
}

type Repository[T any] interface {
	CheckHealth(ctx context.Context) bool
	Find(ctx context.Context, q epa.Query) (*Result[T], error)
}

type repository[T any] struct {
	*elasticsearch.Client
	Schema
}

func NewRepository[T Schema](client *elasticsearch.Client, entity T) Repository[T] {
	return repository[T]{
		client,
		entity,
	}
}

func (r repository[T]) CheckHealth(ctx context.Context) bool {
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

func (r repository[T]) Find(ctx context.Context, q epa.Query) (*Result[T], error) {
	var indexes []string
	if len(q.Indexes) > 0 {
		indexes = q.Indexes
	} else {
		indexes = []string{r.IndexName()}
	}

	searchRequest := make([]func(*esapi.SearchRequest), 0)
	searchRequest = append(searchRequest,
		r.Search.WithContext(ctx),
		r.Client.Search.WithIndex(indexes...),
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
		q, err := epa.ApplyFilters(query.NewBoolQuery(), q.Filters).Source()
		if err != nil {
			return nil, err
		}
		queryMap["query"] = q
	}

	if len(q.Aggregations) > 0 {
		aggs := make(map[string]interface{})

		for name, agg := range q.Aggregations {
			var err error
			aggs[name], err = agg.Source()
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

	resBuf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 404 {
		return &Result[T]{
			TotalHits: 0,
			Values:    []*T{},
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

	values := make([]*T, 0)

	for _, hit := range response.Hits.Hits {
		var value *T
		if err := json.Unmarshal(hit.Source, &value); err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	aggregations := make(aggregation.Aggregations)
	if len(q.Aggregations) > 0 {
		if err := json.Unmarshal(response.Aggregations, &aggregations); err != nil {
			return nil, err
		}
	}

	return &Result[T]{
		TotalHits:    response.Hits.Total.Value,
		Values:       values,
		Aggregations: aggregations,
	}, nil
}
