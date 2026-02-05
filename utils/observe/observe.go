// Package observe provides an observer pattern implementation for deduplicating
// subscriptions. This is the Go equivalent of viem's observe() utility.
//
// When multiple consumers want to observe the same data source (e.g., watching
// blocks from the same client), the observer ensures that only one underlying
// subscription exists, and fan-outs the data to all consumers.
package observe

import (
	"sync"
	"sync/atomic"
)

// callbackID is a unique identifier for each callback registration.
var callbackID atomic.Int64

// Observer manages shared subscriptions and fan-out to multiple listeners.
// It is generic over the event type T.
//
// Example:
//
//	observer := observe.New[BlockEvent]()
//
//	// First subscriber - creates the underlying subscription
//	ch1 := observer.Subscribe("blocks-client1", func() (<-chan BlockEvent, func()) {
//	    ch := make(chan BlockEvent)
//	    // ... setup subscription ...
//	    return ch, cleanup
//	})
//
//	// Second subscriber - reuses the existing subscription
//	ch2 := observer.Subscribe("blocks-client1", nil) // setupFn ignored
type Observer[T any] struct {
	mu        sync.RWMutex
	listeners map[string][]listenerEntry[T]
	sources   map[string]<-chan T
	cleanup   map[string]func()
}

// listenerEntry holds a listener's channel and its unique ID.
type listenerEntry[T any] struct {
	id int64
	ch chan T
}

// New creates a new Observer instance.
func New[T any]() *Observer[T] {
	return &Observer[T]{
		listeners: make(map[string][]listenerEntry[T]),
		sources:   make(map[string]<-chan T),
		cleanup:   make(map[string]func()),
	}
}

// Subscribe registers a listener for the given observer ID.
//
// If this is the first subscriber for the ID, setupFn is called to create
// the underlying data source. Subsequent subscribers share the same source.
//
// The returned channel receives all events from the source. To unsubscribe,
// call Unsubscribe with the same channel.
//
// Parameters:
//   - observerID: Unique identifier for the subscription (e.g., "blocks-clientUID")
//   - setupFn: Function to create the underlying source. Only called for first subscriber.
//     Returns the source channel and a cleanup function.
//
// Returns a channel that receives events. The channel is closed when
// all subscribers have unsubscribed.
func (o *Observer[T]) Subscribe(
	observerID string,
	setupFn func() (<-chan T, func()),
) <-chan T {
	o.mu.Lock()
	defer o.mu.Unlock()

	// Create a channel for this listener
	id := callbackID.Add(1)
	ch := make(chan T, 10) // Buffered to prevent blocking

	entry := listenerEntry[T]{id: id, ch: ch}

	// Check if we already have listeners for this ID
	if existing, ok := o.listeners[observerID]; ok {
		// Add to existing listeners
		o.listeners[observerID] = append(existing, entry)
		return ch
	}

	// First subscriber - setup the source
	o.listeners[observerID] = []listenerEntry[T]{entry}

	if setupFn != nil {
		sourceCh, cleanupFn := setupFn()
		o.sources[observerID] = sourceCh
		if cleanupFn != nil {
			o.cleanup[observerID] = cleanupFn
		}

		// Start fan-out goroutine
		go o.fanOut(observerID, sourceCh)
	}

	return ch
}

// fanOut distributes events from the source to all listeners.
func (o *Observer[T]) fanOut(observerID string, sourceCh <-chan T) {
	for event := range sourceCh {
		o.mu.RLock()
		listeners := o.listeners[observerID]
		// Copy to avoid holding lock during send
		listenersCopy := make([]listenerEntry[T], len(listeners))
		copy(listenersCopy, listeners)
		o.mu.RUnlock()

		// Send to all listeners (non-blocking)
		for _, entry := range listenersCopy {
			select {
			case entry.ch <- event:
			default:
				// Listener channel is full, skip
				// This prevents slow consumers from blocking others
			}
		}
	}

	// Source closed - clean up
	o.cleanupObserver(observerID)
}

// Unsubscribe removes a listener for the given observer ID.
// If this was the last listener, the cleanup function is called.
func (o *Observer[T]) Unsubscribe(observerID string, ch chan T) {
	o.mu.Lock()
	defer o.mu.Unlock()

	listeners, ok := o.listeners[observerID]
	if !ok {
		return
	}

	// Find and remove the listener
	for i, entry := range listeners {
		if entry.ch == ch {
			// Remove from slice
			o.listeners[observerID] = append(listeners[:i], listeners[i+1:]...)
			close(ch)
			break
		}
	}

	// If no more listeners, cleanup
	if len(o.listeners[observerID]) == 0 {
		o.cleanupObserverLocked(observerID)
	}
}

// cleanupObserver cleans up an observer (acquires lock).
func (o *Observer[T]) cleanupObserver(observerID string) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.cleanupObserverLocked(observerID)
}

// cleanupObserverLocked cleans up an observer (must hold lock).
func (o *Observer[T]) cleanupObserverLocked(observerID string) {
	// Call cleanup function if exists
	if cleanupFn, ok := o.cleanup[observerID]; ok {
		cleanupFn()
		delete(o.cleanup, observerID)
	}

	// Close all listener channels
	for _, entry := range o.listeners[observerID] {
		close(entry.ch)
	}

	// Remove from maps
	delete(o.listeners, observerID)
	delete(o.sources, observerID)
}

// UnsubscribeAll removes all listeners for the given observer ID
// and calls the cleanup function.
func (o *Observer[T]) UnsubscribeAll(observerID string) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.cleanupObserverLocked(observerID)
}

// ListenerCount returns the number of listeners for the given observer ID.
func (o *Observer[T]) ListenerCount(observerID string) int {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return len(o.listeners[observerID])
}

// HasListeners returns true if there are any listeners for the given observer ID.
func (o *Observer[T]) HasListeners(observerID string) bool {
	return o.ListenerCount(observerID) > 0
}

// Global observer instances for common use cases.
// These can be used when you want to share subscriptions across
// different parts of your application.
var (
	// BlockObserver is a global observer for block events.
	blockObserverOnce sync.Once
	blockObserver     *Observer[any]
)

// GetBlockObserver returns the global block observer.
// This is useful for sharing block subscriptions across the application.
func GetBlockObserver() *Observer[any] {
	blockObserverOnce.Do(func() {
		blockObserver = New[any]()
	})
	return blockObserver
}
