package repository

import "gorm.io/gorm"

var (
	Create Action = func(tx *gorm.DB, models interface{}) *gorm.DB {
		return tx.Create(models)
	}

	Save Action = func(tx *gorm.DB, models interface{}) *gorm.DB {
		return tx.Save(models)
	}

	Updates Action = func(tx *gorm.DB, models interface{}) *gorm.DB {
		return tx.Updates(models)
	}

	Delete Action = func(tx *gorm.DB, models interface{}) *gorm.DB {
		return tx.Delete(models)
	}
)
