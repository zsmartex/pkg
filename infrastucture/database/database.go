package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/zsmartex/pkg/v3/infrastucture/event_api"
	"github.com/zsmartex/pkg/v3/infrastucture/kafka"
	"github.com/zsmartex/pkg/v3/infrastucture/rango"
	"github.com/zsmartex/pkg/v3/log"

	_ "github.com/lib/pq"
)

type CallbackConfig struct {
	Producer *kafka.Producer
	Rango    *rango.Client
	EventAPI *event_api.EventAPI
}

type Config struct {
	Host                 string
	Port                 int
	User                 string
	Password             string
	DBName               string
	PreferSimpleProtocol bool
	Callback             *CallbackConfig
}

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

func New(config *Config) (*gorm.DB, error) {
	dialector := postgres.New(postgres.Config{
		DSN:                  fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.User, config.Password, config.DBName),
		PreferSimpleProtocol: config.PreferSimpleProtocol,
	})

	db, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction:   true,
		DisableNestedTransaction: true,
		Logger: &DBlogger{
			SlowThreshold:         200 * time.Millisecond,
			SkipErrRecordNotFound: true,
		},
	})

	if err != nil {
		return nil, err
	}

	if config.Callback != nil {
		if config.Callback.Producer != nil && config.Callback.Rango != nil {
			db.Callback().Update().Register("model:updated", func(db *gorm.DB) {
				if db.Statement.Schema != nil {
					if methodValue := db.Statement.ReflectValue.MethodByName("CustomAfterUpdate"); methodValue.IsValid() {
						switch methodValue.Type().String() {
						case "func(context.Context, *kafka.Producer, *rango.Client) error":
							methodValue.Call([]reflect.Value{reflect.ValueOf(db.Statement.Context), reflect.ValueOf(config.Callback.Producer), reflect.ValueOf(config.Callback.Rango)})
						default:
							log.Warnf("Model %v don't match %v Interface, should be `%v(context.Context, *kafka.Producer, *rango.RangoClient) error`. Please see https://gorm.io/docs/hooks.html", db.Statement.Schema, db.Statement.Schema.Name, db.Statement.Schema.Name)
						}
					}
				}
			})

			db.Callback().Create().Register("model:created", func(db *gorm.DB) {
				if db.Statement.Schema != nil {
					if methodValue := db.Statement.ReflectValue.MethodByName("CustomAfterCreate"); methodValue.IsValid() {
						switch methodValue.Type().String() {
						case "func(context.Context, *kafka.Producer, *rango.Client) error":
							methodValue.Call([]reflect.Value{reflect.ValueOf(db.Statement.Context), reflect.ValueOf(config.Callback.Producer), reflect.ValueOf(config.Callback.Rango)})
						default:
							log.Warnf("Model %v don't match %v Interface, should be `%v(context.Context, *kafka.Producer, *rango.RangoClient) error`. Please see https://gorm.io/docs/hooks.html", db.Statement.Schema, db.Statement.Schema.Name, db.Statement.Schema.Name)
						}
					}
				}
			})
		}

		if config.Callback.EventAPI != nil {
			db.Callback().Create().After("gorm:create").Register("eventapi:created", func(db *gorm.DB) {
				if db.Statement.Schema != nil {
					methodTableNameValue := db.Statement.ReflectValue.MethodByName("TableName")

					if !methodTableNameValue.IsValid() {
						log.Errorf("Failed to find TableName method in %v", db.Statement.Schema)
						return
					}

					if methodValue := db.Statement.ReflectValue.MethodByName("AsMapForEventAPI"); methodValue.IsValid() {
						switch methodValue.Type().String() {
						case "func() map[string]interface{}":
							tableName := methodTableNameValue.Call([]reflect.Value{})[0].Interface().(string)
							results := methodValue.Call([]reflect.Value{})

							eventName := fmt.Sprintf("model.%s.created", tableName)
							err := config.Callback.EventAPI.Notify(db.Statement.Context, eventName, event_api.EventAPIPayload{
								Record: results[0].Interface().(map[string]interface{}),
							})

							if err != nil {
								log.Errorf("Failed to send event to event api event_name: %s, err: %v", eventName, err)
							}
						default:
							log.Warnf("Model %v don't match %v Interface, should be `%v() map[string]interface{}`. Please see https://gorm.io/docs/hooks.html", db.Statement.Schema, db.Statement.Schema.Name, db.Statement.Schema.Name)
						}
					}
				}
			})

			db.Callback().Update().After("gorm:update").Register("eventapi:updated", func(db *gorm.DB) {
				if db.Statement.Schema != nil {
					methodTableNameValue := db.Statement.ReflectValue.MethodByName("TableName")

					if !methodTableNameValue.IsValid() {
						log.Errorf("Failed to find TableName method in %v", db.Statement.Schema)
						return
					}

					if methodValue := db.Statement.ReflectValue.MethodByName("AsMapForEventAPI"); methodValue.IsValid() {
						switch methodValue.Type().String() {
						case "func() map[string]interface{}":
							tableName := methodTableNameValue.Call([]reflect.Value{})[0].Interface().(string)
							results := methodValue.Call([]reflect.Value{})

							eventName := fmt.Sprintf("model.%s.updated", tableName)
							err := config.Callback.EventAPI.Notify(db.Statement.Context, eventName, event_api.EventAPIPayload{
								Record: results[0].Interface().(map[string]interface{}),
							})

							if err != nil {
								log.Errorf("Failed to send event to event api event_name: %s, err: %v", eventName, err)
							}
						default:
							log.Warnf("Model %v don't match %v Interface, should be `%v() map[string]interface{}`. Please see https://gorm.io/docs/hooks.html", db.Statement.Schema, db.Statement.Schema.Name, db.Statement.Schema.Name)
						}
					}
				}
			})
		}
	}

	return db, nil
}

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
