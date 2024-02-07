package aggregation

type CountAggregation struct{}

func NewCountAggregation() *CountAggregation {
	return &CountAggregation{}
}

func (a *CountAggregation) Source() (interface{}, error) {
	source := make(map[string]interface{})
	sum := make(map[string]interface{})
	source["value_count"] = sum
	sum["field"] = "id"

	return source, nil
}
