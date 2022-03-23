package function

type pipeline[T any] struct {
	value T
}

type unsafePipeline[T any] struct {
	value interface{}
}

// `Pipe` starts a pipeline safely.
// The `T` generic determines the desired final result type.
//
// If you need to start a pipeline with a type that is different to the desired
// final result type, then you can use `PipeUnsafe`.
func Pipe[T any](value T) *pipeline[T] {
	return &pipeline[T]{
		value: value,
	}
}

// `Then` can be used to chain computations that have input and output types
// that match the pipeline's return type
func (p *pipeline[T]) Then(f func(T) T) *pipeline[T] {
	p.value = f(p.value)
	return p
}

// `ThenUnsafe` is used when the chained function has a different return type to
// the pipeline's result type.
func (p *pipeline[T]) ThenUnsafe(f func(T) interface{}) *unsafePipeline[T] {
	return &unsafePipeline[T]{
		value: f(p.value),
	}
}

// `ThenSafe` is used when the chained function has the same return type as the
// pipeline's result type.
//
// It marks the pipeline as "safe" again, so the result can be accessed.
func (p *pipeline[T]) ThenSafe(f func(interface{}) T) *pipeline[T] {
	p.value = f(p.value)
	return p
}

// `Unsafe` is used when the chained function has both input out return types
// different to the pipeline's desired result type.
func (p *pipeline[T]) Unsafe(f func(interface{}) interface{}) *unsafePipeline[T] {
	return &unsafePipeline[T]{
		value: f(p.value),
	}
}

// Access the result of the pipeline
func (p *pipeline[T]) Result() T {
	return p.value
}

// `PipeUnsafe` starts a pipeline with a value different to the pipeline's
// desired result type.
// The `T` generic determines the desired final result type.
//
// If you need to start a pipeline with a type that is the same as the desired
// final result type, then you can use `Pipe`.
func PipeUnsafe[T any](value interface{}) *unsafePipeline[T] {
	return &unsafePipeline[T]{
		value: value,
	}
}

// `ThenSafe` is used when the chained function has the same return type as the
// pipeline's result type.
//
// It marks the pipeline as "safe" again, so the result can be accessed.
func (p *unsafePipeline[T]) ThenSafe(f func(interface{}) T) *pipeline[T] {
	return &pipeline[T]{
		value: f(p.value),
	}
}

// `Unsafe` is used when the chained function has both input out return types
// different to the pipeline's desired result type.
func (p *unsafePipeline[T]) Unsafe(f func(interface{}) interface{}) *unsafePipeline[T] {
	p.value = f(p.value)
	return p
}
