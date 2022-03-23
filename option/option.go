package option

import "errors"

type tag uint8

const (
	some tag = iota
	none
)

type Option[T any] struct {
	tag   tag
	value T
}

func FromNilable[A any](a *A) Option[A] {
	if a == nil {
		return Option[A]{
			tag: none,
		}
	}

	return Option[A]{
		tag:   some,
		value: *a,
	}
}

func IsSome[T any](o Option[T]) bool {
	return o.tag == some
}

func IsNone[T any](o Option[T]) bool {
	return o.tag == none
}

func Some[T any](a T) Option[T] {
	return Option[T]{tag: some, value: a}
}

func None[T any]() Option[T] {
	return Option[T]{tag: none}
}

func Fold[A any, B any](
	o Option[A],
	onNone func() B,
	onSome func(A) B,
) B {
	if o.tag == none {
		return onNone()
	}

	return onSome(o.value)
}

func Unwrap[A any](o Option[A]) (error, A) {
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

		return Some(f(oa.value))
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
			return None[B]()
		}

		return f(o.value)
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
			return Some(*b)
		}
		return None[B]()
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
