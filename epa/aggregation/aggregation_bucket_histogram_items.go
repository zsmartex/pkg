package aggregation

import "encoding/json"

// AggregationBucketHistogramItem is a single bucket of an AggregationBucketHistogramItems structure.
type AggregationBucketHistogramItem struct {
	Aggregations

	Key         float64 //`json:"key"`
	KeyAsString *string //`json:"key_as_string"`
	DocCount    int64   //`json:"doc_count"`
}

// UnmarshalJSON decodes JSON data and initializes an AggregationBucketHistogramItem structure.
func (a *AggregationBucketHistogramItem) UnmarshalJSON(data []byte) error {
	var aggs map[string]json.RawMessage
	if err := json.Unmarshal(data, &aggs); err != nil {
		return err
	}
	if v, ok := aggs["key"]; ok && v != nil {
		if err := json.Unmarshal(v, &a.Key); err != nil {
			return err
		}
	}
	if v, ok := aggs["key_as_string"]; ok && v != nil {
		if err := json.Unmarshal(v, &a.KeyAsString); err != nil {
			return err
		}
	}
	if v, ok := aggs["doc_count"]; ok && v != nil {
		if err := json.Unmarshal(v, &a.DocCount); err != nil {
			return err
		}
	}
	a.Aggregations = aggs
	return nil
}

type AggregationBucketHistogramItems struct {
	Aggregations

	Buckets []*AggregationBucketHistogramItem `json:"buckets,omitempty"`
	Meta    map[string]interface{}            `json:"meta,omitempty"`
}

// UnmarshalJSON decodes JSON data and initializes an AggregationBucketHistogramItems structure.
func (a *AggregationBucketHistogramItems) UnmarshalJSON(data []byte) error {
	var aggs map[string]json.RawMessage
	if err := json.Unmarshal(data, &aggs); err != nil {
		return err
	}
	if v, ok := aggs["buckets"]; ok && v != nil {
		if err := json.Unmarshal(v, &a.Buckets); err != nil {
			return err
		}
	}
	if v, ok := aggs["meta"]; ok && v != nil {
		if err := json.Unmarshal(v, &a.Meta); err != nil {
			return err
		}
	}
	a.Aggregations = aggs
	return nil
}
