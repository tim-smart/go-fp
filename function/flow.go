package function

type flow[I any, O any] struct {
	fn func(I) O
}

type unsafeFlow[I any, O any] struct {
	fn func(I) any
}

// `Flow` starts a function composition pipeline.
// The `I` generic determines the final function input type, while `O`
// determines the return type.
//
// If you need to start a pipeline with a type that is different to the desired
// final result types, then you can use `FlowUnsafe`.
func Flow[I any, O any](fn func(I) O) *flow[I, O] {
	return &flow[I, O]{fn: fn}
}

// `Then` can be used to chain computations that have input and output types
// that match the pipeline's return type
func (f *flow[I, O]) Then(fb func(O) O) *flow[I, O] {
	fa := f.fn
	f.fn = func(i I) O { return fb(fa(i)) }
	return f
}

// `ThenUnsafe` is used when the chained function has a different return type to
// the pipeline's result type.
func (f *flow[I, O]) ThenUnsafe(fb func(O) interface{}) *unsafeFlow[I, O] {
	fa := f.fn
	return &unsafeFlow[I, O]{
		fn: func(i I) any {
			return fb(fa(i))
		},
	}
}

// `ThenSafe` is used when the chained function has the same return type as the
// pipeline's result type.
//
// It marks the pipeline as "safe" again, so the result can be accessed.
func (f *flow[I, O]) ThenSafe(fb func(any) O) *flow[I, O] {
	fa := f.fn
	f.fn = func(i I) O { return fb(fa(i)) }
	return f
}

// `Unsafe` is used when the chained function has both input out return types
// different to the pipeline's desired result type.
func (f *flow[I, O]) Unsafe(fb func(any) any) *unsafeFlow[I, O] {
	fa := f.fn
	return &unsafeFlow[I, O]{
		fn: func(i I) any { return fb(fa(i)) },
	}
}

// Access the result of the pipeline
func (f *flow[I, O]) Result() func(I) O {
	return f.fn
}

// `FlowUnsafe` starts a pipeline with a value different to the pipeline's
// desired result type.
//
// If you need to start a function composition pipeline with a type that is the
// same as the desired final result type, then you can use `Flow`.
func FlowUnsafe[O any, I any](fa func(I) any) *unsafeFlow[I, O] {
	return &unsafeFlow[I, O]{fn: fa}
}

// `ThenSafe` is used when the chained function has the same return type as the
// pipeline's result type.
//
// It marks the pipeline as "safe" again, so the result can be accessed.
func (f *unsafeFlow[I, O]) ThenSafe(fb func(any) O) *flow[I, O] {
	fa := f.fn
	return &flow[I, O]{
		fn: func(i I) O { return fb(fa(i)) },
	}
}

// `Unsafe` is used when the chained function has both input out return types
// different to the pipeline's desired result type.
func (f *unsafeFlow[I, O]) Unsafe(fb func(any) any) *unsafeFlow[I, O] {
	fa := f.fn
	f.fn = func(i I) any { return fb(fa(i)) }
	return f
}
