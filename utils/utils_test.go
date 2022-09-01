package utils

import (
	"testing"

	"github.com/volatiletech/null/v9"
)

func TestCompareDiff(t *testing.T) {
	type Label struct {
		Key string
	}

	type Test struct {
		Name   string
		Age    int
		Level  null.Int
		Labels []*Label
	}

	user1 := &Test{
		Name: "John",
		Age:  20,
		Level: null.Int{
			Int: 1,
		},
		Labels: nil,
	}

	user2 := &Test{
		Name: "John",
		Age:  21,
		Level: null.Int{
			Int: 2,
		},
		Labels: []*Label{
			{
				Key: "label2",
			},
		},
	}

	if err := CompareDiff(user1, user2, user1); err != nil {
		t.Error(err)
	}

	t.Log(user1)
}
