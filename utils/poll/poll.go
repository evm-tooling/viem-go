// Package poll provides a generic polling utility for Go.
// This is the Go equivalent of viem's poll() utility, using channels
// and context for idiomatic Go concurrency.
package poll

import (
	"context"
	"time"
)

// Options configures the polling behavior.
type Options struct {
	// Interval is the time between poll attempts.
	// Required.
	Interval time.Duration

	// EmitOnBegin determines whether to emit immediately when polling starts.
	// If true, the polling function is called immediately before waiting for the first interval.
	// Default: false
	EmitOnBegin bool

	// InitialWaitTime optionally specifies a different initial wait time.
	// If zero, uses Interval.
	InitialWaitTime time.Duration
}

// Result wraps the result of a poll iteration.
type Result[T any] struct {
	// Value is the value returned from the poll function.
	Value T

	// Error is any error that occurred during polling.
	Error error
}

// Poll executes the given function at regular intervals and sends results to a channel.
// The channel is closed when the context is canceled.
//
// This is the Go equivalent of viem's poll() utility, leveraging Go's native
// concurrency features:
//   - Channels for streaming results (instead of callbacks)
//   - context.Context for cancellation (instead of unwatch function)
//   - Goroutines for non-blocking operation
//
// Example:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	results := poll.Poll(ctx, func(ctx context.Context) (uint64, error) {
//	    return getBlockNumber(ctx)
//	}, poll.Options{
//	    Interval:    time.Second,
//	    EmitOnBegin: true,
//	})
//
//	for result := range results {
//	    if result.Error != nil {
//	        log.Printf("error: %v", result.Error)
//	        continue
//	    }
//	    fmt.Printf("block number: %d\n", result.Value)
//	}
func Poll[T any](ctx context.Context, fn func(ctx context.Context) (T, error), opts Options) <-chan Result[T] {
	ch := make(chan Result[T], 1) // Buffered to prevent blocking on slow consumers

	go func() {
		defer close(ch)

		// Emit on begin if requested
		if opts.EmitOnBegin {
			value, err := fn(ctx)
			select {
			case ch <- Result[T]{Value: value, Error: err}:
			case <-ctx.Done():
				return
			}
		}

		// Determine initial wait time
		initialWait := opts.InitialWaitTime
		if initialWait == 0 {
			initialWait = opts.Interval
		}

		// Wait for initial interval
		select {
		case <-ctx.Done():
			return
		case <-time.After(initialWait):
		}

		// Start polling loop
		ticker := time.NewTicker(opts.Interval)
		defer ticker.Stop()

		for {
			// Execute the poll function
			value, err := fn(ctx)

			// Send result (non-blocking with context check)
			select {
			case ch <- Result[T]{Value: value, Error: err}:
			case <-ctx.Done():
				return
			}

			// Wait for next tick or cancellation
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}
		}
	}()

	return ch
}

// PollWithUnpoll is similar to Poll but provides an unpoll channel that can be used
// to stop polling from within the poll function itself.
//
// The poll function receives an unpoll channel. Sending to this channel
// will stop the polling loop.
//
// Example:
//
//	results := poll.PollWithUnpoll(ctx, func(ctx context.Context, unpoll chan<- struct{}) (int, error) {
//	    value := getValue()
//	    if value >= 100 {
//	        close(unpoll) // Stop polling
//	    }
//	    return value, nil
//	}, opts)
func PollWithUnpoll[T any](
	ctx context.Context,
	fn func(ctx context.Context, unpoll chan<- struct{}) (T, error),
	opts Options,
) <-chan Result[T] {
	ch := make(chan Result[T], 1)
	unpoll := make(chan struct{})

	go func() {
		defer close(ch)

		// Create a merged context that cancels on either ctx or unpoll
		innerCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		go func() {
			select {
			case <-unpoll:
				cancel()
			case <-innerCtx.Done():
			}
		}()

		// Emit on begin if requested
		if opts.EmitOnBegin {
			value, err := fn(innerCtx, unpoll)
			select {
			case ch <- Result[T]{Value: value, Error: err}:
			case <-innerCtx.Done():
				return
			}
		}

		// Determine initial wait time
		initialWait := opts.InitialWaitTime
		if initialWait == 0 {
			initialWait = opts.Interval
		}

		// Wait for initial interval
		select {
		case <-innerCtx.Done():
			return
		case <-time.After(initialWait):
		}

		// Start polling loop
		ticker := time.NewTicker(opts.Interval)
		defer ticker.Stop()

		for {
			// Execute the poll function
			value, err := fn(innerCtx, unpoll)

			// Send result
			select {
			case ch <- Result[T]{Value: value, Error: err}:
			case <-innerCtx.Done():
				return
			}

			// Wait for next tick or cancellation
			select {
			case <-innerCtx.Done():
				return
			case <-ticker.C:
			}
		}
	}()

	return ch
}

// PollUntil polls until the predicate returns true or context is canceled.
// Returns the final value when the predicate is satisfied.
//
// Example:
//
//	result := poll.PollUntil(ctx, func(ctx context.Context) (int, error) {
//	    return getCount()
//	}, func(value int) bool {
//	    return value >= 10
//	}, opts)
func PollUntil[T any](
	ctx context.Context,
	fn func(ctx context.Context) (T, error),
	predicate func(T) bool,
	opts Options,
) Result[T] {
	results := Poll(ctx, fn, opts)

	for result := range results {
		if result.Error != nil {
			return result
		}
		if predicate(result.Value) {
			return result
		}
	}

	// Context was canceled
	var zero T
	return Result[T]{Value: zero, Error: ctx.Err()}
}
