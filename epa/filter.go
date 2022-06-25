package epa

import (
	"github.com/olivere/elastic/v7"
)

type Ordering string

const (
	OrderingAscending  Ordering = "asc"
	OrderingDescending Ordering = "desc"
)

type Query struct {
	Page     int
	Limit    int
	OrderBy  string
	Ordering Ordering
	Filters  []Filter
}

type Filter func(query *elastic.BoolQuery) *elastic.BoolQuery

func ApplyFilters(query *elastic.BoolQuery, filters []Filter) *elastic.BoolQuery {
	for _, f := range filters {
		query = f(query)
	}
	return query
}

func ChainFilters(filters ...Filter) Filter {
	return func(query *elastic.BoolQuery) *elastic.BoolQuery {
		for _, f := range filters {
			query = f(query)
		}
		return query
	}
}
