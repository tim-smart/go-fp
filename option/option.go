package option

import (
	"encoding/json"
	"errors"
)

type tag uint8

const (
	some tag = iota
	none
)

type Option[A any] struct {
	tag   tag
	value *A
}

var _ json.Marshaler = (*Option[int])(nil)
var _ json.Unmarshaler = (*Option[int])(nil)

func (o Option[A]) MarshalJSON() ([]byte, error) {
	if o.tag == none {
		return []byte("null"), nil
	}

	return json.Marshal(o.value)
}

func (o *Option[A]) UnmarshalJSON(bytes []byte) error {
	v := (*A)(nil)
	json.Unmarshal(bytes, &v)

	if v != nil {
		o.tag = some
		o.value = v
	} else {
		o.tag = none
		o.value = nil
	}

	return nil
}

func FromNilable[A any](a *A) Option[A] {
	if a == nil {
		return Option[A]{
			tag: none,
		}
	}

	return Option[A]{
		tag:   some,
		value: a,
	}
}

func FromNilableI[A any](a *A) any {
	return FromNilable(a)
}

type eitherable[E any, A any] interface {
	Unwrap() (*E, *A)
}

func FromEither[A any](either eitherable[any, A]) Option[A] {
	e, a := either.Unwrap()

	if e != nil {
		return None[A]()
	}

	return Some(*a)
}

func (o Option[A]) IsSome() bool {
	return o.tag == some
}

func (o Option[A]) IsNone() bool {
	return o.tag == none
}

func Some[A any](a A) Option[A] {
	return Option[A]{tag: some, value: &a}
}

func None[A any]() Option[A] {
	return Option[A]{tag: none}
}

func Fold[A any, B any](
	onNone func() B,
	onSome func(A) B,
) func(Option[A]) B {
	return func(o Option[A]) B {
		if o.tag == none {
			return onNone()
		}

		return onSome(*o.value)
	}
}

func (o Option[A]) GetOrElse(onNone func() A) A {
	if o.tag == none {
		return onNone()
	}

	return *o.value
}

func (o Option[A]) GetOrElseValue(orElse A) A {
	if o.tag == none {
		return orElse
	}

	return *o.value
}

func (o Option[A]) Unwrap() (error, *A) {
	if o.tag == none {
		return errors.New("Unwrap: got none"), o.value
	}

	return nil, o.value
}

func Map[A any, B any](
	f func(A) B,
) func(Option[A]) Option[B] {
	return func(oa Option[A]) Option[B] {
		if oa.tag == none {
			return None[B]()
		}

		return Some(f(*oa.value))
	}
}

func MapI[A any, B any](
	f func(A) B,
) func(interface{}) Option[B] {
	fMap := Map(f)
	return func(o interface{}) Option[B] {
		return fMap(o.(Option[A]))
	}
}

func Chain[A any, B any](
	f func(A) Option[B],
) func(Option[A]) Option[B] {
	return func(o Option[A]) Option[B] {
		if o.tag == none {
			return Option[B]{tag: none}
		}

		return f(*o.value)
	}
}

func ChainI[A any, B any](
	f func(A) Option[B],
) func(interface{}) Option[B] {
	chain := Chain(f)
	return func(o interface{}) Option[B] {
		return chain(o.(Option[A]))
	}
}

func ChainNilable[A any, B any](
	f func(A) *B,
) func(Option[A]) Option[B] {
	return Chain(func(a A) Option[B] {
		if b := f(a); b != nil {
			return Option[B]{tag: some, value: b}
		}
		return Option[B]{tag: none}
	})
}

func ChainNilableI[A any, B any](
	f func(A) *B,
) func(interface{}) Option[B] {
	chain := ChainNilable(f)
	return func(o interface{}) Option[B] {
		return chain(o.(Option[A]))
	}
}

func Filter[A any](f func(A) bool) func(Option[A]) Option[A] {
	return func(option Option[A]) Option[A] {
		if option.tag == none {
			return option
		} else if f(*option.value) == true {
			return option
		}

		return Option[A]{tag: none}
	}
}

func FilterI[A any](f func(A) bool) func(interface{}) Option[A] {
	fn := Filter(f)
	return func(option interface{}) Option[A] {
		return fn(option.(Option[A]))
	}
}

func Try[A any](f func() (A, error)) Option[A] {
	a, err := f()
	if err != nil {
		return Option[A]{tag: none}
	}
	return Option[A]{tag: some, value: &a}
}

func TryK[A any, B any](f func(A) (B, error)) func(A) Option[B] {
	return func(a A) Option[B] {
		return Try(func() (B, error) { return f(a) })
	}
}

func TryKI[A any, B any](f func(A) (B, error)) func(interface{}) Option[B] {
	return func(a interface{}) Option[B] {
		return Try(func() (B, error) { return f(a.(A)) })
	}
}

func ChainTryK[A any, B any](f func(A) (B, error)) func(Option[A]) Option[B] {
	return func(option Option[A]) Option[B] {
		if option.tag == none {
			return Option[B]{tag: none}
		}

		return Try(func() (B, error) { return f(*option.value) })
	}
}

func ChainTryKI[A any, B any](f func(A) (B, error)) func(interface{}) Option[B] {
	fn := ChainTryK(f)
	return func(either interface{}) Option[B] {
		return fn(either.(Option[A]))
	}
}
