package function_test

import (
	"fmt"
	"testing"

	f "github.com/tim-smart/go-fp/function"
	o "github.com/tim-smart/go-fp/option"
)

func TestPipe(t *testing.T) {
	i := f.PipeUnsafe[o.Option[string]](o.Some(1)).
		ThenSafe(o.MapI(func(a int) string {
			return fmt.Sprintf("got: %d", a)
		})).
		Then(o.Map(func(a string) string {
			return fmt.Sprintf("test: %s", a)
		})).
		Result()

	err, result := o.Unwrap(i)
	if err != nil {
		t.Fail()
	}
	if result != "test: got: 1" {
		t.Fail()
	}
}
