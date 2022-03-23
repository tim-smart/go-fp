package pipeline

func Identity[T any](a T) T {
	return a
}

type pineline[T any] struct {
	value T
}

type unsafePipeline[T any] struct {
	value interface{}
}

func Pipe[T any](value T) *pipeline[T] {
	return &pipeline[T]{
		value: value,
	}
}

func (p *pipeline[T]) Then(f func(T) T) *pipeline[T] {
	p.value = f(p.value)
	return p
}

func (p *pipeline[T]) ThenUnsafe(f func(T) interface{}) *unsafePipeline[T] {
	return &unsafePipeline{
		value: f(p.value),
	}
}

func (p *pipeline[T]) ThenSafe(f func(interface{}) T) *pipeline[T] {
	p.value = f(p.value)
	return p
}

func (p *pipeline[T]) Unsafe(f func(interface{}) interface{}) *unsafePipeline[T] {
	return &unsafePipeline{
		value: f(p.value),
	}
}

func (p *pipeline[T]) Result() T {
	return p.value
}

func PipeUnsafe[T any](value interface{}) *unsafePipeline[T] {
	return &unsafePipeline[T]{
		value: value,
	}
}

func (p *unsafePipeline[T]) ThenSafe(f func(interface{}) T) *pipeline[T] {
	return &pipeline{
		value: f(p.value),
	}
}

func (p *unsafePipeline[T]) Unsafe(f func(interface{}) interface{}) *unsafePipeline[T] {
	p.value = f(p.value)
	return p
}
