package either

type tag uint8

const (
	left tag = iota
	right
)

type Either[E any, A any] struct {
	tag   tag
	left  E
	right A
}

func IsLeft[E any, A any](e Either[E, A]) bool {
	return e.tag == left
}

func IsRight[E any, A any](e Either[E, A]) bool {
	return e.tag == right
}

func Left[E any, A any](e E) Either[E, A] {
	return Either[E, A]{tag: left, left: e}
}

func Right[E any, A any](a A) Either[E, A] {
	return Either[E, A]{tag: right, right: a}
}

func Fold[E any, A any, B any](
	e Either[E, A],
	onLeft func(E) B,
	onRight func(A) B,
) B {
	if e.tag == left {
		return onLeft(e.left)
	}

	return onRight(e.right)
}

func Unwrap[E any, A any](e Either[E, A]) (E, A) {
	return e.left, e.right
}

func Chain[E any, A any, B any](
	fab func(A) Either[E, B],
) func(Either[E, A]) Either[E, B] {
	return func(e Either[E, A]) Either[E, B] {
		if e.tag == left {
			return Either[E, B]{tag: left, left: e.left}
		}

		return fab(e.right)
	}
}

func ChainI[E any, A any, B any](
	fab func(A) Either[E, B],
) func(interface{}) Either[E, B] {
	chain := Chain(fab)
	return func(o interface{}) Either[E, B] {
		return chain(o.(Either[E, A]))
	}
}
