package validate

import (
	"github.com/gookit/validate"
	"github.com/zsmartex/pkg/queries"
)

func InitValidation() {
	validate.AddValidator("ordering", func(val interface{}) bool {
		v := val.(queries.Ordering)

		return v == queries.OrderingAsc || v == queries.OrderingDesc
	})
}