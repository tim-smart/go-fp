package option_test

import (
	"bytes"
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

func TestJsonMarshal(t *testing.T) {
	json, _ := o.Some(123).MarshalJSON()

	if !bytes.Equal(json, []byte("123")) {
		t.Errorf("got: %s", json)
	}
}

func TestJsonMarshalNone(t *testing.T) {
	json, _ := o.None[int]().MarshalJSON()

	if !bytes.Equal(json, []byte("null")) {
		t.Errorf("got: %s", json)
	}
}

func TestJsonUnmarshal(t *testing.T) {
	oa := o.None[int]()
	oa.UnmarshalJSON([]byte("123"))

	if o.IsNone(oa) {
		t.Fail()
	}
}

func TestJsonUnmarshalNull(t *testing.T) {
	oa := o.Some(123)
	oa.UnmarshalJSON([]byte("null"))

	if o.IsSome(oa) {
		t.Fail()
	}
}
