package elasticsearch_fx_test

// type Order struct {
// 	ID        int64     `json:"id"`
// 	CreatedAt time.Time `json:"created_at"`
// }

// func (o Order) IndexName() string {
// 	return "pg.orders"
// }

// func newRepo() (epa.Repository[Order], error) {
// 	client, err := elasticsearch_fx.New(config.Elasticsearch{
// 		URL:      []string{"http://zsmartex.com:9200"},
// 		Username: "elastic",
// 		Password: "elastic",
// 	})
// 	if err != nil {
// 		return epa.Repository[Order]{}, err
// 	}

// 	return epa.New(client, Order{}), nil
// }

// func TestCreate(t *testing.T) {
//repo, err := newRepo()
//if err != nil {
//	t.Error(err)
//}
//
//orders := []*Order{
//	{
//		ID:        1,
//		CreatedAt: time.Now(),
//	},
//	{
//		ID:        2,
//		CreatedAt: time.Now(),
//	},
//	{
//		ID:        3,
//		CreatedAt: time.Now(),
//	},
//}
//
//for _, order := range orders {
//	r, err := repo.Create(context.Background(), fmt.Sprint(order.ID), order)
//	if err != nil {
//		t.Error(err)
//	}
//
//	t.Error(r)
//}
// }

// func TestFind(t *testing.T) {
// 	repo, err := newRepo()
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	result, err := repo.Find(
// 		context.Background(),
// 		epa.Query{
// 			Limit: 10,
// 			Filters: []epa.Filter{
// 				epa.WithCreatedAtBefore(time.Now()),
// 				epa.WithFieldEqual("price", 10),
// 			},
// 		},
// 	)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	t.Log(result.Values)
// 	t.Log(result.TotalHits)
// }
