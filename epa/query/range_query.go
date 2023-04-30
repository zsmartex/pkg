package query

type RangeQuery struct {
	field string
	lte   interface{}
	lt    interface{}
	gte   interface{}
	gt    interface{}
}

func NewRangeQuery(field string) *RangeQuery {
	return &RangeQuery{
		field: field,
	}
}

func (q *RangeQuery) Lte(lte interface{}) *RangeQuery {
	q.lte = lte
	return q
}

func (q *RangeQuery) Lt(lt interface{}) *RangeQuery {
	q.lt = lt
	return q
}

func (q *RangeQuery) Gte(gte interface{}) *RangeQuery {
	q.gte = gte
	return q
}

func (q *RangeQuery) Gt(gt interface{}) *RangeQuery {
	q.gt = gt
	return q
}

func (q *RangeQuery) Source() (interface{}, error) {
	//
	// {
	//   "range" : {
	//     "age" : {
	//       "gte" : 10,
	//       "lte" : 20
	//     }
	//   }
	// }
	//
	source := make(map[string]interface{})
	rangeQ := make(map[string]interface{})
	source["range"] = rangeQ
	rangeQ[q.field] = make(map[string]interface{})
	if q.gte != nil {
		rangeQ[q.field].(map[string]interface{})["gte"] = q.gte
	}
	if q.lte != nil {
		rangeQ[q.field].(map[string]interface{})["lte"] = q.lte
	}
	if q.gt != nil {
		rangeQ[q.field].(map[string]interface{})["gt"] = q.gt
	}
	if q.lt != nil {
		rangeQ[q.field].(map[string]interface{})["lt"] = q.lt
	}
	return source, nil
}
