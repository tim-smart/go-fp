package function

func Identity[T any](a T) T {
	return a
}

type Pipeline[T any] struct {
	value T
}

func Pipe[T any](value T) *Pipeline[T] {
	return &Pipeline[T]{
		value: value,
	}
}

func (pipeline *Pipeline[T]) Then(f func(T) T) *Pipeline[T] {
	pipeline.value = f(pipeline.value)
	return pipeline
}

func (pipeline *Pipeline[T]) Result() T {
	return pipeline.value
}

type UnsafePipeline[T any] struct {
	value interface{}
}

func PipeUnsafe[T any](value interface{}) *UnsafePipeline[T] {
	return &UnsafePipeline[T]{
		value: value,
	}
}

func (pipeline *UnsafePipeline[T]) Then(f func(interface{}) interface{}) *UnsafePipeline[T] {
	pipeline.value = f(pipeline.value)
	return pipeline
}

func (pipeline *UnsafePipeline[T]) MakeSafe() *Pipeline[T] {
	return &Pipeline[T]{value: pipeline.value.(T)}
}

func (pipeline *UnsafePipeline[T]) Result() T {
	return pipeline.value.(T)
}
