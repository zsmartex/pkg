package queries

import (
	"fmt"
	"time"

	"github.com/zsmartex/pkg/v2/gpa"
	"github.com/zsmartex/pkg/v2/gpa/filters"
)

// Pagination is the pagination query parameters.
type Pagination struct {
	Page  int `query:"page" validate:"int|uint" default:"1"`
	Limit int `query:"limit" validate:"int|max:1000" default:"100"`
}

func (p *Pagination) GetFilter() gpa.Filter {
	return filters.WithPagination(p.Page, p.Limit)
}

// Period is a query parameter for time period.
type Period struct {
	TimeFrom int64 `query:"time_from" validate:"int"`
	TimeTo   int64 `query:"time_to" validate:"int"`
}

// GetFilter return the filter for the query.
func (p *Period) GetFilter() gpa.Filter {
	fs := make([]gpa.Filter, 0)

	if p.TimeFrom > 0 {
		fs = append(fs, filters.WithCreatedAtAfter(time.Unix(p.TimeFrom, 0)))
	}

	if p.TimeTo > 0 {
		fs = append(fs, filters.WithCreateAtBefore(time.Unix(p.TimeTo, 0)))
	}

	return gpa.ChainFilters(
		fs...,
	)
}

// Validate validate the query parameters and return the error.
//
// parameters:
//   - limitMonths is the maximum number of months that can be queried.
//   - limitUntil is whether to limit the time_from parameter.
func (p *Period) Validate(limitMonths int, limitUntil bool) error {
	if p.TimeFrom > 0 && p.TimeTo > 0 && p.TimeFrom > p.TimeTo {
		return fmt.Errorf("time_from must be less than time_to")
	}

	// time to and time from must between in 3 months
	if limitMonths > 0 && time.Unix(p.TimeTo, 0).Sub(time.Unix(p.TimeFrom, 0)).Hours() > float64(limitMonths*30*24) {
		return fmt.Errorf("time_to and time_from must be less than 3 months")
	}

	if limitUntil {
		if time.Since(time.Unix(p.TimeTo, 0)) > time.Duration(24*30*limitMonths)*time.Hour {
			return fmt.Errorf("time_to must be less than %d months", limitMonths)
		}
	}

	return nil
}

type Ordering string

var (
	OrderingAsc  Ordering = "asc"
	OrderingDesc Ordering = "desc"
)

// Order is the order query parameters.
type Order struct {
	OrderBy  string   `query:"order_by" default:"created_at"`
	Ordering Ordering `query:"ordering" validate:"ordering" default:"asc"`
}

// GetFilter return the filter for the query.
func (o *Order) GetFilter() gpa.Filter {
	return filters.WithOrder(fmt.Sprintf("%s %s", o.OrderBy, o.Ordering))
}
