package option_test

import (
	"testing"

	"github.com/tim-smart/go-fp/option"
)

func TestMap(t *testing.T) {
	r1 := option.Some(1)
	r2 := option.Map(func(a int) string {
		return "asdc"
	})(r1)
	err, value := option.Unwrap(r2)

	if err != nil {
		t.Error("option not nil")
	}

	if value != "asdc" {
		t.Error("value not asdc")
	}
}
