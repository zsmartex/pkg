package aggregation

type Aggregation interface {
	Source() (interface{}, error)
}

type Aggregations map[string]Aggregation

func (a Aggregations) Source() (interface{}, error) {
	//
	// {
	//   "aggs" : {
	//     "avg_grade" : { "avg" : { "field" : "grade" } }
	//   }
	// }
	//
	source := make(map[string]interface{})
	for name, agg := range a {
		src, err := agg.Source()
		if err != nil {
			return nil, err
		}

		source[name] = src
	}
	return source, nil
}
