package aggregation

import "github.com/zsmartex/pkg/v3/queries"

type DateHistogramAggregation struct {
	field            string
	fixedInterval    string
	calendarInterval string
	order            string
	ordering         queries.Ordering
	minDocCount      *int64
	timeZone         string
	format           string
	offset           string
}

func NewDateHistogramAggregation(field string) *DateHistogramAggregation {
	return &DateHistogramAggregation{
		field:    field,
		ordering: queries.OrderingAsc,
	}
}

func (a *DateHistogramAggregation) FixedInterval(fixedInterval string) *DateHistogramAggregation {
	a.fixedInterval = fixedInterval
	return a
}

func (a *DateHistogramAggregation) CalendarInterval(calendarInterval string) *DateHistogramAggregation {
	a.calendarInterval = calendarInterval
	return a
}

func (a *DateHistogramAggregation) Order(order string, ordering queries.Ordering) *DateHistogramAggregation {
	a.order = order
	a.ordering = ordering
	return a
}

func (a *DateHistogramAggregation) OrderByCount(ordering queries.Ordering) *DateHistogramAggregation {
	a.order = "_count"
	a.ordering = ordering
	return a
}

func (a *DateHistogramAggregation) MinDocCount(minDocCount int64) *DateHistogramAggregation {
	a.minDocCount = &minDocCount
	return a
}

func (a *DateHistogramAggregation) TimeZone(timeZone string) *DateHistogramAggregation {
	a.timeZone = timeZone
	return a
}

func (a *DateHistogramAggregation) Format(format string) *DateHistogramAggregation {
	a.format = format
	return a
}

func (a *DateHistogramAggregation) Offset(offset string) *DateHistogramAggregation {
	a.offset = offset
	return a
}

func (a *DateHistogramAggregation) Source() (interface{}, error) {
	// Example:
	// {
	//     "date_histogram" : {
	//         "field" : "date",
	//         "fixed_interval" : "month"
	//     }
	// }
	//
	// This method returns only the { "date_histogram" : { ... } } part.
	source := make(map[string]interface{})
	dateHistogram := make(map[string]interface{})
	source["date_histogram"] = dateHistogram
	dateHistogram["field"] = a.field
	if a.fixedInterval != "" {
		dateHistogram["fixed_interval"] = a.fixedInterval
	}
	if a.calendarInterval != "" {
		dateHistogram["calendar_interval"] = a.calendarInterval
	}
	if a.order != "" {
		dateHistogram["order"] = map[string]interface{}{
			a.order: a.ordering,
		}
	}
	if a.minDocCount != nil {
		dateHistogram["min_doc_count"] = *a.minDocCount
	}
	if a.timeZone != "" {
		dateHistogram["time_zone"] = a.timeZone
	}
	if a.format != "" {
		dateHistogram["format"] = a.format
	}
	if a.offset != "" {
		dateHistogram["offset"] = a.offset
	}

	return source, nil
}
