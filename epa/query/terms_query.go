package query

type TermsQuery struct {
	field  string
	values []interface{}
}

func NewTermsQuery(field string, values ...interface{}) *TermsQuery {
	return &TermsQuery{
		field:  field,
		values: values,
	}
}

func (q *TermsQuery) Source() (interface{}, error) {
	//
	// {
	//   "terms" : {
	//     "user" : ["kimchy", "elasticsearch"]
	//   }
	// }
	//
	source := make(map[string]interface{})
	terms := make(map[string]interface{})
	source["terms"] = terms
	terms[q.field] = q.values
	return source, nil
}
