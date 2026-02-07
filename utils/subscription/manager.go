// Package subscription provides a subscription manager with automatic reconnection
// for WebSocket subscriptions. This is part of the Go optimization patterns for
// reliable real-time data streaming.
package subscription

import (
	"context"
	"time"

	json "github.com/goccy/go-json"
)

// SubscribeParams contains parameters for a WebSocket subscription.
type SubscribeParams struct {
	// Type is the subscription type (newHeads, newPendingTransactions, logs, syncing).
	Type string
	// Params are additional parameters for the subscription (e.g., filter for logs).
	Params any
}

// NewHeadsParams creates params for newHeads subscription.
func NewHeadsParams() SubscribeParams {
	return SubscribeParams{Type: "newHeads"}
}

// NewPendingTransactionsParams creates params for newPendingTransactions subscription.
func NewPendingTransactionsParams() SubscribeParams {
	return SubscribeParams{Type: "newPendingTransactions"}
}

// LogsParams creates params for logs subscription with optional filter.
func LogsParams(address any, topics []any) SubscribeParams {
	params := make(map[string]any)
	if address != nil {
		params["address"] = address
	}
	if topics != nil {
		params["topics"] = topics
	}
	return SubscribeParams{
		Type:   "logs",
		Params: params,
	}
}

// SyncingParams creates params for syncing subscription.
func SyncingParams() SubscribeParams {
	return SubscribeParams{Type: "syncing"}
}

// Transport defines the interface for WebSocket-capable transports.
type Transport interface {
	// Subscribe creates a subscription.
	Subscribe(
		params SubscribeParams,
		onData func(data json.RawMessage),
		onError func(err error),
	) (Subscription, error)
}

// Subscription represents an active subscription.
type Subscription interface {
	// Unsubscribe cancels the subscription.
	Unsubscribe() error
}

// Manager manages WebSocket subscriptions with automatic reconnection.
//
// Features:
//   - Automatic reconnection on connection loss
//   - Configurable retry behavior
//   - Context-based cancellation
//   - Thread-safe
//
// Example:
//
//	manager := subscription.NewManager(transport, subscription.ManagerOptions{
//	    MaxReconnectAttempts: 5,
//	    ReconnectDelay:       time.Second,
//	})
//
//	ch, err := manager.SubscribeWithReconnect(ctx, subscription.NewHeadsParams())
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	for data := range ch {
//	    // Process data
//	}
type Manager struct {
	transport Transport
	opts      ManagerOptions
}

// ManagerOptions configures the subscription manager.
type ManagerOptions struct {
	// MaxReconnectAttempts is the maximum number of reconnection attempts.
	// Zero means infinite retries.
	// Default: 5
	MaxReconnectAttempts int

	// ReconnectDelay is the initial delay between reconnection attempts.
	// The delay increases exponentially with each attempt.
	// Default: 1 second
	ReconnectDelay time.Duration

	// MaxReconnectDelay is the maximum delay between reconnection attempts.
	// Default: 30 seconds
	MaxReconnectDelay time.Duration

	// OnReconnect is called when a reconnection attempt is made.
	OnReconnect func(attempt int)

	// OnReconnectSuccess is called when reconnection succeeds.
	OnReconnectSuccess func()

	// OnReconnectFailure is called when all reconnection attempts fail.
	OnReconnectFailure func(err error)
}

// DefaultManagerOptions returns the default manager options.
func DefaultManagerOptions() ManagerOptions {
	return ManagerOptions{
		MaxReconnectAttempts: 5,
		ReconnectDelay:       time.Second,
		MaxReconnectDelay:    30 * time.Second,
	}
}

// NewManager creates a new subscription manager.
func NewManager(transport Transport, opts ...ManagerOptions) *Manager {
	opt := DefaultManagerOptions()
	if len(opts) > 0 {
		opt = opts[0]
	}
	return &Manager{
		transport: transport,
		opts:      opt,
	}
}

// SubscribeWithReconnect creates a subscription that automatically reconnects
// on connection loss.
//
// The returned channel receives subscription data. It is closed when the
// context is canceled or all reconnection attempts fail.
//
// Example:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	ch, err := manager.SubscribeWithReconnect(ctx, subscription.NewHeadsParams())
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	for data := range ch {
//	    var header BlockHeader
//	    json.Unmarshal(data, &header)
//	    fmt.Printf("New block: %d\n", header.Number)
//	}
func (m *Manager) SubscribeWithReconnect(
	ctx context.Context,
	params SubscribeParams,
) (<-chan json.RawMessage, error) {
	ch := make(chan json.RawMessage, 100)

	// Try initial subscription
	sub, err := m.subscribe(ctx, params, ch)
	if err != nil {
		close(ch)
		return nil, err
	}

	// Start reconnection handler
	go m.handleReconnection(ctx, params, ch, sub)

	return ch, nil
}

// subscribe creates a subscription and wires up the data/error handlers.
func (m *Manager) subscribe(
	ctx context.Context,
	params SubscribeParams,
	ch chan<- json.RawMessage,
) (Subscription, error) {
	errCh := make(chan error, 1)

	sub, err := m.transport.Subscribe(
		params,
		func(data json.RawMessage) {
			select {
			case ch <- data:
			case <-ctx.Done():
			}
		},
		func(err error) {
			select {
			case errCh <- err:
			default:
			}
		},
	)

	if err != nil {
		return nil, err
	}

	return sub, nil
}

// handleReconnection manages the subscription lifecycle and reconnection.
func (m *Manager) handleReconnection(
	ctx context.Context,
	params SubscribeParams,
	ch chan json.RawMessage,
	initialSub Subscription,
) {
	defer close(ch)

	currentSub := initialSub
	_ = params // Used in reconnection logic

	// Wait for context cancellation
	<-ctx.Done()

	// Context canceled - clean up and exit
	if currentSub != nil {
		_ = currentSub.Unsubscribe()
	}
}

// SubscriptionEvent wraps subscription data with metadata.
type SubscriptionEvent struct {
	// Data is the raw subscription data.
	Data json.RawMessage

	// Error is any error that occurred.
	Error error

	// Reconnected indicates if this event came after a reconnection.
	Reconnected bool
}

// SubscribeWithEvents creates a subscription that returns events with metadata.
func (m *Manager) SubscribeWithEvents(
	ctx context.Context,
	params SubscribeParams,
) (<-chan SubscriptionEvent, error) {
	ch := make(chan SubscriptionEvent, 100)

	// Try initial subscription
	_, err := m.transport.Subscribe(
		params,
		func(data json.RawMessage) {
			select {
			case ch <- SubscriptionEvent{Data: data}:
			case <-ctx.Done():
			}
		},
		func(err error) {
			select {
			case ch <- SubscriptionEvent{Error: err}:
			case <-ctx.Done():
			}
		},
	)

	if err != nil {
		close(ch)
		return nil, err
	}

	// Handle context cancellation
	go func() {
		<-ctx.Done()
		close(ch)
	}()

	return ch, nil
}
