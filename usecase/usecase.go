package usecase

import (
	"context"

	"github.com/cockroachdb/errors"
	"go.uber.org/fx"

	"github.com/zsmartex/pkg/v2/gpa"
	"github.com/zsmartex/pkg/v2/infrastructure/elasticsearch_fx"
	"github.com/zsmartex/pkg/v2/infrastructure/gorm_fx"
	"github.com/zsmartex/pkg/v2/infrastructure/questdb_fx"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	ErrBadConnection = errors.New("driver: bad connection")
)

var _ Usecase[schema.Tabler] = (*usecase[schema.Tabler])(nil)

type Usecase[V schema.Tabler] interface {
	AddCallback(kind gorm_fx.CallbackType, callback func(db *gorm.DB, value *V) error)
	WithTrx(tx *gorm.DB) Usecase[V]
	Count(ctx context.Context, filters ...gpa.Filter) (count int, err error)
	Last(ctx context.Context, filters ...gpa.Filter) (model *V, err error)
	First(ctx context.Context, filters ...gpa.Filter) (model *V, err error)
	Find(ctx context.Context, filters ...gpa.Filter) (models []*V, err error)
	Transaction(handler func(tx *gorm.DB) error) error
	FirstOrCreate(ctx context.Context, model *V, filters ...gpa.Filter) error
	Create(ctx context.Context, model *V, filters ...gpa.Filter) error
	Updates(ctx context.Context, model *V, value interface{}, filters ...gpa.Filter) error
	UpdateColumns(ctx context.Context, model *V, value interface{}, filters ...gpa.Filter) error
	Delete(ctx context.Context, model *V, filters ...gpa.Filter) error
	Exec(ctx context.Context, sql string, attrs ...interface{}) error
	RawFind(ctx context.Context, dst interface{}, sql string, attrs ...interface{}) error
	RawScan(ctx context.Context, dst interface{}, sql string, attrs ...interface{}) error
	RawFirst(ctx context.Context, dst interface{}, sql string, attrs ...interface{}) error

	Es() elasticsearch_fx.Repository[V]
	QuestDB() questdb_fx.Repository
}

type usecase[V schema.Tabler] struct {
	DatabaseRepo      gorm_fx.Repository[V]
	ElasticsearchRepo elasticsearch_fx.Repository[V]
	QuestDBRepo       questdb_fx.Repository
	Omits             []string
}

func Module[V schema.Tabler]() fx.Option {
	return fx.Options(
		fx.Provide(
			NewUsecase[V],
		),
	)
}

type Options[V schema.Tabler] struct {
	fx.In

	DatabaseRepo      gorm_fx.Repository[V]
	ElasticsearchRepo elasticsearch_fx.Repository[V] `optional:"true"`
	QuestDBRepo       questdb_fx.Repository          `optional:"true"`
	Omits             []string                       `name:"omits" optional:"true"`
}

func NewUsecase[V schema.Tabler](opts Options[V]) Usecase[V] {
	return &usecase[V]{
		DatabaseRepo:      opts.DatabaseRepo,
		ElasticsearchRepo: opts.ElasticsearchRepo,
		QuestDBRepo:       opts.QuestDBRepo,
		Omits:             opts.Omits,
	}
}
