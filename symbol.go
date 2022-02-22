package pkg

import "strings"

type Symbol struct {
	BaseCurrency  string `json:"base_currency"`
	QuoteCurrency string `json:"quote_currency"`
}

func (s *Symbol) String() string {
	return s.ToSymbol("_")
}

func (s *Symbol) ToSymbol(joinChar string) string {
	return strings.Join([]string{s.BaseCurrency, s.QuoteCurrency}, joinChar)
}
