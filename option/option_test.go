package option_test

import (
	"testing"

	f "github.com/tim-smart/go-fp/function"
	o "github.com/tim-smart/go-fp/option"
)

func TestMap(t *testing.T) {
	result := f.
		PipeUnsafe[o.Option[string]](o.Some(1)).
		ThenSafe(o.MapI(func(a int) string { return "asdc" })).
		Result()

	err, value := o.Unwrap(result)

	if err != nil {
		t.Error("err not nil")
	}

	if *value != "asdc" {
		t.Error("value not asdc")
	}
}
