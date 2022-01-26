package GrpcEngine

import (
	"github.com/shopspring/decimal"
)

func (d *Decimal) ToDecimal() decimal.Decimal {
	return decimal.New(d.Val, d.Exp)
}

func (d *Decimal) ToNullDecimal() decimal.NullDecimal {
	return decimal.NewNullDecimal(d.ToDecimal())
}
