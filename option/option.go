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

func FromNilable[T any](a *T) Option[T] {
	if a == nil {
		return Option[T]{
			tag: none,
		}
	}

	return Option[T]{
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

func Fold[T any, A any](
	onNone func() A,
	onSome func(T) A,
) func(Option[T]) A {
	return func(o Option[T]) A {
		if o.tag == none {
			return onNone()
		}

		return onSome(o.value)
	}
}

func Unwrap[A any](o Option[A]) (error, A) {
	if o.tag == none {
		return errors.New("Unwrap: got none"), o.value
	}

	return nil, o.value
}

func Chain_[A any, B any](
	oa Option[A],
	f func(A) Option[B],
) Option[B] {
	if oa.tag == none {
		return None[B]()
	}

	return f(oa.value)
}

func Chain[A any, B any](
	f func(A) Option[B],
) func(Option[A]) Option[B] {
	return func(o Option[A]) Option[B] {
		return Chain_(o, f)
	}
}

func Map_[A any, B any](
	oa Option[A],
	f func(A) B,
) Option[B] {
	if oa.tag == none {
		return None[B]()
	}

	return Some(f(oa.value))
}

func Map[A any, B any](
	f func(A) B,
) func(Option[A]) Option[B] {
	return func(o Option[A]) Option[B] {
		return Map_(o, f)
	}
}
