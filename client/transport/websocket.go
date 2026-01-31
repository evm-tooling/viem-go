package transport

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ChefBingbong/viem-go/utils/rpc"
)

// WebSocketTransportConfig contains configuration for the WebSocket transport.
type WebSocketTransportConfig struct {
	// URL is the WebSocket RPC endpoint URL.
	URL string
	// Key is the transport key.
	Key string
	// Name is the transport name.
	Name string
	// Methods specifies which RPC methods to allow/block.
	Methods *MethodFilter
	// KeepAlive enables keep-alive pings.
	KeepAlive *KeepAliveConfig
	// Reconnect enables automatic reconnection.
	Reconnect *ReconnectConfig
	// RetryCount is the maximum number of retry attempts.
	RetryCount int
	// RetryDelay is the base delay between retries.
	RetryDelay time.Duration
	// Timeout is the request timeout.
	Timeout time.Duration
}

// KeepAliveConfig contains keep-alive configuration.
type KeepAliveConfig struct {
	// Enabled specifies whether keep-alive is enabled.
	Enabled bool
	// Interval is the interval between keep-alive pings.
	Interval time.Duration
}

// ReconnectConfig contains reconnection configuration.
type ReconnectConfig struct {
	// Enabled specifies whether reconnection is enabled.
	Enabled bool
	// MaxAttempts is the maximum number of reconnection attempts.
	MaxAttempts int
	// Delay is the delay between reconnection attempts.
	Delay time.Duration
}

// DefaultWebSocketTransportConfig returns default WebSocket transport configuration.
func DefaultWebSocketTransportConfig() WebSocketTransportConfig {
	return WebSocketTransportConfig{
		Key:        "webSocket",
		Name:       "WebSocket JSON-RPC",
		RetryCount: 3,
		RetryDelay: 150 * time.Millisecond,
		Timeout:    10 * time.Second,
		KeepAlive: &KeepAliveConfig{
			Enabled:  true,
			Interval: 30 * time.Second,
		},
		Reconnect: &ReconnectConfig{
			Enabled:     true,
			MaxAttempts: 5,
			Delay:       2 * time.Second,
		},
	}
}

// WebSocketTransport implements Transport over WebSocket.
type WebSocketTransport struct {
	config WebSocketTransportConfig
	client *rpc.WebSocketClient
}

// WebSocket creates a new WebSocket transport factory.
func WebSocket(url string, config ...WebSocketTransportConfig) TransportFactory {
	return func(params TransportParams) (Transport, error) {
		cfg := DefaultWebSocketTransportConfig()
		if len(config) > 0 {
			cfg = config[0]
		}

		// Use URL from config or parameter, or from chain
		finalURL := url
		if finalURL == "" {
			finalURL = cfg.URL
		}
		if finalURL == "" && params.Chain != nil {
			if len(params.Chain.RPCUrls.Default.WebSocket) > 0 {
				finalURL = params.Chain.RPCUrls.Default.WebSocket[0]
			}
		}
		if finalURL == "" {
			return nil, ErrURLRequired
		}

		cfg.URL = finalURL

		// Apply parameter overrides
		if params.RetryCount != nil {
			cfg.RetryCount = *params.RetryCount
		}
		if params.Timeout != nil {
			cfg.Timeout = *params.Timeout
		}

		return NewWebSocketTransport(cfg)
	}
}

// NewWebSocketTransport creates a new WebSocket transport.
func NewWebSocketTransport(config WebSocketTransportConfig) (*WebSocketTransport, error) {
	// Build client options
	clientOpts := rpc.WebSocketClientOptions{}

	if config.KeepAlive != nil {
		clientOpts.KeepAlive = &rpc.KeepAliveConfig{
			Enabled:  config.KeepAlive.Enabled,
			Interval: config.KeepAlive.Interval,
		}
	}

	if config.Reconnect != nil {
		clientOpts.Reconnect = &rpc.ReconnectConfig{
			Enabled:     config.Reconnect.Enabled,
			MaxAttempts: config.Reconnect.MaxAttempts,
			Delay:       config.Reconnect.Delay,
		}
	}

	// Create WebSocket client
	client, err := rpc.NewWebSocketClient(config.URL, clientOpts)
	if err != nil {
		return nil, err
	}

	return &WebSocketTransport{
		config: config,
		client: client,
	}, nil
}

// Config returns the transport configuration.
func (t *WebSocketTransport) Config() TransportConfig {
	return TransportConfig{
		Name:       t.config.Name,
		Key:        t.config.Key,
		Type:       "webSocket",
		Methods:    t.config.Methods,
		RetryCount: t.config.RetryCount,
		RetryDelay: t.config.RetryDelay,
		Timeout:    t.config.Timeout,
	}
}

