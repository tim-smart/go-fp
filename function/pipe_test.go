package function_test

import (
	"fmt"
	"testing"

	f "github.com/tim-smart/go-fp/function"
	o "github.com/tim-smart/go-fp/option"
)

func TestPipe(t *testing.T) {
	i := f.Pipe(o.Some(1)).
		Then(o.Map(func(a int) int {
			return a + 10
		})).
		Result()

	result := o.Fold(
		func() string { return "nothing" },
		func(a int) string { return fmt.Sprintf("got: %d", a) },
	)(i)

	if result != "got: 11" {
		t.Fail()
	}
}

func TestPipeUnsafe(t *testing.T) {
	result, err := f.PipeUnsafe[o.Option[string]](o.Some(1)).
		ThenSafe(o.MapI(func(a int) string {
			return fmt.Sprintf("got: %d", a)
		})).
		Then(o.Map(func(a string) string {
			return fmt.Sprintf("test: %s", a)
		})).
		Result().
		Unwrap()

	if err != nil {
		t.Fail()
	}
	if *result != "test: got: 1" {
		t.Fail()
	}
}
