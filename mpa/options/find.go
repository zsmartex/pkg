package options

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Ordering string

var (
	OrderingAsc  Ordering = "asc"
	OrderingDesc Ordering = "desc"
)

type findOptions struct {
	options.FindOptions
}

func Find() *findOptions {
	return &findOptions{}
}

func (o *findOptions) WithPage(page int) *findOptions {
	p := int64(page)
	o.Skip = &p
	return o
}

func (o *findOptions) WithLimit(limit int) *findOptions {
	l := int64(limit)
	o.Limit = &l
	return o
}

func (o *findOptions) WithOrder(orderBy string, ordering Ordering) *findOptions {
	order := 1
	if ordering == OrderingDesc {
		order = -1
	}

	o.Sort = bson.M{orderBy: order}
	return o
}
