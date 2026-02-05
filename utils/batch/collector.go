// Package batch provides utilities for batching streaming data.
// This is used to efficiently batch logs and other events for
// processing in groups rather than individually.
package batch

import (
	"context"
	"sync"
	"time"
)

// Collector batches items from a channel into groups.
//
// Items are collected until either:
//   - The batch size is reached
//   - The timeout expires
//   - The input channel is closed
//
// This is useful for:
//   - Batching logs for efficient processing
//   - Reducing callback overhead when watching events
//   - Aggregating updates for batch database inserts
//
// Example:
//
//	collector := batch.NewCollector[Log](batch.CollectorOptions{
//	    BatchSize: 100,
//	    Timeout:   time.Second,
//	})
//
//	batches := collector.Collect(ctx, logChannel)
//	for batch := range batches {
//	    processBatch(batch)
//	}
type Collector[T any] struct {
	opts CollectorOptions
}

// CollectorOptions configures the batch collector.
type CollectorOptions struct {
	// BatchSize is the maximum number of items per batch.
	// When this many items are collected, the batch is emitted immediately.
	// Default: 100
	BatchSize int

	// Timeout is the maximum time to wait before emitting a partial batch.
	// If no new items arrive within this duration, the current batch is emitted.
	// Default: 1 second
	Timeout time.Duration

	// MinBatchSize is the minimum number of items to collect before emitting.
	// If set, batches smaller than this won't be emitted (except on timeout or close).
	// Default: 1
	MinBatchSize int
}

// DefaultCollectorOptions returns the default collector options.
func DefaultCollectorOptions() CollectorOptions {
	return CollectorOptions{
		BatchSize:    100,
		Timeout:      time.Second,
		MinBatchSize: 1,
	}
}

// NewCollector creates a new batch collector.
func NewCollector[T any](opts ...CollectorOptions) *Collector[T] {
	opt := DefaultCollectorOptions()
	if len(opts) > 0 {
		opt = opts[0]
	}
	if opt.BatchSize <= 0 {
		opt.BatchSize = 100
	}
	if opt.Timeout <= 0 {
		opt.Timeout = time.Second
	}
	if opt.MinBatchSize <= 0 {
		opt.MinBatchSize = 1
	}
	return &Collector[T]{opts: opt}
}

// Collect reads items from the input channel and batches them.
// The returned channel emits batches of items.
//
// The output channel is closed when:
//   - The input channel is closed (after emitting any remaining items)
//   - The context is canceled
func (c *Collector[T]) Collect(ctx context.Context, input <-chan T) <-chan []T {
	output := make(chan []T, 1)

	go func() {
		defer close(output)

		batch := make([]T, 0, c.opts.BatchSize)
		timer := time.NewTimer(c.opts.Timeout)
		defer timer.Stop()

		emit := func() {
			if len(batch) >= c.opts.MinBatchSize {
				select {
				case output <- batch:
				case <-ctx.Done():
					return
				}
			}
			batch = make([]T, 0, c.opts.BatchSize)
			timer.Reset(c.opts.Timeout)
		}

		for {
			select {
			case <-ctx.Done():
				// Emit remaining items before exit
				if len(batch) > 0 {
					select {
					case output <- batch:
					default:
					}
				}
				return

			case item, ok := <-input:
				if !ok {
					// Input closed - emit remaining and exit
					if len(batch) > 0 {
						select {
						case output <- batch:
						case <-ctx.Done():
						}
					}
					return
				}

				batch = append(batch, item)
				if len(batch) >= c.opts.BatchSize {
					emit()
				}

			case <-timer.C:
				// Timeout - emit partial batch
				if len(batch) > 0 {
					emit()
				} else {
					timer.Reset(c.opts.Timeout)
				}
			}
		}
	}()

	return output
}

// CollectWithFlush is like Collect but also provides a flush channel.
// Sending on the flush channel will immediately emit the current batch.
func (c *Collector[T]) CollectWithFlush(
	ctx context.Context,
	input <-chan T,
	flush <-chan struct{},
) <-chan []T {
	output := make(chan []T, 1)

	go func() {
		defer close(output)

		batch := make([]T, 0, c.opts.BatchSize)
		timer := time.NewTimer(c.opts.Timeout)
		defer timer.Stop()

		emit := func() {
			if len(batch) > 0 {
				select {
				case output <- batch:
				case <-ctx.Done():
					return
				}
				batch = make([]T, 0, c.opts.BatchSize)
			}
			timer.Reset(c.opts.Timeout)
		}

		for {
			select {
			case <-ctx.Done():
				if len(batch) > 0 {
					select {
					case output <- batch:
					default:
					}
				}
				return

			case item, ok := <-input:
				if !ok {
					if len(batch) > 0 {
						select {
						case output <- batch:
						case <-ctx.Done():
						}
					}
					return
				}

				batch = append(batch, item)
				if len(batch) >= c.opts.BatchSize {
					emit()
				}

			case <-timer.C:
				emit()

			case <-flush:
				emit()
			}
		}
	}()

	return output
}

// RingBuffer is a fixed-size circular buffer that drops oldest items when full.
// This is useful for scenarios where you want to keep only recent items
// without blocking the producer.
type RingBuffer[T any] struct {
	mu       sync.RWMutex
	items    []T
	size     int
	writeIdx int
	count    int
}

// NewRingBuffer creates a new ring buffer with the given capacity.
func NewRingBuffer[T any](capacity int) *RingBuffer[T] {
	return &RingBuffer[T]{
		items: make([]T, capacity),
		size:  capacity,
	}
}

// Push adds an item to the buffer, dropping the oldest item if full.
func (rb *RingBuffer[T]) Push(item T) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	rb.items[rb.writeIdx] = item
	rb.writeIdx = (rb.writeIdx + 1) % rb.size

	if rb.count < rb.size {
		rb.count++
	}
}

// Drain removes and returns all items from the buffer in order (oldest first).
func (rb *RingBuffer[T]) Drain() []T {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if rb.count == 0 {
		return nil
	}

	result := make([]T, rb.count)

	// Calculate start index
	startIdx := 0
	if rb.count == rb.size {
		startIdx = rb.writeIdx
	}

	for i := 0; i < rb.count; i++ {
		result[i] = rb.items[(startIdx+i)%rb.size]
	}

	// Reset buffer
	rb.writeIdx = 0
	rb.count = 0

	return result
}

// Peek returns all items without removing them (oldest first).
func (rb *RingBuffer[T]) Peek() []T {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	if rb.count == 0 {
		return nil
	}

	result := make([]T, rb.count)

	startIdx := 0
	if rb.count == rb.size {
		startIdx = rb.writeIdx
	}

	for i := 0; i < rb.count; i++ {
		result[i] = rb.items[(startIdx+i)%rb.size]
	}

	return result
}

// Len returns the current number of items in the buffer.
func (rb *RingBuffer[T]) Len() int {
	rb.mu.RLock()
	defer rb.mu.RUnlock()
	return rb.count
}

// Cap returns the capacity of the buffer.
func (rb *RingBuffer[T]) Cap() int {
	return rb.size
}

// Clear removes all items from the buffer.
func (rb *RingBuffer[T]) Clear() {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	rb.writeIdx = 0
	rb.count = 0
}
