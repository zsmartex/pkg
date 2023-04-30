package query

type MultiMatchQuery struct {
	match  interface{}
	fields []string
}

func NewMultiMatchQuery(match interface{}, fields ...string) *MultiMatchQuery {
	return &MultiMatchQuery{
		match:  match,
		fields: fields,
	}
}

func (q *MultiMatchQuery) Source() (interface{}, error) {
	//
	// {
	//   "multi_match" : {
	//     "query" : "this is a test",
	//     "fields" : [ "subject", "message" ]
	//   }
	// }
	//

	source := make(map[string]interface{})
	multiMatch := make(map[string]interface{})
	source["multi_match"] = multiMatch

	multiMatch["query"] = q.match

	var fields []string
	for _, field := range q.fields {
		fields = append(fields, field)
	}
	if fields == nil {
		multiMatch["fields"] = []string{}
	} else {
		multiMatch["fields"] = fields
	}

	return source, nil
}
