package gorm_fx

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/sirupsen/logrus"
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
	logrus.WithContext(ctx).Infof(s, args...)
}

func (d *DBlogger) Warn(ctx context.Context, s string, args ...interface{}) {
	logrus.WithContext(ctx).Warnf(s, args...)
}

func (d *DBlogger) Error(ctx context.Context, s string, args ...interface{}) {
	logrus.WithContext(ctx).Errorf(s, args...)
}

func (d *DBlogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sqlStr, _ := fc()
	space := regexp.MustCompile(`\s+`)
	sqlStr = space.ReplaceAllString(sqlStr, " ")
	fields := logrus.Fields{}
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && d.SkipErrRecordNotFound) {
		fields[logrus.ErrorKey] = err
		logrus.WithContext(ctx).WithFields(fields).Errorf("%s [%s]", sqlStr, elapsed)
		return
	}

	if d.SlowThreshold != 0 && elapsed > d.SlowThreshold {
		logrus.WithContext(ctx).WithFields(fields).Warnf("%s [%s]", sqlStr, elapsed)
		return
	}

	logrus.WithContext(ctx).WithFields(fields).Tracef("%s [%s]", sqlStr, elapsed)
}
