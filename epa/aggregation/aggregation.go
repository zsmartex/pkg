package aggregation

import "encoding/json"

type Aggregation interface {
	Source() (interface{}, error)
}

type Aggregations map[string]json.RawMessage

func (a Aggregations) DateHistogram(name string) (items *AggregationBucketHistogramItems, found bool) {
	if raw, found := a[name]; found {
		agg := new(AggregationBucketHistogramItems)
		if raw == nil {
			return agg, true
		}
		if err := json.Unmarshal(raw, agg); err == nil {
			return agg, true
		}
	}

	return nil, false
}
