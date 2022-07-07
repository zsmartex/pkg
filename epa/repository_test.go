package epa

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/zsmartex/pkg/v2/infrastucture/elasticsearch"
)

type Order struct {
	ID        int64           `json:"id"`
	UserID    int64           `json:"user_id"`
	Price     decimal.Decimal `json:"price"`
	State     int64           `json:"state"`
	CreatedAt time.Time       `json:"created_at"`
}

func (o Order) IndexName() string {
	return "orders"
}

func newRepo() (Repository[Order], error) {
	client, err := elasticsearch.New(&elasticsearch.Config{
		URL:      "http://localhost:9200",
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
			UserID:    1,
			Price:     decimal.NewFromFloat(20.7),
			State:     100,
			CreatedAt: time.Now(),
		},
		{
			ID:        2,
			UserID:    2,
			Price:     decimal.NewFromFloat(48.4),
			State:     200,
			CreatedAt: time.Now(),
		},
		{
			ID:        3,
			UserID:    1,
			Price:     decimal.NewFromFloat(15.0),
			State:     200,
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
			Filters: []Filter{
				WithFieldLessThanOrEqualTo("price", "30"),
			},
		},
	)
	if err != nil {
		t.Error(err)
	}

	t.Error(result.Values)
	t.Error(result.TotalHits)
}
