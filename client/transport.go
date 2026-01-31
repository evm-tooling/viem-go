package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"
)

// Transport defines the interface for JSON-RPC transport.
type Transport interface {
	// Call sends a JSON-RPC request and returns the raw result.
	Call(ctx context.Context, method string, params ...any) (json.RawMessage, error)
	// Close closes the transport connection.
	Close() error
}

// HTTPTransport implements Transport over HTTP.
type HTTPTransport struct {
	url       string
	client    *http.Client
	requestID uint64
}

// HTTPTransportOptions contains configuration options for HTTPTransport.
type HTTPTransportOptions struct {
	// Timeout is the HTTP request timeout. Default is 30 seconds.
	Timeout time.Duration
	// Headers are additional HTTP headers to include in requests.
	Headers map[string]string
}

// DefaultHTTPTransportOptions returns the default options.
func DefaultHTTPTransportOptions() HTTPTransportOptions {
	return HTTPTransportOptions{
		Timeout: 30 * time.Second,
	}
}

// NewHTTPTransport creates a new HTTP transport with default options.
func NewHTTPTransport(url string) *HTTPTransport {
	return NewHTTPTransportWithOptions(url, DefaultHTTPTransportOptions())
}

// NewHTTPTransportWithOptions creates a new HTTP transport with custom options.
func NewHTTPTransportWithOptions(url string, opts HTTPTransportOptions) *HTTPTransport {
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &HTTPTransport{
		url: url,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Call sends a JSON-RPC request over HTTP.
func (t *HTTPTransport) Call(ctx context.Context, method string, params ...any) (json.RawMessage, error) {
	// Build request
	id := atomic.AddUint64(&t.requestID, 1)
	req := RPCRequest{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  params,
	}

	// Marshal request
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, t.url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	// Send request
	resp, err := t.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check HTTP status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse JSON-RPC response
	var rpcResp RPCResponse
	if err := json.Unmarshal(respBody, &rpcResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check for RPC error
	if rpcResp.Error != nil {
		return nil, rpcResp.Error
	}

	return rpcResp.Result, nil
}

// Close closes the HTTP transport.
func (t *HTTPTransport) Close() error {
	t.client.CloseIdleConnections()
	return nil
}

// URL returns the transport URL.
func (t *HTTPTransport) URL() string {
	return t.url
}

// BatchTransport wraps a transport to support batch requests.
type BatchTransport struct {
	transport Transport
}

// NewBatchTransport creates a new batch transport wrapper.
func NewBatchTransport(transport Transport) *BatchTransport {
	return &BatchTransport{transport: transport}
}

// BatchRequest represents a single request in a batch.
type BatchRequest struct {
	Method string
	Params []any
}

// BatchResponse represents a single response in a batch.
type BatchResponse struct {
	Result json.RawMessage
	Error  error
}

// BatchCall sends multiple requests in a single HTTP call.
// Note: This requires the transport to support batch requests.
func (b *BatchTransport) BatchCall(ctx context.Context, requests []BatchRequest) ([]BatchResponse, error) {
	// For HTTP transport, we need to implement batch support
	httpTransport, ok := b.transport.(*HTTPTransport)
	if !ok {
		// Fall back to sequential calls
		return b.sequentialCall(ctx, requests)
	}

	// Build batch request
	var batch []RPCRequest
	for i, req := range requests {
		batch = append(batch, RPCRequest{
			JSONRPC: "2.0",
			ID:      uint64(i + 1),
			Method:  req.Method,
			Params:  req.Params,
		})
	}

	// Marshal request
	reqBody, err := json.Marshal(batch)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal batch request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, httpTransport.url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	// Send request
	resp, err := httpTransport.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check HTTP status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse JSON-RPC batch response
	var batchResp []RPCResponse
	if err := json.Unmarshal(respBody, &batchResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal batch response: %w", err)
	}

	// Map responses by ID
	responseMap := make(map[uint64]RPCResponse)
	for _, r := range batchResp {
		// Handle ID which can be any type (usually float64 from JSON or uint64)
		var id uint64
		switch v := r.ID.(type) {
		case float64:
			id = uint64(v)
		case int64:
			id = uint64(v)
		case uint64:
			id = v
		case int:
			id = uint64(v)
		default:
			continue // Skip if ID is not a number
		}
		responseMap[id] = r
	}

	// Build result in original order
	results := make([]BatchResponse, len(requests))
	for i := range requests {
		id := uint64(i + 1)
		if r, ok := responseMap[id]; ok {
			if r.Error != nil {
				results[i] = BatchResponse{Error: r.Error}
			} else {
				results[i] = BatchResponse{Result: r.Result}
			}
		} else {
			results[i] = BatchResponse{Error: fmt.Errorf("missing response for request %d", i)}
		}
	}

	return results, nil
}

// sequentialCall falls back to sequential calls when batch is not supported.
//
//nolint:unparam // error is always nil as errors are captured in BatchResponse.Error
func (b *BatchTransport) sequentialCall(ctx context.Context, requests []BatchRequest) ([]BatchResponse, error) {
	results := make([]BatchResponse, len(requests))
	for i, req := range requests {
		result, err := b.transport.Call(ctx, req.Method, req.Params...)
		if err != nil {
			results[i] = BatchResponse{Error: err}
		} else {
			results[i] = BatchResponse{Result: result}
		}
	}
	return results, nil
}

// Close closes the underlying transport.
func (b *BatchTransport) Close() error {
	return b.transport.Close()
}
