package epa

import (
	"github.com/zsmartex/pkg/v2/epa/aggregation"
	"github.com/zsmartex/pkg/v2/epa/query"
	"github.com/zsmartex/pkg/v2/queries"
)

type Query struct {
	Indexes      []string
	Page         int
	Limit        int
	OrderBy      string
	Ordering     queries.Ordering
	Filters      []Filter
	SearchAfter  []interface{}
	Aggregations map[string]aggregation.Aggregation
}

type Filter func(*query.BoolQuery) *query.BoolQuery

func ApplyFilters(q *query.BoolQuery, filters []Filter) *query.BoolQuery {
	for _, f := range filters {
		q = f(q)
	}
	return q
}

func ChainFilters(filters ...Filter) Filter {
	return func(q *query.BoolQuery) *query.BoolQuery {
		for _, f := range filters {
			q = f(q)
		}
		return q
	}
}
