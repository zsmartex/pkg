package query

type BoolQuery struct {
	Query
	mustClauses      []Query
	mustNotClauses   []Query
	filterClauses    []Query
	shouldClauses    []Query
	shouldNotClauses []Query
	boost            *float64
}

// Creates a new bool query.
func NewBoolQuery() *BoolQuery {
	return &BoolQuery{
		mustClauses:      make([]Query, 0),
		mustNotClauses:   make([]Query, 0),
		filterClauses:    make([]Query, 0),
		shouldClauses:    make([]Query, 0),
		shouldNotClauses: make([]Query, 0),
	}
}

func (q *BoolQuery) Must(queries ...Query) *BoolQuery {
	q.mustClauses = queries

	return q
}

func (q *BoolQuery) MustNot(queries ...Query) *BoolQuery {
	q.mustNotClauses = queries

	q.MustNot(NewMultiMatchQuery(""))

	return q
}

func (q *BoolQuery) Filter(queries ...Query) *BoolQuery {
	q.filterClauses = queries

	return q
}

func (q *BoolQuery) Should(queries ...Query) *BoolQuery {
	q.shouldClauses = queries

	return q
}

func (q *BoolQuery) ShouldNot(queries ...Query) *BoolQuery {
	q.shouldNotClauses = queries

	return q
}

func (q *BoolQuery) Boost(boost float64) *BoolQuery {
	q.boost = &boost

	return q
}

// Creates the query source for the bool query.
func (q *BoolQuery) Source() (interface{}, error) {
	// {
	//	"bool" : {
	//		"must" : {
	//			"term" : { "user" : "kimchy" }
	//		},
	//		"must_not" : {
	//			"range" : {
	//				"age" : { "from" : 10, "to" : 20 }
	//			}
	//		},
	//    "filter" : [
	//      ...
	//    ]
	//		"should" : [
	//			{
	//				"term" : { "tag" : "wow" }
	//			},
	//			{
	//				"term" : { "tag" : "elasticsearch" }
	//			}
	//		],
	//		"boost" : 1.0
	//	}
	// }
	query := make(map[string]interface{})
	boolClause := make(map[string]interface{})
	query["bool"] = boolClause

	if len(q.mustClauses) > 0 {
		var clauses []interface{}
		for _, subQuery := range q.mustClauses {
			src, err := subQuery.Source()
			if err != nil {
				return nil, err
			}
			clauses = append(clauses, src)
		}
		boolClause["must"] = clauses
	}

	if len(q.mustNotClauses) > 0 {
		var clauses []interface{}
		for _, subQuery := range q.mustNotClauses {
			src, err := subQuery.Source()
			if err != nil {
				return nil, err
			}
			clauses = append(clauses, src)
		}
		boolClause["must_not"] = clauses
	}

	if len(q.filterClauses) > 0 {
		var clauses []interface{}
		for _, subQuery := range q.filterClauses {
			src, err := subQuery.Source()
			if err != nil {
				return nil, err
			}
			clauses = append(clauses, src)
		}
		boolClause["filter"] = clauses
	}

	if len(q.shouldClauses) > 0 {
		var clauses []interface{}
		for _, subQuery := range q.shouldClauses {
			src, err := subQuery.Source()
			if err != nil {
				return nil, err
			}
			clauses = append(clauses, src)
		}
		boolClause["should"] = clauses
	}

	if len(q.shouldNotClauses) > 0 {
		var clauses []interface{}
		for _, subQuery := range q.shouldNotClauses {
			src, err := subQuery.Source()
			if err != nil {
				return nil, err
			}
			clauses = append(clauses, src)
		}
		boolClause["should_not"] = clauses
	}

	if q.boost != nil {
		boolClause["boost"] = *q.boost
	}

	return query, nil
}
