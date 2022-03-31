package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBlogger struct {
	SlowThreshold         time.Duration
	SkipErrRecordNotFound bool
}

func (d *DBlogger) LogMode(logger.LogLevel) logger.Interface {
	return d
}

func (d *DBlogger) Info(ctx context.Context, s string, args ...interface{}) {
	log.WithContext(ctx).Infof(s, args)
}

func (d *DBlogger) Warn(ctx context.Context, s string, args ...interface{}) {
	log.WithContext(ctx).Warnf(s, args)
}

func (d *DBlogger) Error(ctx context.Context, s string, args ...interface{}) {
	log.WithContext(ctx).Errorf(s, args)
}

func (d *DBlogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, _ := fc()
	fields := log.Fields{}
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && d.SkipErrRecordNotFound) {
		fields[log.ErrorKey] = err
		log.WithContext(ctx).WithFields(fields).Errorf("%s [%s]", sql, elapsed)
		return
	}

	if d.SlowThreshold != 0 && elapsed > d.SlowThreshold {
		log.WithContext(ctx).WithFields(fields).Warnf("%s [%s]", sql, elapsed)
		return
	}

	log.WithContext(ctx).WithFields(fields).Debugf("%s [%s]", sql, elapsed)
}

func NewDatabase(host string, port int, user, password, dbname string) (*gorm.DB, error) {
	var dialector gorm.Dialector

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	dialector = postgres.Open(dsn)

	db, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
		Logger: &DBlogger{
			SlowThreshold:         200 * time.Millisecond,
			SkipErrRecordNotFound: true,
		},
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}
