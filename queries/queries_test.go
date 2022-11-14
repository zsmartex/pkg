package queries

import (
	"testing"
)

func TestPeriodValidate(t *testing.T) {
	// Arrange
	period := &Period{}
	period.Init(3)
	// Act
	err := period.Validate(3, true)
	if err != nil {
		t.Fatal(err)
	}
}
