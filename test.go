package pkg

import (
	"context"

	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {
	c, _ := kgo.NewClient()

	kgo.ConsumeTopics()

	fetches := c.PollFetches(context.Background())

	fetches.EachRecord(func(r *kgo.Record) {
		c.CommitRecords(context.Background(), r)
	})
}
