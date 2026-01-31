// Package transport provides JSON-RPC transport implementations for Ethereum clients.
//
// This package mirrors viem's transport architecture, providing:
// - HTTP transport with batching, retry logic, and timeout support
// - WebSocket transport with keep-alive, reconnection, and subscription support
// - Custom transport for user-defined request handlers
// - Fallback transport for trying multiple transports in sequence
//
// Example usage:
//
//	// Create HTTP transport
//	httpTransport := transport.HTTP("https://eth.llamarpc.com")
//
//	// Create HTTP transport with batching
//	httpTransport := transport.HTTP("https://eth.llamarpc.com", transport.HTTPTransportConfig{
//	    Batch: &transport.BatchConfig{
//	        Enabled:   true,
//	        BatchSize: 100,
//	        Wait:      10 * time.Millisecond,
//	    },
//	})
//
//	// Create WebSocket transport with subscriptions
//	wsTransport := transport.WebSocket("wss://eth.llamarpc.com")
//
//	// Create fallback transport
//	fallbackTransport := transport.Fallback(
//	    transport.HTTP("https://eth.llamarpc.com"),
//	    transport.HTTP("https://rpc.ankr.com/eth"),
//	)
package transport

import (
	"context"
	"time"
)

// CreateTransportConfig contains configuration for creating a transport.
type CreateTransportConfig struct {
	// Name is a human-readable name for the transport.
	Name string
	// Key is a unique identifier for the transport type.
	Key string
	// Type is the transport type (e.g., "http", "webSocket", "custom").
	Type string
	// Methods specifies which RPC methods to allow/block.
	Methods *MethodFilter
	// Request is the request function.
	Request func(ctx context.Context, req RPCRequest) (*RPCResponse, error)
	// RetryCount is the maximum number of retry attempts.
	RetryCount int
	// RetryDelay is the base delay between retries.
	RetryDelay time.Duration
	// Timeout is the request timeout.
	Timeout time.Duration
}

// CreateTransportResult contains the result of creating a transport.
type CreateTransportResult struct {
	// Config is the transport configuration.
	Config TransportConfig
	// Request is the request function with retry logic.
	Request func(ctx context.Context, method string, params ...any) (*RPCResponse, error)
	// Value contains transport-specific attributes.
	Value *TransportValue
}

// CreateTransport creates a transport from configuration.
// This function wraps the request function with retry logic and method filtering.
func CreateTransport(config CreateTransportConfig, value *TransportValue) *CreateTransportResult {
	// Apply defaults
	if config.RetryCount == 0 {
		config.RetryCount = 3
	}
	if config.RetryDelay == 0 {
		config.RetryDelay = 150 * time.Millisecond
	}
	if config.Timeout == 0 {
		config.Timeout = 10 * time.Second
	}

	transportConfig := TransportConfig{
		Name:       config.Name,
		Key:        config.Key,
		Type:       config.Type,
		Methods:    config.Methods,
		RetryCount: config.RetryCount,
		RetryDelay: config.RetryDelay,
		Timeout:    config.Timeout,
	}

	// Create request function with retry logic
	request := func(ctx context.Context, method string, params ...any) (*RPCResponse, error) {
		// Check method filter
		if config.Methods != nil && !config.Methods.IsAllowed(method) {
			return nil, ErrMethodNotSupported
		}

		// Build request
		req := RPCRequest{
			JSONRPC: "2.0",
			ID:      NextID(),
			Method:  method,
			Params:  params,
		}

		// Execute with retry
		return executeWithRetry(ctx, config.Request, req, config.RetryCount, config.RetryDelay, config.Timeout)
	}

	return &CreateTransportResult{
		Config:  transportConfig,
		Request: request,
		Value:   value,
	}
}

// executeWithRetry executes a request with retry logic.
func executeWithRetry(
	ctx context.Context,
	requestFn func(ctx context.Context, req RPCRequest) (*RPCResponse, error),
	req RPCRequest,
	retryCount int,
	retryDelay time.Duration,
	timeout time.Duration,
) (*RPCResponse, error) {
	var lastErr error

	for attempt := 0; attempt <= retryCount; attempt++ {
		// Create timeout context
		timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
		resp, err := requestFn(timeoutCtx, req)
		cancel()

		if err == nil {
			if resp.Error != nil {
				// Check if RPC error is retryable
				if !IsRetryableError(resp.Error) || attempt >= retryCount {
					return nil, resp.Error
				}
				lastErr = resp.Error
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
		if attempt < retryCount {
			delay := calculateDelay(attempt, retryDelay, lastErr)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}
	}

	return nil, lastErr
}

// calculateDelay calculates the retry delay.
func calculateDelay(attempt int, baseDelay time.Duration, err error) time.Duration {
	// Check for Retry-After header in HTTP errors
	if httpErr, ok := err.(*HTTPRequestError); ok {
		if retryAfter, exists := httpErr.Headers["Retry-After"]; exists {
			var seconds int
			if _, parseErr := parseRetryAfter(retryAfter); parseErr == nil {
				return time.Duration(seconds) * time.Second
			}
		}
	}

	// Exponential backoff: baseDelay * 2^attempt
	return baseDelay * time.Duration(1<<attempt)
}

// parseRetryAfter parses a Retry-After header value.
func parseRetryAfter(value string) (int, error) {
	var seconds int
	_, err := parseUint(value)
	return seconds, err
}

// parseUint is a simple uint parser for Retry-After.
func parseUint(s string) (uint64, error) {
	var n uint64
	for _, c := range s {
		if c < '0' || c > '9' {
			break
		}
		n = n*10 + uint64(c-'0')
	}
	return n, nil
}

// TransportInstance wraps a Transport with additional functionality.
type TransportInstance struct {
	transport Transport
}

// NewTransportInstance creates a new transport instance.
func NewTransportInstance(factory TransportFactory, params TransportParams) (*TransportInstance, error) {
	t, err := factory(params)
	if err != nil {
		return nil, err
	}
	return &TransportInstance{transport: t}, nil
}

// Request sends a JSON-RPC request.
func (t *TransportInstance) Request(ctx context.Context, method string, params ...any) (*RPCResponse, error) {
	req := RPCRequest{
		JSONRPC: "2.0",
		ID:      NextID(),
		Method:  method,
		Params:  params,
	}
	return t.transport.Request(ctx, req)
}

// Config returns the transport configuration.
func (t *TransportInstance) Config() TransportConfig {
	return t.transport.Config()
}

// Value returns transport-specific attributes.
func (t *TransportInstance) Value() *TransportValue {
	return t.transport.Value()
}

// Close closes the transport.
func (t *TransportInstance) Close() error {
	return t.transport.Close()
}

// Transport returns the underlying transport.
func (t *TransportInstance) Transport() Transport {
	return t.transport
}
