package filters

import (
	"github.com/zsmartex/pkg/v2/mpa"
	"go.mongodb.org/mongo-driver/bson"
)

func applyFilter(filters ...mpa.Filter) []bson.M {
	value := []bson.M{}
	for _, filter := range filters {
		k, v := filter()
		value = append(value, bson.M{k: v})
	}

	return value
}

func WithAnd(filters ...mpa.Filter) mpa.Filter {
	return func() (k string, v interface{}) {
		return "$and", applyFilter(filters...)
	}
}

func WithOr(filters ...mpa.Filter) mpa.Filter {
	return func() (k string, v interface{}) {
		return "$or", applyFilter(filters...)
	}
}

func WithNot(filters ...mpa.Filter) mpa.Filter {
	return func() (k string, v interface{}) {
		return "$not", applyFilter(filters...)
	}
}

func WithNor(filters ...mpa.Filter) mpa.Filter {
	return func() (k string, v interface{}) {
		return "$nor", applyFilter(filters...)
	}
}
