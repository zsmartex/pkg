package queries

import (
	"testing"
)

func TestPeriodValidate(t *testing.T) {
	period := &Period{}
	period.Init(3)

	err := period.Validate(3, true)
	if err != nil {
		t.Fatal(err)
		t.Fail()
	}
}
