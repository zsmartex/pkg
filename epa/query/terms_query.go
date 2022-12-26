package query

type TermsQuery[T any] struct {
	field  string
	values []T
}

func NewTermsQuery[T any](field string, values ...T) *TermsQuery[T] {
	return &TermsQuery[T]{
		field:  field,
		values: values,
	}
}

func (q *TermsQuery[T]) Source() (interface{}, error) {
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
