package trade

import (
	"github.com/shopspring/decimal"

	"github.com/zsmartex/pkg/order"
)

// Trade .
type Trade struct {
	ID         uint64          `json:"id"`
	Symbol     string          `json:"symbol"`
	Price      decimal.Decimal `json:"price"`
	Quantity   decimal.Decimal `json:"quantity"`
	Total      decimal.Decimal `json:"total"`
	MakerOrder order.Order     `json:"maker"`
	TakerOrder order.Order     `json:"taker"`
}

func (t *Trade) BuyOrder() order.Order {
	if t.MakerOrder.Side == order.SideBuy {
		return t.MakerOrder
	} else {
		return t.TakerOrder
	}
}

func (t *Trade) SellOrder() order.Order {
	if t.MakerOrder.Side == order.SideSell {
		return t.MakerOrder
	} else {
		return t.TakerOrder
	}
}
