package pkg

import (
	"github.com/shopspring/decimal"
)

// Trade .
type Trade struct {
	ID         int64           `json:"id"`
	Symbol     Symbol          `json:"symbol"`
	Price      decimal.Decimal `json:"price"`
	Quantity   decimal.Decimal `json:"quantity"`
	Total      decimal.Decimal `json:"total"`
	MakerOrder Order           `json:"maker"`
	TakerOrder Order           `json:"taker"`
}

func (t *Trade) BuyOrder() Order {
	if t.MakerOrder.Side == SideBuy {
		return t.MakerOrder
	} else {
		return t.TakerOrder
	}
}

func (t *Trade) SellOrder() Order {
	if t.MakerOrder.Side == SideSell {
		return t.MakerOrder
	} else {
		return t.TakerOrder
	}
}
