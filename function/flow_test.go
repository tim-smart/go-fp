package function_test

import (
	"strconv"
	"strings"
	"testing"

	f "github.com/tim-smart/go-fp/function"
	o "github.com/tim-smart/go-fp/option"
)

var maybeParseString = f.FlowUnsafe[o.Option[string]](o.FromNilableI[string]).
	ThenSafe(o.MapI(strings.TrimSpace)).
	Then(o.Filter(func(a string) bool { return len(a) > 0 })).
	Result()

var maybeParseInt = f.FlowUnsafe[o.Option[int]](func(i *string) any {
	return maybeParseString(i)
}).
	ThenSafe(o.ChainTryKI(strconv.Atoi)).
	Result()

func TestFlow(t *testing.T) {
	if maybeParseInt(nil).IsSome() {
		t.Fail()
	}

	s := "jaklsdjf"
	if maybeParseInt(&s).IsSome() {
		t.Fail()
	}

	s = "123"
	if maybeParseInt(&s).GetOrElseValue(-1) != 123 {
		t.Fail()
	}
}
