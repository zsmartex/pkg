package epa

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/olivere/elastic/v7"
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
		URL:      "http://zsmartex.com:9200",
		Username: "elastic",
		Password: "elastic",
	})
	if err != nil {
		return Repository[Order]{}, err
	}

	return New(client, Order{}), nil
}

func TestCreate(t *testing.T) {
	repo, err := newRepo()
	if err != nil {
		t.Error(err)
	}

	orders := []*Order{
		{
			ID:        1,
			CreatedAt: time.Now(),
		},
		{
			ID:        2,
			CreatedAt: time.Now(),
		},
		{
			ID:        3,
			CreatedAt: time.Now(),
		},
	}

	for _, order := range orders {
		r, err := repo.Create(context.Background(), fmt.Sprint(order.ID), order)
		if err != nil {
			t.Error(err)
		}

		t.Error(r)
	}
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
			Addons: func(searchService *elastic.SearchService) *elastic.SearchService {
				return searchService.Aggregation("sales", elastic.NewDateHistogramAggregation().Field("created_at").CalendarInterval("day").Format("yyyy-MM-dd"))
			},
		},
	)
	if err != nil {
		t.Error(err)
	}

	t.Error(result.Values)
	t.Error(result.TotalHits)
}
