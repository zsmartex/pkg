package order

import (
	"time"

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
	ID             uint64          `json:"id"`
	MemberID       uint64          `json:"member_id"`
	Symbol         string          `json:"symbol"`
	Side           OrderSide       `json:"side"`
	Type           OrderType       `json:"type"`
	Price          decimal.Decimal `json:"price"`
	StopPrice      decimal.Decimal `json:"stop_price"`
	Quantity       decimal.Decimal `json:"quantity"`
	FilledQuantity decimal.Decimal `json:"filled_quantity"`
	Cancelled      bool            `json:"canceled"`
	Fake           bool            `json:"is_fake"`
	CreatedAt      time.Time       `json:"created_at"`
}

type OrderKey struct {
	ID        uint64
	Symbol    string
	Side      OrderSide
	Price     decimal.Decimal
	StopPrice decimal.Decimal
	CreatedAt time.Time
}

func (o *Order) Key() *OrderKey {
	return &OrderKey{
		ID:        o.ID,
		Symbol:    o.Symbol,
		Side:      o.Side,
		Price:     o.Price,
		StopPrice: o.StopPrice,
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
