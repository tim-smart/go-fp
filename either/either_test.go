package either_test

import (
	"errors"
	"testing"

	e "github.com/tim-smart/go-fp/either"
	f "github.com/tim-smart/go-fp/function"
	"github.com/tim-smart/go-fp/option"
)

func TestTry(t *testing.T) {
	either := e.Try(
		func() (int, error) {
			return 123, nil
		},
		func(err error) string { return "fail" },
	)

	if either.IsLeft() {
		t.Fail()
	}

	_, value := either.Unwrap()
	if *value != 123 {
		t.Fail()
	}
}

func TestTryFail(t *testing.T) {
	either := e.Try(
		func() (int, error) { return -1, errors.New("asdf") },
		func(err error) string { return "fail" },
	)

	if either.IsRight() {
		t.Fail()
	}

	err, _ := either.Unwrap()
	if *err != "fail" {
		t.Fail()
	}
}

func TestChain(t *testing.T) {
	result := f.Pipe(e.Right[string](1)).
		Then(e.Chain(func(a int) e.Either[string, int] {
			return e.Left[int]("fail")
		})).
		Result().
		GetOrElse(func(s string) int {
			return 42
		})

	if result != 42 {
		t.Fail()
	}
}

func TestAlt(t *testing.T) {
	result := f.Pipe(e.Right[string](1)).
		Then(e.Chain(func(_ int) e.Either[string, int] {
			return e.Left[int]("fail")
		})).
		Then(e.Alt(func(_ string) e.Either[string, int] {
			return e.Right[string](10)
		})).
		Result().
		GetOrElse(func(s string) int {
			return 42
		})

	if result != 10 {
		t.Fail()
	}
}

func TestFromOption(t *testing.T) {
	o := option.Some(1)
	result := e.FromOption[int](func() string {
		return "was none"
	})(o).GetOrElseValue(-1)

	if result != 1 {
		t.Fail()
	}
}
