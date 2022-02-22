package pkg

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OrderSide string
type OrderType string

var (
	SideSell OrderSide = "sell"
	SideBuy  OrderSide = "buy"
)

var (
	TypeLimit  OrderType = "limit"
	TypeMarket OrderType = "market"
)

type Order struct {
	ID             int64           `json:"id"`
	UUID           uuid.UUID       `json:"uuid"`
	MemberID       int64           `json:"member_id"`
	Symbol         Symbol          `json:"symbol"`
	Side           OrderSide       `json:"side"`
	Type           OrderType       `json:"type"`
	Price          decimal.Decimal `json:"price"`
	StopPrice      decimal.Decimal `json:"stop_price"`
	Quantity       decimal.Decimal `json:"quantity"`
	FilledQuantity decimal.Decimal `json:"filled_quantity"`
	Cancelled      bool            `json:"canceled"`
	Fake           bool            `json:"fake"`
	CreatedAt      time.Time       `json:"created_at"`
}

type OrderKey struct {
	ID        int64           `json:"id"`
	UUID      uuid.UUID       `json:"uuid"`
	Symbol    Symbol          `json:"symbol"`
	Side      OrderSide       `json:"side"`
	Price     decimal.Decimal `json:"price"`
	StopPrice decimal.Decimal `json:"stop_price"`
	Fake      bool            `json:"fake"`
	CreatedAt time.Time       `json:"created_at"`
}

func (o *Order) Key() *OrderKey {
	return &OrderKey{
		ID:        o.ID,
		UUID:      o.UUID,
		Symbol:    o.Symbol,
		Side:      o.Side,
		Price:     o.Price,
		StopPrice: o.StopPrice,
		Fake:      o.Fake,
		CreatedAt: o.CreatedAt,
	}
}

func (o *Order) IsBid() bool {
	return o.Side == SideBuy
}

func (o *Order) IsAsk() bool {
	return o.Side == SideSell
}

func (o *Order) IsFake() bool {
	return o.Fake
}

func (o *Order) Fill(quantity decimal.Decimal) {
	o.FilledQuantity = o.FilledQuantity.Add(quantity)
	if o.Filled() {
		o.Cancel()
	}
}

func (o *Order) Filled() bool {
	return o.Quantity.Equal(o.FilledQuantity)
}

func (o *Order) UnfilledQuantity() decimal.Decimal {
	return o.Quantity.Sub(o.FilledQuantity)
}

func (o *Order) IsCrossed(price decimal.Decimal) bool {
	if o.Side == SideSell {
		return price.GreaterThanOrEqual(o.Price)
	} else {
		return price.LessThanOrEqual(o.Price)
	}
}

func (o *Order) Cancel() {
	o.Cancelled = true
}
