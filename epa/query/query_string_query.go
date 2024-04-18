package query

type QueryStringQuery struct {
	query  string
	fields []string
}

func NewQueryStringQuery(query string, fields ...string) *QueryStringQuery {
	return &QueryStringQuery{
		query:  query,
		fields: fields,
	}
}

func (q *QueryStringQuery) Source() (interface{}, error) {
	//
	// {
	//   "query_string" : {
	//     "query" : "this is a test",
	//     "fields" : [ "subject", "message" ]
	//   }
	// }
	//
	source := make(map[string]interface{})
	queryString := make(map[string]interface{})
	source["query_string"] = queryString
	queryString["query"] = q.query
	var fields []string
	for _, field := range q.fields {
		fields = append(fields, field)
	}
	if fields == nil {
		queryString["fields"] = []string{}
	} else {
		queryString["fields"] = fields
	}
	return source, nil
}
