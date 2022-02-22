package GrpcSymbol

import "github.com/zsmartex/pkg"

func (s *Symbol) ToSymbol() pkg.Symbol {
	return pkg.Symbol{BaseCurrency: s.BaseCurrency, QuoteCurrency: s.QuoteCurrency}
}
