package function_test

import (
	"testing"

	"github.com/tim-smart/go-fp/function"
)

func TestMap_(t *testing.T) {
	i := function.Pipe(1).Then(func(i int) int {
		return i + 1
	}).Then(func(i int) int {
		return i * 2
	}).Result()

	if i != 4 {
		t.Fail()
	}
}
