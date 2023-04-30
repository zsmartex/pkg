package query

type ExistsQuery struct {
	field string
}

func NewExistsQuery(field string) *ExistsQuery {
	return &ExistsQuery{
		field: field,
	}
}

func (q *ExistsQuery) Source() (interface{}, error) {
	//
	// {
	//   "exists" : {
	//     "field" : "user"
	//   }
	// }
	//
	source := make(map[string]interface{})
	exists := make(map[string]interface{})
	source["exists"] = exists
	exists["field"] = q.field
	return source, nil
}
