package epa_test

// type Commission struct {
// 	ID        int64     `json:"id"`
// 	CreatedAt time.Time `json:"created_at"`
// }

// func (c Commission) IndexName() string {
// 	return "pg.commissions"
// }

// func newCommissionRepo() (epa.Repository[Commission], error) {
// 	client, err := elasticsearch_fx.New(config.Elasticsearch{
// 		URL:      []string{"http://demo.zsmartex.com:9200"},
// 		Username: "",
// 		Password: "",
// 	})
// 	if err != nil {
// 		return epa.Repository[Commission]{}, err
// 	}

// 	return epa.New(client, Commission{}), nil
// }

// func TestCommissionFind(t *testing.T) {
// 	repo, err := newCommissionRepo()
// 	if err != nil {
// 		panic(err)
// 	}

// 	result, err := repo.Find(
// 		context.Background(),
// 		epa.Query{
// 			Page:  1,
// 			Limit: 10,
// 			Aggregations: map[string]aggregation.Aggregation{
// 				"tong": aggregation.NewSumAggregation("total"),
// 			},
// 			Filters: []epa.Filter{},
// 		},
// 	)
// 	if err != nil {
// 		panic(err)
// 	}

// 	value, ecx := result.Aggregations.Sum("tong")
// 	log.Println(value.Value, ecx)
// }
