package epa

import (
	"context"
	"log"
	"time"

	"github.com/zsmartex/pkg/v3/epa/aggregation"
	"github.com/zsmartex/pkg/v3/infrastucture/elasticsearch"
)

type Commission struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

func (c Commission) IndexName() string {
	return "pg.commissions"
}

func newCommissionRepo() (Repository[Commission], error) {
	client, err := elasticsearch.New(&elasticsearch.Config{
		URL:      []string{"http://demo.zsmartex.com:9200"},
		Username: "",
		Password: "",
	})
	if err != nil {
		return Repository[Commission]{}, err
	}

	return New(client, Commission{}), nil
}

func TestCommissionFind() {
	repo, err := newCommissionRepo()
	if err != nil {
		panic(err)
	}

	result, err := repo.Find(
		context.Background(),
		Query{
			Page:  1,
			Limit: 10,
			Aggregations: map[string]aggregation.Aggregation{
				"tong": aggregation.NewSumAggregation("total"),
			},
			Filters: []Filter{},
		},
	)
	if err != nil {
		panic(err)
	}

	value, ecx := result.Aggregations.Sum("tong")
	log.Println(value.Value, ecx)
}

func main() {
	TestCommissionFind()
}
