package pipeline

import "context"

// Pipe is a middleware-like function that receives a generic payload and the next pipe in the chain.
type Pipe[T any] func(ctx context.Context, payload T, next func(context.Context, T) (T, error)) (T, error)

// Pipeline manages a sequence of pipes.
type Pipeline[T any] struct {
	pipes []Pipe[T]
}

// New creates a new pipeline.
func New[T any]() *Pipeline[T] {
	return &Pipeline[T]{}
}

// Send sets the initial payload (not strictly required if passing directly to Then, but mirrors Laravel).
// We'll just append pipes.
func (p *Pipeline[T]) Through(pipes ...Pipe[T]) *Pipeline[T] {
	p.pipes = append(p.pipes, pipes...)
	return p
}

// Then executes the pipeline with a given destination callback.
func (p *Pipeline[T]) Then(ctx context.Context, payload T, destination func(context.Context, T) (T, error)) (T, error) {
	// Build the chain from the inside out
	next := destination

	for i := len(p.pipes) - 1; i >= 0; i-- {
		pipe := p.pipes[i]
		currentNext := next
		next = func(c context.Context, t T) (T, error) {
			return pipe(c, t, currentNext)
		}
	}

	return next(ctx, payload)
}
