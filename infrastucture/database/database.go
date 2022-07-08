package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"database/sql"

	_ "github.com/lib/pq"
)

type DBlogger struct {
	SlowThreshold         time.Duration
	SkipErrRecordNotFound bool
}

func (d *DBlogger) LogMode(logger.LogLevel) logger.Interface {
	return d
}

func (d *DBlogger) Info(ctx context.Context, s string, args ...interface{}) {
	logrus.WithContext(ctx).Infof(s, args)
}

func (d *DBlogger) Warn(ctx context.Context, s string, args ...interface{}) {
	logrus.WithContext(ctx).Warnf(s, args)
}

func (d *DBlogger) Error(ctx context.Context, s string, args ...interface{}) {
	logrus.WithContext(ctx).Errorf(s, args)
}

func (d *DBlogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, _ := fc()
	fields := logrus.Fields{}
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && d.SkipErrRecordNotFound) {
		fields[logrus.ErrorKey] = err
		logrus.WithContext(ctx).WithFields(fields).Errorf("%s [%s]", sql, elapsed)
		return
	}

	if d.SlowThreshold != 0 && elapsed > d.SlowThreshold {
		logrus.WithContext(ctx).WithFields(fields).Warnf("%s [%s]", sql, elapsed)
		return
	}

	logrus.WithContext(ctx).WithFields(fields).Debugf("%s [%s]", sql, elapsed)
}

func New(host string, port int, user, password, dbname string) (*gorm.DB, error) {
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

	// TODO: add support EventAPI here

	return db, nil
}

// create database if it doesn't exist
func CreateDatabase(host string, port int, user, password, dbname string) error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	if _, err := db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname)); err != nil {
		return err
	}

	return nil
}