// Request sends a JSON-RPC request.
func (t *WebSocketTransport) Request(ctx context.Context, req RPCRequest) (*RPCResponse, error) {
	// Check method filter
	if t.config.Methods != nil && !t.config.Methods.IsAllowed(req.Method) {
		return nil, ErrMethodNotSupported
	}

	// Convert to transport request
	body := RPCRequest{
		JSONRPC: "2.0",
		ID:      req.ID,
		Method:  req.Method,
		Params:  req.Params,
	}
	if body.ID == nil {
		body.ID = NextID()
	}

	// Send request with retry
	return t.retryRequest(ctx, body)
}

// retryRequest sends a request with retry logic.
func (t *WebSocketTransport) retryRequest(ctx context.Context, body RPCRequest) (*RPCResponse, error) {
	var lastErr error

	for attempt := 0; attempt <= t.config.RetryCount; attempt++ {
		resp, err := t.client.RequestAsync(ctx, body, t.config.Timeout)

		if err == nil {
			// Check for RPC error
			if resp.Error != nil {
				rpcErr := &RPCRequestError{
					URL:      t.config.URL,
					Body:     body,
					RPCError: resp.Error,
				}

				// Check if RPC error is retryable
				if !IsRetryableError(resp.Error) || attempt >= t.config.RetryCount {
					return nil, rpcErr
				}

				lastErr = rpcErr
			} else {
				return resp, nil
			}
		} else {
			lastErr = err
		}

		// Check if error is retryable
		if lastErr != nil && !IsRetryableError(lastErr) {
			return nil, lastErr
		}

		// Check context
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// Wait before retry
		if attempt < t.config.RetryCount {
			delay := t.config.RetryDelay * time.Duration(1<<attempt)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}
	}

	return nil, lastErr
}

// Value returns transport-specific attributes.
func (t *WebSocketTransport) Value() *TransportValue {
	return &TransportValue{
		URL: t.config.URL,
		Attributes: map[string]any{
			"getRpcClient": t.GetRpcClient,
			"subscribe":    t.Subscribe,
		},
	}
}

// Close closes the transport.
func (t *WebSocketTransport) Close() error {
	return t.client.Close()
}

// URL returns the transport URL.
func (t *WebSocketTransport) URL() string {
	return t.config.URL
}

// GetRpcClient returns the underlying WebSocket client.
func (t *WebSocketTransport) GetRpcClient() *rpc.WebSocketClient {
	return t.client
}

// IsConnected returns true if the transport is connected.
func (t *WebSocketTransport) IsConnected() bool {
	return t.client.IsConnected()
}

// SubscribeParams contains parameters for a subscription.
type SubscribeParams struct {
	// Type is the subscription type (newHeads, newPendingTransactions, logs, syncing).
	Type string
	// Params are additional parameters for the subscription.
	Params any
}

// NewHeadsSubscribeParams creates params for newHeads subscription.
func NewHeadsSubscribeParams() SubscribeParams {
	return SubscribeParams{
		Type: "newHeads",
	}
}

// NewPendingTransactionsSubscribeParams creates params for newPendingTransactions subscription.
func NewPendingTransactionsSubscribeParams() SubscribeParams {
	return SubscribeParams{
		Type: "newPendingTransactions",
	}
}

// LogsSubscribeParams creates params for logs subscription.
func LogsSubscribeParams(address any, topics []any) SubscribeParams {
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

// SyncingSubscribeParams creates params for syncing subscription.
func SyncingSubscribeParams() SubscribeParams {
	return SubscribeParams{
		Type: "syncing",
	}
}

// Subscribe creates a subscription on the WebSocket transport.
func (t *WebSocketTransport) Subscribe(
	params SubscribeParams,
	onData func(data json.RawMessage),
	onError func(err error),
) (*Subscription, error) {
	// Build subscription params
	subParams := []any{params.Type}
	if params.Params != nil {
		subParams = append(subParams, params.Params)
	}

	return t.client.Subscribe(subParams, onData, onError)
}

// SubscribeNewHeads subscribes to new block headers.
func (t *WebSocketTransport) SubscribeNewHeads(
	onData func(data json.RawMessage),
	onError func(err error),
) (*Subscription, error) {
	return t.Subscribe(NewHeadsSubscribeParams(), onData, onError)
}

// SubscribeNewPendingTransactions subscribes to new pending transactions.
func (t *WebSocketTransport) SubscribeNewPendingTransactions(
	onData func(data json.RawMessage),
	onError func(err error),
) (*Subscription, error) {
	return t.Subscribe(NewPendingTransactionsSubscribeParams(), onData, onError)
}

// SubscribeLogs subscribes to logs.
func (t *WebSocketTransport) SubscribeLogs(
	address any,
	topics []any,
	onData func(data json.RawMessage),
	onError func(err error),
) (*Subscription, error) {
	return t.Subscribe(LogsSubscribeParams(address, topics), onData, onError)
}

// SubscribeSyncing subscribes to syncing status.
func (t *WebSocketTransport) SubscribeSyncing(
	onData func(data json.RawMessage),
	onError func(err error),
) (*Subscription, error) {
	return t.Subscribe(SyncingSubscribeParams(), onData, onError)
}
