package query

type MatchQuery struct {
	query interface{}
	field string
}

func NewMatchQuery(field string, query interface{}) *MatchQuery {
	return &MatchQuery{
		query: query,
		field: field,
	}
}

func (q *MatchQuery) Source() (interface{}, error) {
	//
	// {
	//   "match" : {
	//     "subject" :  "this is a test",
	//   }
	// }
	//

	source := make(map[string]interface{})
	match := make(map[string]interface{})
	source["match"] = match
	match[q.field] = q.query

	return source, nil
}
