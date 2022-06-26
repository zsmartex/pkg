package epa

import (
	"time"

	"github.com/olivere/elastic/v7"
)

func WithFieldEqual(field string, value interface{}) Filter {
	return func(query *elastic.BoolQuery) *elastic.BoolQuery {
		return query.Must(elastic.NewMatchQuery(field, value))
	}
}

func WithFieldNotEqual(field string, value interface{}) Filter {
	return func(query *elastic.BoolQuery) *elastic.BoolQuery {
		return query.MustNot(elastic.NewMatchQuery(field, value))
	}
}

func WithFieldGreaterThan(field string, value interface{}) Filter {
	return func(query *elastic.BoolQuery) *elastic.BoolQuery {
		return query.Must(elastic.NewRangeQuery(field).Gt(value))
	}
}

func WithFieldLessThan(field string, value interface{}) Filter {
	return func(query *elastic.BoolQuery) *elastic.BoolQuery {
		return query.Must(elastic.NewRangeQuery(field).Lt(value))
	}
}

func WithFieldGreaterThanOrEqualTo(field string, value interface{}) Filter {
	return func(query *elastic.BoolQuery) *elastic.BoolQuery {
		return query.Must(elastic.NewRangeQuery(field).Gte(value))
	}
}

func WithFieldLessThanOrEqualTo(field string, value interface{}) Filter {
	return func(query *elastic.BoolQuery) *elastic.BoolQuery {
		return query.Must(elastic.NewRangeQuery(field).Lte(value))
	}
}

func WithFieldIn(field string, value interface{}) Filter {
	return func(query *elastic.BoolQuery) *elastic.BoolQuery {
		return query.Must(elastic.NewTermsQuery(field, value))
	}
}

func WithFieldNotIn(field string, value interface{}) Filter {
	return func(query *elastic.BoolQuery) *elastic.BoolQuery {
		return query.MustNot(elastic.NewTermsQuery(field, value))
	}
}

func WithFieldLike(field string, value interface{}) Filter {
	return func(query *elastic.BoolQuery) *elastic.BoolQuery {
		return query.Should(elastic.NewMatchQuery(field, value))
	}
}

func WithNotLIKE(field string, value interface{}) Filter {
	return func(query *elastic.BoolQuery) *elastic.BoolQuery {
		return query.MustNot(elastic.NewMatchQuery(field, value))
	}
}

func WithFieldIsNull(field string) Filter {
	return func(query *elastic.BoolQuery) *elastic.BoolQuery {
		return query.Must(elastic.NewExistsQuery(field))
	}
}

func WithID(id string) Filter {
	return func(query *elastic.BoolQuery) *elastic.BoolQuery {
		return query.Must(elastic.NewMatchQuery("_id", id))
	}
}

func WithDateRange(field string, from time.Time, to time.Time) *elastic.BoolQuery {
	return elastic.
		NewBoolQuery().
		Filter(elastic.NewRangeQuery(field).Gte(from)).
		Filter(elastic.NewRangeQuery(field).Lte(to))
}

func WithCreatedAtBy(created_at time.Time) Filter {
	return func(query *elastic.BoolQuery) *elastic.BoolQuery {
		return query.Must(elastic.NewMatchQuery("created_at", created_at))
	}
}

func WithUpdatedAtBy(updated_at time.Time) Filter {
	return func(query *elastic.BoolQuery) *elastic.BoolQuery {
		return query.Must(elastic.NewMatchQuery("updated_by", updated_at))
	}
}

func WithCreatedAtAfter(t time.Time) Filter {
	return func(query *elastic.BoolQuery) *elastic.BoolQuery {
		return query.Filter(elastic.NewRangeQuery("created_at").Gt(t))
	}
}

func WithCreateAtBefore(t time.Time) Filter {
	return func(query *elastic.BoolQuery) *elastic.BoolQuery {
		return query.Filter(elastic.NewRangeQuery("created_at").Lt(t))
	}
}

func WithUpdatedAtAfter(t time.Time) Filter {
	return func(query *elastic.BoolQuery) *elastic.BoolQuery {
		return query.Filter(elastic.NewRangeQuery("updated_at").Gt(t))
	}
}

func WithUpdateAtBefore(t time.Time) Filter {
	return func(query *elastic.BoolQuery) *elastic.BoolQuery {
		return query.Filter(elastic.NewRangeQuery("updated_at").Lt(t))
	}
}
