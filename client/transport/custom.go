package transport

import (
	"context"
	"time"
)

// CustomTransportConfig contains configuration for a custom transport.
type CustomTransportConfig struct {
	// Key is the transport key.
	Key string
	// Name is the transport name.
	Name string
	// Methods specifies which RPC methods to allow/block.
	Methods *MethodFilter
	// Request is the custom request function.
	Request func(ctx context.Context, req RPCRequest) (*RPCResponse, error)
	// RetryCount is the maximum number of retry attempts.
	RetryCount int
	// RetryDelay is the base delay between retries.
	RetryDelay time.Duration
	// Timeout is the request timeout.
	Timeout time.Duration
}

// DefaultCustomTransportConfig returns default custom transport configuration.
func DefaultCustomTransportConfig() CustomTransportConfig {
	return CustomTransportConfig{
		Key:        "custom",
		Name:       "Custom JSON-RPC",
		RetryCount: 3,
		RetryDelay: 150 * time.Millisecond,
		Timeout:    10 * time.Second,
	}
}

// CustomTransport implements a custom transport.
type CustomTransport struct {
	config CustomTransportConfig
}

// Custom creates a new custom transport factory.
func Custom(config CustomTransportConfig) TransportFactory {
	return func(params TransportParams) (Transport, error) {
		// Apply parameter overrides
		if params.RetryCount != nil {
			config.RetryCount = *params.RetryCount
		}
		if params.Timeout != nil {
			config.Timeout = *params.Timeout
		}

		return NewCustomTransport(config), nil
	}
}

// NewCustomTransport creates a new custom transport.
func NewCustomTransport(config CustomTransportConfig) *CustomTransport {
	if config.Key == "" {
		config.Key = "custom"
	}
	if config.Name == "" {
		config.Name = "Custom JSON-RPC"
	}

	return &CustomTransport{
		config: config,
	}
}

// Config returns the transport configuration.
func (t *CustomTransport) Config() TransportConfig {
	return TransportConfig{
		Name:       t.config.Name,
		Key:        t.config.Key,
		Type:       "custom",
		Methods:    t.config.Methods,
		RetryCount: t.config.RetryCount,
		RetryDelay: t.config.RetryDelay,
		Timeout:    t.config.Timeout,
	}
}

// Request sends a JSON-RPC request.
func (t *CustomTransport) Request(ctx context.Context, req RPCRequest) (*RPCResponse, error) {
	// Check method filter
	if t.config.Methods != nil && !t.config.Methods.IsAllowed(req.Method) {
		return nil, ErrMethodNotSupported
	}

	// Ensure request has required fields
	if req.ID == nil {
		req.ID = NextID()
	}
	if req.JSONRPC == "" {
		req.JSONRPC = "2.0"
	}

	// Send request with retry
	return t.retryRequest(ctx, req)
}

// retryRequest sends a request with retry logic.
func (t *CustomTransport) retryRequest(ctx context.Context, req RPCRequest) (*RPCResponse, error) {
	var lastErr error

	for attempt := 0; attempt <= t.config.RetryCount; attempt++ {
		// Create timeout context
		timeoutCtx, cancel := context.WithTimeout(ctx, t.config.Timeout)
		resp, err := t.config.Request(timeoutCtx, req)
		cancel()

		if err == nil {
			// Check for RPC error
			if resp.Error != nil {
				// Check if RPC error is retryable
				if !IsRetryableError(resp.Error) || attempt >= t.config.RetryCount {
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
func (t *CustomTransport) Value() *TransportValue {
	return &TransportValue{}
}

// Close closes the transport.
func (t *CustomTransport) Close() error {
	return nil
}
