package either

type tag uint8

const (
	left tag = iota
	right
)

type Either[E any, A any] struct {
	tag   tag
	left  *E
	right *A
}

func (e Either[E, A]) IsLeft() bool {
	return e.tag == left
}

func (e Either[E, A]) IsRight() bool {
	return e.tag == right
}

func Left[A any, E any](e E) Either[E, A] {
	return Either[E, A]{tag: left, left: &e}
}

func Right[E any, A any](a A) Either[E, A] {
	return Either[E, A]{tag: right, right: &a}
}

type optionable[A any] interface {
	Unwrap() (error, *A)
}

func FromOption[A any, E any](onNone func() E) func(optionable[A]) Either[E, A] {
	return func(o optionable[A]) Either[E, A] {
		err, a := o.Unwrap()

		if err != nil {
			return Left[A](onNone())
		}

		return Right[E](*a)
	}
}

func Fold[E any, A any, B any](
	onLeft func(E) B,
	onRight func(A) B,
) func(Either[E, A]) B {
	return func(e Either[E, A]) B {
		if e.tag == left {
			return onLeft(*e.left)
		}

		return onRight(*e.right)
	}
}

func (e Either[E, A]) GetOrElse(onLeft func(E) A) A {
	if e.tag == left {
		return onLeft(*e.left)
	}

	return *e.right
}

func (e Either[E, A]) GetOrElseValue(orElse A) A {
	if e.tag == left {
		return orElse
	}

	return *e.right
}

func (e Either[E, A]) Unwrap() (*E, *A) {
	return e.left, e.right
}

func Chain[E any, A any, B any](
	fab func(A) Either[E, B],
) func(Either[E, A]) Either[E, B] {
	return func(e Either[E, A]) Either[E, B] {
		if e.tag == left {
			return Either[E, B]{tag: left, left: e.left}
		}

		return fab(*e.right)
	}
}

func ChainI[E any, A any, B any](
	fab func(A) Either[E, B],
) func(any) Either[E, B] {
	chain := Chain(fab)
	return func(o any) Either[E, B] {
		return chain(o.(Either[E, A]))
	}
}

func Alt[E any, A any](
	f func(E) Either[E, A],
) func(Either[E, A]) Either[E, A] {
	return func(e Either[E, A]) Either[E, A] {
		if e.tag == right {
			return Either[E, A]{tag: right, right: e.right}
		}

		return f(*e.left)
	}
}

func AltI[E any, A any](
	f func(E) Either[E, A],
) func(any) Either[E, A] {
	fn := Alt(f)
	return func(o any) Either[E, A] {
		return fn(o.(Either[E, A]))
	}
}

func Map[E any, A any, B any](
	fab func(A) B,
) func(Either[E, A]) Either[E, B] {
	return func(e Either[E, A]) Either[E, B] {
		if e.tag == left {
			return Either[E, B]{tag: left, left: e.left}
		}

		b := fab(*e.right)
		return Either[E, B]{tag: right, right: &b}
	}
}

func MapI[E any, A any, B any](
	fab func(A) B,
) func(any) Either[E, B] {
	mapF := Map[E](fab)
	return func(o any) Either[E, B] {
		return mapF(o.(Either[E, A]))
	}
}

func MapLeft[E any, A any, E1 any](
	fab func(E) E1,
) func(Either[E, A]) Either[E1, A] {
	return func(e Either[E, A]) Either[E1, A] {
		if e.tag == right {
			return Either[E1, A]{tag: right, right: e.right}
		}

		a := fab(*e.left)
		return Either[E1, A]{tag: left, left: &a}
	}
}

func MapLeftI[E any, A any, E1 any](
	fab func(E) E1,
) func(any) Either[E1, A] {
	mapF := MapLeft[E, A](fab)
	return func(o any) Either[E1, A] {
		return mapF(o.(Either[E, A]))
	}
}

func Tap[E any, A any](
	f func(A) any,
) func(Either[E, A]) Either[E, A] {
	return func(e Either[E, A]) Either[E, A] {
		if e.tag == right {
			f(*e.right)
		}
		return e
	}
}

func TapI[E any, A any](
	f func(A) any,
) func(any) Either[E, A] {
	mapF := Tap[E](f)
	return func(o any) Either[E, A] {
		return mapF(o.(Either[E, A]))
	}
}

func Try[E any, A any](
	f func() (A, error),
	onError func(error) E,
) Either[E, A] {
	a, err := f()
	if err != nil {
		e := onError(err)
		return Either[E, A]{tag: left, left: &e}
	}

	return Either[E, A]{tag: right, right: &a}
}

func TryK[E any, A any, B any](
	f func(A) (B, error),
	onError func(error) E,
) func(A) Either[E, B] {
	return func(a A) Either[E, B] {
		return Try(
			func() (B, error) { return f(a) },
			onError,
		)
	}
}

func TryKI[E any, A any, B any](
	f func(A) (B, error),
	onError func(error) E,
) func(any) Either[E, B] {
	return func(a any) Either[E, B] {
		return Try(
			func() (B, error) { return f(a.(A)) },
			onError,
		)
	}
}

func ChainTryK[E any, A any, B any](
	f func(A) (B, error),
	onError func(error) E,
) func(Either[E, A]) Either[E, B] {
	return func(either Either[E, A]) Either[E, B] {
		if either.tag == left {
			return Either[E, B]{tag: left, left: either.left}
		}

		return Try(
			func() (B, error) { return f(*either.right) },
			onError,
		)
	}
}

func ChainTryKI[E any, A any, B any](
	f func(A) (B, error),
	onError func(error) E,
) func(any) Either[E, B] {
	fn := ChainTryK(f, onError)
	return func(either any) Either[E, B] {
		return fn(either.(Either[E, A]))
	}
}

func Filter[E any, A any](
	f func(A) bool,
	onFalse func(A) E,
) func(Either[E, A]) Either[E, A] {
	return func(either Either[E, A]) Either[E, A] {
		if either.tag == left {
			return either
		} else if f(*either.right) == true {
			return either
		}

		e := onFalse(*either.right)
		return Either[E, A]{tag: left, left: &e}
	}
}

func FilterI[E any, A any](
	f func(A) bool,
	onFalse func(A) E,
) func(any) Either[E, A] {
	fn := Filter(f, onFalse)
	return func(either any) Either[E, A] {
		return fn(either.(Either[E, A]))
	}
}
