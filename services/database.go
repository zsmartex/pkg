package services

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase() (*gorm.DB, error) {
	var dialector gorm.Dialector

	dsn := "host=" + os.Getenv("DATABASE_HOST") +
		" port=" + os.Getenv("DATABASE_PORT") +
		" user=" + os.Getenv("DATABASE_USER") +
		" password=" + os.Getenv("DATABASE_PASS") +
		" dbname=" + os.Getenv("DATABASE_NAME") +
		" sslmode=disable"

	dialector = postgres.Open(dsn)

	db, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}
