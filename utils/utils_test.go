package utils

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestCompareDiff(t *testing.T) {
	type Order struct {
		Price decimal.Decimal
	}

	order1 := Order{
		Price: decimal.NewFromFloat(1.1),
	}

	order2 := Order{
		Price: decimal.NewFromFloat(1.2),
	}

	order3 := Order{
		Price: decimal.NewFromFloat(1.3),
	}

	if err := CompareDiff(&order1, order2, order3); err != nil {
		t.Error(err)
	}

	t.Log(order1)
}
