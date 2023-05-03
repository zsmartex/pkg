package aggregation

type SumAggregation struct {
	field string
}

func NewSumAggregation(field string) *SumAggregation {
	return &SumAggregation{
		field: field,
	}
}

func (a *SumAggregation) Source() (interface{}, error) {
	// Example
	//
	// {
	// 		"hat_prices" : {
	// 			"sum": {
	// 				"field": "price"
	// 			}
	// 		}
	// }
	//
	//
	// This method returns only the  { "sum": { "field": "..." } } part.

	source := make(map[string]interface{})
	sum := make(map[string]interface{})
	source["sum"] = sum
	sum["field"] = a.field

	return source, nil
}
