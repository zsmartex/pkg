package repository

import "gorm.io/gorm"

var (
	FirstOrCreate ActionWithNonValue = func(tx *gorm.DB, models interface{}) *gorm.DB {
		return tx.FirstOrCreate(models)
	}

	Create ActionWithNonValue = func(tx *gorm.DB, models interface{}) *gorm.DB {
		return tx.Create(models)
	}

	Save ActionWithNonValue = func(tx *gorm.DB, models interface{}) *gorm.DB {
		return tx.Save(models)
	}

	Updates ActionWithValue = func(tx *gorm.DB, model interface{}, value interface{}) *gorm.DB {
		return tx.Model(model).Updates(value)
	}

	UpdateColumns ActionWithValue = func(tx *gorm.DB, model interface{}, value interface{}) *gorm.DB {
		return tx.Model(model).UpdateColumns(value)
	}

	Delete ActionWithNonValue = func(tx *gorm.DB, models interface{}) *gorm.DB {
		return tx.Delete(models)
	}
)
