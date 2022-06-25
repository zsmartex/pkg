package epa

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	ID        int64           `json:"id"`
	UserID    int64           `json:"user_id"`
	Price     decimal.Decimal `json:"price"`
	State     int64           `json:"state"`
	CreatedAt time.Time       `json:"created_at"`
}

func TestCreate(t *testing.T) {
	repo, err := New[Order](&Config{
		URL:      "http://localhost:9200",
		Username: "elastic",
		Password: "elastic",
	})
	if err != nil {
		t.Error(err)
	}

	orders := []*Order{
		{
			ID:        1,
			UserID:    1,
			Price:     decimal.NewFromFloat(10.0),
			State:     100,
			CreatedAt: time.Now(),
		},
		{
			ID:        2,
			UserID:    2,
			Price:     decimal.NewFromFloat(20.0),
			State:     200,
			CreatedAt: time.Now(),
		},
		{
			ID:        3,
			UserID:    1,
			Price:     decimal.NewFromFloat(30.0),
			State:     200,
			CreatedAt: time.Now(),
		},
	}

	for _, order := range orders {
		r, err := repo.Create(context.Background(), "orders", fmt.Sprint(order.ID), order)
		if err != nil {
			t.Error(err)
		}

		t.Error(r)
	}
}

func TestFind(t *testing.T) {
	repo, err := New[Order](&Config{
		URL:      "http://localhost:9200",
		Username: "elastic",
		Password: "elastic",
	})
	if err != nil {
		t.Error(err)
	}

	result, err := repo.Find(
		context.Background(),
		"orders",
		Query{
			Limit: 1,
			Filters: []Filter{
				WithCreateAtBefore(time.Now()),
			},
		},
	)
	if err != nil {
		t.Error(err)
	}

	t.Error(result.Values)
	t.Error(result.TotalHits)
}
