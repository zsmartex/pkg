package triggers

import (
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"

	"github.com/zsmartex/pkg/repository"
)

func Mapping(target interface{}) repository.TransactionTrigger {
	return func(tx *gorm.DB, models interface{}) error {
		return mapstructure.Decode(models, target)
	}
}
