package epa

import (
	"context"
	"testing"
	"time"

	"github.com/zsmartex/pkg/v2/epa/aggregation"
	"github.com/zsmartex/pkg/v2/infrastucture/elasticsearch"
)

type Order struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

func (o Order) IndexName() string {
	return "pg.orders"
}

func newRepo() (Repository[Order], error) {
	client, err := elasticsearch.New(&elasticsearch.Config{
		URL:      []string{"http://demo.zsmartex.com:9200"},
		Username: "elastic",
		Password: "elastic",
	})
	if err != nil {
		return Repository[Order]{}, err
	}

	return New(client, Order{}), nil
}

func TestCreate(t *testing.T) {
	//repo, err := newRepo()
	//if err != nil {
	//	t.Error(err)
	//}
	//
	//orders := []*Order{
	//	{
	//		ID:        1,
	//		CreatedAt: time.Now(),
	//	},
	//	{
	//		ID:        2,
	//		CreatedAt: time.Now(),
	//	},
	//	{
	//		ID:        3,
	//		CreatedAt: time.Now(),
	//	},
	//}
	//
	//for _, order := range orders {
	//	r, err := repo.Create(context.Background(), fmt.Sprint(order.ID), order)
	//	if err != nil {
	//		t.Error(err)
	//	}
	//
	//	t.Error(r)
	//}
}

func TestFind(t *testing.T) {
	repo, err := newRepo()
	if err != nil {
		t.Error(err)
	}

	result, err := repo.Find(
		context.Background(),
		Query{
			Limit: 0,
			Filters: []Filter{
				WithDateRange("created_at", "2022-07-05T00:00:00.551Z", "2022-07-08T17:10:26.697Z"),
				WithFieldLessThan("price", 12),
			},
			Aggregations: map[string]aggregation.Aggregation{
				"price": aggregation.NewDateHistogramAggregation("created_at").FixedInterval("1d"),
			},
		},
	)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(result.Values)
	t.Log(result.TotalHits)
	t.Log(result.Aggregations)
}
