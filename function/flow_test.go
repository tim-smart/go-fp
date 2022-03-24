package function_test

import (
	"strconv"
	"testing"

	f "github.com/tim-smart/go-fp/function"
	o "github.com/tim-smart/go-fp/option"
)

func TestFlow(t *testing.T) {
	parseInt := f.FlowUnsafe[o.Option[int]](o.FromNilableI[string]).
		ThenSafe(o.ChainTryKI(strconv.Atoi)).
		Result()

	if parseInt(nil).IsSome() {
		t.Fail()
	}

	s := "jaklsdjf"
	if parseInt(&s).IsSome() {
		t.Fail()
	}

	s = "123"
	if parseInt(&s).GetOrElseValue(-1) != 123 {
		t.Fail()
	}
}
