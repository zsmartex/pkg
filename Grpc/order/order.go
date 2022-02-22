package GrpcOrder

import (
	"github.com/google/uuid"
	"github.com/zsmartex/pkg"
)

func (r *OrderKey) ToOrderKey() *pkg.OrderKey {
	uuid, _ := uuid.FromBytes(r.Uuid)

	return &pkg.OrderKey{
		ID:        r.Id,
		UUID:      uuid,
		Symbol:    r.Symbol.ToSymbol(),
		Side:      pkg.OrderSide(r.Side),
		Price:     r.Price.ToDecimal(),
		StopPrice: r.StopPrice.ToDecimal(),
		Fake:      r.Fake,
		CreatedAt: r.CreatedAt.AsTime(),
	}
}

func (r *Order) ToOrder() *pkg.Order {
	uuid, _ := uuid.FromBytes(r.Uuid)

	return &pkg.Order{
		ID:             r.Id,
		UUID:           uuid,
		Symbol:         r.Symbol.ToSymbol(),
		MemberID:       r.MemberId,
		Side:           pkg.OrderSide(r.Side),
		Type:           pkg.OrderType(r.Type),
		Price:          r.Price.ToDecimal(),
		StopPrice:      r.StopPrice.ToDecimal(),
		Quantity:       r.Quantity.ToDecimal(),
		FilledQuantity: r.FilledQuantity.ToDecimal(),
		Cancelled:      r.Cancelled,
		Fake:           r.Fake,
		CreatedAt:      r.CreatedAt.AsTime(),
	}
}
