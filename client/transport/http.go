package transport

import (
	"context"
	"fmt"
	"time"

	"github.com/ChefBingbong/viem-go/utils/rpc"
)

// HTTPTransportConfig contains configuration for the HTTP transport.
type HTTPTransportConfig struct {
	// URL is the RPC endpoint URL.
	URL string
	// Batch enables JSON-RPC batching.
	Batch *BatchConfig
	// Key is the transport key.
	Key string
	// Name is the transport name.
	Name string
	// Methods specifies which RPC methods to allow/block.
	Methods *MethodFilter
	// RetryCount is the maximum number of retry attempts.
	RetryCount int
	// RetryDelay is the base delay between retries.
	RetryDelay time.Duration
	// Timeout is the request timeout.
	Timeout time.Duration
	// Headers are additional HTTP headers.
	Headers map[string]string
	// OnRequest is called before each request.
	OnRequest func(req any) error
	// OnResponse is called after each response.
	OnResponse func(resp any) error
	// Raw returns RPC errors as responses instead of throwing.
	Raw bool
}

// BatchConfig contains batching configuration.
type BatchConfig struct {
	// Enabled enables batching.
	Enabled bool
	// BatchSize is the maximum number of requests per batch.
	BatchSize int
	// Wait is the maximum time to wait before sending a batch.
	Wait time.Duration
}

// DefaultHTTPTransportConfig returns default HTTP transport configuration.
func DefaultHTTPTransportConfig() HTTPTransportConfig {
	return HTTPTransportConfig{
		Key:        "http",
		Name:       "HTTP JSON-RPC",
		RetryCount: 3,
		RetryDelay: 150 * time.Millisecond,
		Timeout:    10 * time.Second,
	}
}

// HTTPTransport implements Transport over HTTP.
type HTTPTransport struct {
	config         HTTPTransportConfig
	client         *rpc.HTTPClient
	batchScheduler *rpc.BatchScheduler
}

// HTTP creates a new HTTP transport factory.
func HTTP(url string, config ...HTTPTransportConfig) TransportFactory {
	return func(params TransportParams) (Transport, error) {
		cfg := DefaultHTTPTransportConfig()
		if len(config) > 0 {
			cfg = config[0]
		}

		// Use URL from config or parameter, or from chain
		finalURL := url
		if finalURL == "" {
			finalURL = cfg.URL
		}
		if finalURL == "" && params.Chain != nil {
			if len(params.Chain.RPCUrls.Default.HTTP) > 0 {
				finalURL = params.Chain.RPCUrls.Default.HTTP[0]
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

		return NewHTTPTransport(cfg)
	}
}

// NewHTTPTransport creates a new HTTP transport.
func NewHTTPTransport(config HTTPTransportConfig) (*HTTPTransport, error) {
	// Create HTTP client
	clientOpts := rpc.HTTPClientOptions{
		Timeout: config.Timeout,
		Headers: config.Headers,
	}

	client, err := rpc.NewHTTPClient(config.URL, clientOpts)
	if err != nil {
		return nil, err
	}

	transport := &HTTPTransport{
		config: config,
		client: client,
	}

	// Create batch scheduler if batching is enabled
	if config.Batch != nil && config.Batch.Enabled {
		batchOpts := rpc.BatchSchedulerOptions{
			BatchSize: config.Batch.BatchSize,
			Wait:      config.Batch.Wait,
		}
		if batchOpts.BatchSize == 0 {
			batchOpts.BatchSize = 1000
		}
		transport.batchScheduler = rpc.NewBatchScheduler(client, batchOpts)
	}

	return transport, nil
}

// Config returns the transport configuration.
func (t *HTTPTransport) Config() TransportConfig {
	return TransportConfig{
		Name:       t.config.Name,
		Key:        t.config.Key,
		Type:       "http",
		Methods:    t.config.Methods,
		RetryCount: t.config.RetryCount,
		RetryDelay: t.config.RetryDelay,
		Timeout:    t.config.Timeout,
	}
}

// Request sends a JSON-RPC request.
func (t *HTTPTransport) Request(ctx context.Context, req RPCRequest) (*RPCResponse, error) {
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

	// Use batch scheduler if available
	if t.batchScheduler != nil {
		return t.batchedRequest(ctx, body)
	}

	// Send request with retry
	return t.retryRequest(ctx, body)
}

// batchedRequest sends a request through the batch scheduler.
func (t *HTTPTransport) batchedRequest(ctx context.Context, body RPCRequest) (*RPCResponse, error) {
	resp, err := t.batchScheduler.Schedule(ctx, body)
	if err != nil {
		return nil, err
	}

	// Handle raw mode
	if t.config.Raw {
		return resp, nil
	}

	// Check for RPC error
	if resp.Error != nil {
		return nil, &RPCRequestError{
			URL:      t.config.URL,
			Body:     body,
			RPCError: resp.Error,
		}
	}

	return resp, nil
}

// retryRequest sends a request with retry logic.
func (t *HTTPTransport) retryRequest(ctx context.Context, body RPCRequest) (*RPCResponse, error) {
	var lastErr error

	for attempt := 0; attempt <= t.config.RetryCount; attempt++ {
		resp, err := t.client.Request(ctx, body)

		if err == nil {
			// Handle raw mode
			if t.config.Raw {
				return resp, nil
			}

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
			delay := t.calculateRetryDelay(attempt, lastErr)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}
	}

	return nil, lastErr
}

// calculateRetryDelay calculates the delay before the next retry.
func (t *HTTPTransport) calculateRetryDelay(attempt int, err error) time.Duration {
	// Check for Retry-After header
	if httpErr, ok := err.(*HTTPRequestError); ok {
		if retryAfter, exists := httpErr.Headers["Retry-After"]; exists {
			var seconds int
			if _, parseErr := fmt.Sscanf(retryAfter, "%d", &seconds); parseErr == nil {
				return time.Duration(seconds) * time.Second
			}
		}
	}

	// Exponential backoff: baseDelay * 2^attempt
	return t.config.RetryDelay * time.Duration(1<<attempt)
}

// Value returns transport-specific attributes.
func (t *HTTPTransport) Value() *TransportValue {
	return &TransportValue{
		URL: t.config.URL,
		Attributes: map[string]any{
			"fetchOptions": t.config.Headers,
		},
	}
}

// Close closes the transport.
func (t *HTTPTransport) Close() error {
	if t.batchScheduler != nil {
		t.batchScheduler.Close()
	}
	return t.client.Close()
}

// URL returns the transport URL.
func (t *HTTPTransport) URL() string {
	return t.config.URL
}

// Client returns the underlying HTTP client.
func (t *HTTPTransport) Client() *rpc.HTTPClient {
	return t.client
}
