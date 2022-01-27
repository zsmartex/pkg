package GrpcOrder

import (
	"github.com/google/uuid"
	"github.com/zsmartex/pkg/order"
)

func (r *OrderKey) ToOrderKey() *order.OrderKey {
	uuid, _ := uuid.FromBytes(r.Uuid)

	return &order.OrderKey{
		ID:        r.Id,
		UUID:      uuid,
		Symbol:    r.Symbol,
		Side:      order.OrderSide(r.Side),
		Price:     r.Price.ToDecimal(),
		StopPrice: r.StopPrice.ToDecimal(),
		Fake:      r.Fake,
		CreatedAt: r.CreatedAt.AsTime(),
	}
}

func (r *Order) ToOrder() *order.Order {
	uuid, _ := uuid.FromBytes(r.Uuid)

	return &order.Order{
		ID:             r.Id,
		UUID:           uuid,
		Symbol:         r.Symbol,
		MemberID:       r.MemberId,
		Side:           order.OrderSide(r.Side),
		Type:           order.OrderType(r.Type),
		Price:          r.Price.ToDecimal(),
		StopPrice:      r.StopPrice.ToDecimal(),
		Quantity:       r.Quantity.ToDecimal(),
		FilledQuantity: r.FilledQuantity.ToDecimal(),
		Cancelled:      r.Cancelled,
		Fake:           r.Fake,
		CreatedAt:      r.CreatedAt.AsTime(),
	}
}
