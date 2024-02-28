package gorm_fx

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/zsmartex/pkg/v2/config"
	"github.com/zsmartex/pkg/v2/infrastructure/event_api_fx"
	"github.com/zsmartex/pkg/v2/log"
)

type Config struct {
	Host                 string
	Port                 int
	User                 string
	Password             string
	Name                 string
	PreferSimpleProtocol bool
}

type gormParams struct {
	fx.In

	MaxIdleConns    int           `name:"max_idle_conns"`
	MaxOpenConns    int           `name:"max_open_conns"`
	ConnMaxLifetime time.Duration `name:"conn_max_lifetime"`
	Config          config.Postgres
	EventAPI        *event_api_fx.EventAPI `optional:"true"`
}

func New(params gormParams) (*gorm.DB, error) {
	dialector := postgres.New(postgres.Config{
		DSN: fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			params.Config.Host,
			params.Config.Port,
			params.Config.User,
			params.Config.Pass,
			params.Config.Name,
		),
	})

	db, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction:   false,
		DisableNestedTransaction: false,
		Logger: &DBlogger{
			SlowThreshold:         15000 * time.Millisecond,
			SkipErrRecordNotFound: true,
		},
	})
	if err != nil {
		return nil, err
	}

	if params.EventAPI != nil {
		db.Callback().Create().After("gorm:create").Register("eventapi:created", func(db *gorm.DB) {
			if db.Statement.Schema != nil {
				methodTableNameValue := db.Statement.ReflectValue.MethodByName("TableName")

				if !methodTableNameValue.IsValid() {
					log.Errorf("Failed to find TableName method in %v", db.Statement.Schema)
					return
				}

				if methodValue := db.Statement.ReflectValue.MethodByName("ToEventAPI"); methodValue.IsValid() {
					switch methodValue.Type().String() {
					case "func() map[string]interface{}":
						tableName := methodTableNameValue.Call([]reflect.Value{})[0].Interface().(string)
						results := methodValue.Call([]reflect.Value{})

						eventName := fmt.Sprintf("model.%s.created", tableName)
						err := params.EventAPI.Notify(db.Statement.Context, eventName, event_api_fx.EventAPIPayload{
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
						err := params.EventAPI.Notify(db.Statement.Context, eventName, event_api_fx.EventAPIPayload{
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
