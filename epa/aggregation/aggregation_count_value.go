package aggregation

import (
	"encoding/json"
)

type AggregationCountValue struct {
	Aggregations

	Value int64 `json:"value,omitempty"`
}

func (a *AggregationCountValue) UnmarshalJSON(data []byte) error {
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
