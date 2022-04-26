package gpa

import "gorm.io/gorm"

type Filter func(query *gorm.DB) *gorm.DB

func ApplyFilters(query *gorm.DB, filters []Filter) *gorm.DB {
	for _, f := range filters {
		query = f(query)
	}
	return query
}

func ChainFilters(filters ...Filter) Filter {
	return func(query *gorm.DB) *gorm.DB {
		for _, f := range filters {
			query = f(query)
		}
		return query
	}
}
