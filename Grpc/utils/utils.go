package GrpcUtils

import (
	"github.com/shopspring/decimal"
)

func (d *Decimal) ToDecimal() decimal.Decimal {
	if d == nil {
		d = &Decimal{
			Val: 0,
			Exp: 0,
		}
	}

	return decimal.New(d.Val, d.Exp)
}

func (d *Decimal) ToNullDecimal() decimal.NullDecimal {
	return decimal.NewNullDecimal(d.ToDecimal())
}
