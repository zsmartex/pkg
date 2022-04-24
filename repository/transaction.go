package repository

import "gorm.io/gorm"

type ActionWithNonValue func(tx *gorm.DB, models interface{}) *gorm.DB
type ActionWithValue func(tx *gorm.DB, model interface{}, value interface{}) *gorm.DB

type Action interface {
	ActionWithNonValue | ActionWithValue
}

type Transaction func(tx *gorm.DB) error
type TransactionTrigger func(tx *gorm.DB, models interface{}) error
type TransactionOption func(optional *Optional)
type Optional struct {
	filters []Filter
}

func WithFilters(filters ...Filter) TransactionOption {
	return func(opt *Optional) {
		opt.filters = append(opt.filters, filters...)
	}
}

func MakeTransactionWithActionNonValue(action ActionWithNonValue, models interface{}, opts []TransactionOption) Transaction {
	var options Optional
	for _, opt := range opts {
		opt(&options)
	}
	return func(tx *gorm.DB) (err error) {
		for _, filter := range options.filters {
			tx = filter(tx)
		}
		if err = action(tx, models).Error; err != nil {
			return err
		}
		return
	}
}

func MakeTransactionWithActionValue(action ActionWithValue, model interface{}, value interface{}, opts []TransactionOption) Transaction {
	var options Optional
	for _, opt := range opts {
		opt(&options)
	}
	return func(tx *gorm.DB) (err error) {
		for _, filter := range options.filters {
			tx = filter(tx)
		}
		if err = action(tx, model, value).Error; err != nil {
			return err
		}
		return
	}
}
