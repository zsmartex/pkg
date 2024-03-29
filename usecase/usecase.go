package usecase

import (
	"context"

	"github.com/cockroachdb/errors"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/zsmartex/pkg/v2/gpa"
	"github.com/zsmartex/pkg/v2/infrastructure/elasticsearch_fx"
	"github.com/zsmartex/pkg/v2/infrastructure/gorm_fx"
	"github.com/zsmartex/pkg/v2/infrastructure/mongo_fx"
	"github.com/zsmartex/pkg/v2/infrastructure/questdb_fx"
)

var (
	ErrBadConnection = errors.New("driver: bad connection")
)

var _ IUsecase[schema.Tabler] = (*Usecase[schema.Tabler])(nil)

type IUsecase[V schema.Tabler] interface {
	Repository() gorm_fx.Repository[V]
	Count(ctx context.Context, filters ...gpa.Filter) (count int, err error)
	Last(ctx context.Context, filters ...gpa.Filter) (model *V, err error)
	First(ctx context.Context, filters ...gpa.Filter) (model *V, err error)
	Find(ctx context.Context, filters ...gpa.Filter) (models []*V, err error)
	FindInBatches(ctx context.Context, batch int, filters ...gpa.Filter) (models []*V, err error)
	Transaction(handler func(tx *gorm.DB) error) error
	FirstOrCreate(ctx context.Context, model *V, filters ...gpa.Filter) error
	CreateInBatches(ctx context.Context, models []*V, batchSize int, filters ...gpa.Filter) error
	Create(ctx context.Context, model *V, filters ...gpa.Filter) error
	UpdateInBatches(ctx context.Context, value interface{}, filters ...gpa.Filter) error
	Updates(ctx context.Context, model *V, value interface{}, filters ...gpa.Filter) error
	UpdateColumns(ctx context.Context, model *V, value interface{}, filters ...gpa.Filter) error
	Delete(ctx context.Context, model *V, filters ...gpa.Filter) error
	Exec(ctx context.Context, sql string, attrs ...interface{}) error
	RawFind(ctx context.Context, dst interface{}, sql string, attrs ...interface{}) error
	RawScan(ctx context.Context, dst interface{}, sql string, attrs ...interface{}) error
	RawFirst(ctx context.Context, dst interface{}, sql string, attrs ...interface{}) error

	MongoDBRead() mongo_fx.ReadRepository[V]
	MongoDBWrite() mongo_fx.WriteRepository[V]
	Es() elasticsearch_fx.Repository[V]
	QuestDB() questdb_fx.Repository[V]
}

type Usecase[V schema.Tabler] struct {
	DatabaseRepo      gorm_fx.Repository[V]
	MongoDBReadRepo   mongo_fx.ReadRepository[V]
	MongoDBWriteRepo  mongo_fx.WriteRepository[V]
	ElasticsearchRepo elasticsearch_fx.Repository[V]
	QuestDBRepo       questdb_fx.Repository[V]
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
	MongoDBReadRepo   mongo_fx.ReadRepository[V]     `optional:"true"`
	MongoDBWriteRepo  mongo_fx.WriteRepository[V]    `optional:"true"`
	ElasticsearchRepo elasticsearch_fx.Repository[V] `optional:"true"`
	QuestDBRepo       questdb_fx.Repository[V]       `optional:"true"`
}

func NewUsecase[V schema.Tabler](opts Options[V]) Usecase[V] {
	return Usecase[V]{
		DatabaseRepo:      opts.DatabaseRepo,
		MongoDBReadRepo:   opts.MongoDBReadRepo,
		MongoDBWriteRepo:  opts.MongoDBWriteRepo,
		ElasticsearchRepo: opts.ElasticsearchRepo,
		QuestDBRepo:       opts.QuestDBRepo,
	}
}
