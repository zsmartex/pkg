package mpa

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Filter func() bson.E

type Ordering string

var (
	OrderingAsc  Ordering = "asc"
	OrderingDesc Ordering = "desc"
)

type OptionsFind struct {
	Page     int
	Limit    int
	OrderBy  string
	Ordering Ordering
}

func NewOptionsFind(optionsFind OptionsFind) *options.FindOptions {
	opts := options.Find()

	if optionsFind.Page > 0 {
		opts.SetLimit(int64(optionsFind.Page))
	}

	if optionsFind.Limit > 0 {
		opts.SetLimit(int64(optionsFind.Limit))
	}

	if len(optionsFind.OrderBy) > 0 && len(optionsFind.Ordering) > 0 {
		order := 1
		if optionsFind.Ordering == OrderingDesc {
			order = -1
		}

		opts.SetSort(bson.M{optionsFind.OrderBy: order})
	}

	return opts
}

func ApplyFilters(fs ...Filter) bson.D {
	var result []bson.E

	if len(fs) == 0 {
		return bson.D{}
	}

	for _, f := range fs {
		result = append(result, f())
	}

	return result
}
