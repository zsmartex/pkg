package services

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)

	db, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 newLogger,
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}
