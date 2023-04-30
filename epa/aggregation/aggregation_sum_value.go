package aggregation

import (
	"encoding/json"

	"github.com/shopspring/decimal"
)

type AggregationSumValue struct {
	Aggregations

	Value decimal.Decimal `json:"value,omitempty"`
}

func (a *AggregationSumValue) UnmarshalJSON(data []byte) error {
	var aggs map[string]json.RawMessage
	if err := json.Unmarshal(data, &aggs); err != nil {
		return err
	}
	if v, ok := aggs["value"]; ok && v != nil {
		if err := json.Unmarshal(v, &a.Value); err != nil {
			return err
		}
	}
	a.Aggregations = aggs
	return nil
}
