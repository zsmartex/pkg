package epa

import (
	"time"

	"github.com/zsmartex/pkg/v3/epa/query"
)

func WithMultiFieldsNotEqual(value interface{}, fields ...string) Filter {
	return func(q *query.BoolQuery) *query.BoolQuery {
		return q.MustNot(query.NewMultiMatchQuery(value, fields...))
	}
}

func WithMultiFieldsEqual(value interface{}, fields ...string) Filter {
	return func(q *query.BoolQuery) *query.BoolQuery {
		return q.Must(query.NewMultiMatchQuery(value, fields...))
	}
}

func WithFieldEqual(field string, value interface{}) Filter {
	return func(q *query.BoolQuery) *query.BoolQuery {
		return q.Must(query.NewMatchQuery(field, value))
	}
}

func WithFieldExtractEqual(field string, value string) Filter {
	return func(q *query.BoolQuery) *query.BoolQuery {
		return q.Must(query.NewQueryStringQuery(value, field))
	}
}

func WithFieldNotEqual(field string, value interface{}) Filter {
	return func(q *query.BoolQuery) *query.BoolQuery {
		return q.MustNot(query.NewMatchQuery(field, value))
	}
}

func WithFieldGreaterThan(field string, value interface{}) Filter {
	return func(q *query.BoolQuery) *query.BoolQuery {
		return q.Must(query.NewRangeQuery(field).Gt(value))
	}
}

func WithFieldLessThan(field string, value interface{}) Filter {
	return func(q *query.BoolQuery) *query.BoolQuery {
		return q.Must(query.NewRangeQuery(field).Lt(value))
	}
}

func WithFieldGreaterThanOrEqualTo(field string, value interface{}) Filter {
	return func(q *query.BoolQuery) *query.BoolQuery {
		return q.Must(query.NewRangeQuery(field).Gte(value))
	}
}

func WithFieldLessThanOrEqualTo(field string, value interface{}) Filter {
	return func(q *query.BoolQuery) *query.BoolQuery {
		return q.Must(query.NewRangeQuery(field).Lte(value))
	}
}

func WithFieldIn(field string, values ...interface{}) Filter {
	return func(q *query.BoolQuery) *query.BoolQuery {
		return q.Must(query.NewTermsQuery(field, values...))
	}
}

func WithFieldNotIn(field string, values ...interface{}) Filter {
	return func(q *query.BoolQuery) *query.BoolQuery {
		return q.MustNot(query.NewTermsQuery(field, values...))
	}
}

func WithFieldLike(field string, value interface{}) Filter {
	return func(q *query.BoolQuery) *query.BoolQuery {
		return q.Should(query.NewMatchQuery(field, value))
	}
}

func WithNotLIKE(field string, value interface{}) Filter {
	return func(q *query.BoolQuery) *query.BoolQuery {
		return q.MustNot(query.NewMatchQuery(field, value))
	}
}

func WithFieldIsNull(field string) Filter {
	return func(q *query.BoolQuery) *query.BoolQuery {
		return q.Must(query.NewExistsQuery(field))
	}
}

func WithID(id string) Filter {
	return func(q *query.BoolQuery) *query.BoolQuery {
		return q.Must(query.NewMatchQuery("_id", id))
	}
}

func WithDateRange(field string, from interface{}, to interface{}) Filter {
	return func(q *query.BoolQuery) *query.BoolQuery {
		return q.Filter(query.NewRangeQuery(field).Lte(to).Gte(from))
	}
}

func WithCreatedAtBy(created_at time.Time) Filter {
	return func(q *query.BoolQuery) *query.BoolQuery {
		return q.Must(query.NewMatchQuery("created_at", created_at))
	}
}

func WithUpdatedAtBy(updated_at time.Time) Filter {
	return func(q *query.BoolQuery) *query.BoolQuery {
		return q.Must(query.NewMatchQuery("updated_by", updated_at))
	}
}

func WithCreatedAtAfter(t time.Time) Filter {
	return func(q *query.BoolQuery) *query.BoolQuery {
		return q.Filter(query.NewRangeQuery("created_at").Gt(t))
	}
}

func WithCreatedAtBefore(t time.Time) Filter {
	return func(q *query.BoolQuery) *query.BoolQuery {
		return q.Filter(query.NewRangeQuery("created_at").Lt(t))
	}
}

func WithUpdatedAtAfter(t time.Time) Filter {
	return func(q *query.BoolQuery) *query.BoolQuery {
		return q.Filter(query.NewRangeQuery("updated_at").Gt(t))
	}
}

func WithUpdateAtBefore(t time.Time) Filter {
	return func(q *query.BoolQuery) *query.BoolQuery {
		return q.Filter(query.NewRangeQuery("updated_at").Lt(t))
	}
}
