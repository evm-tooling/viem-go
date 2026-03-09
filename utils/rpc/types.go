package rpc

import (
	"errors"
	"fmt"
	"strings"
	"sync/atomic"

	json "github.com/goccy/go-json"
)

// Common RPC errors
var (
	// ErrURLRequired is returned when a URL is required but not provided.
	ErrURLRequired = errors.New("url is required")
	// ErrSocketClosed is returned when attempting to use a closed socket.
	ErrSocketClosed = errors.New("socket is closed")
	// ErrTimeout is returned when a request times out.
	ErrTimeout = errors.New("request timeout")
)

// RPCRequest represents a JSON-RPC request.
type RPCRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      any    `json:"id"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
}

// RPCResponse represents a JSON-RPC response.
type RPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
	// For subscriptions
	Method string             `json:"method,omitempty"`
	Params *SubscriptionParam `json:"params,omitempty"`
}

// SubscriptionParam represents subscription notification parameters.
type SubscriptionParam struct {
	Subscription string          `json:"subscription"`
	Result       json.RawMessage `json:"result,omitempty"`
}

// RPCError represents a JSON-RPC error.
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (e *RPCError) Error() string {
	if e.Data != nil {
		return fmt.Sprintf("RPC error %d: %s (data: %v)", e.Code, e.Message, e.Data)
	}
	return fmt.Sprintf("RPC error %d: %s", e.Code, e.Message)
}

// HTTPRequestError represents an HTTP request error.
// Mirrors viem's HttpRequestError: body is the JSON-RPC request sent, details is the error message from the response.
type HTTPRequestError struct {
	URL         string
	Status      int
	StatusText  string
	Body        any   // JSON-RPC request body that was sent (for "Request body" in error message)
	Details     string // Error message from response (for "Details" in error message)
	Headers     map[string]string
	Cause error
}

func (e *HTTPRequestError) Error() string {
	msg := "HTTP request failed."
	var parts []string
	if e.Status > 0 {
		parts = append(parts, fmt.Sprintf("Status: %d", e.Status))
	}
	if e.URL != "" {
		parts = append(parts, fmt.Sprintf("URL: %s", e.URL))
	}
	if e.Body != nil {
		bodyStr := formatRequestBody(e.Body)
		if bodyStr != "" {
			parts = append(parts, "Request body: "+bodyStr)
		}
	}
	if len(parts) > 0 {
		msg += "\n\n" + strings.Join(parts, "\n")
	}
	if e.Details != "" {
		msg += "\n\nDetails: \"" + e.Details + "\""
	} else if e.Cause != nil {
		msg += "\n\nDetails: " + e.Cause.Error()
	}
	return msg
}

func (e *HTTPRequestError) Unwrap() error {
	return e.Cause
}

// formatRequestBody formats the RPC request body as JSON for error display (mirrors viem's stringify).
func formatRequestBody(body any) string {
	if body == nil {
		return ""
	}
	// Use json.Marshal for consistent JSON output
	b, err := json.Marshal(body)
	if err != nil {
		return fmt.Sprintf("%v", body)
	}
	return string(b)
}

// NewHTTPRequestError creates a new HTTPRequestError.
// body: the JSON-RPC request that was sent (shown as "Request body").
// details: error message from the response (shown as "Details"); if empty and cause is set, cause.Error() is used.
func NewHTTPRequestError(url string, status int, statusText string, body any, cause error) *HTTPRequestError {
	return &HTTPRequestError{
		URL:        url,
		Status:     status,
		StatusText: statusText,
		Body:       body,
		Cause:      cause,
	}
}

// NewHTTPRequestErrorWithDetails creates an HTTPRequestError with explicit request body and response details.
// Use when the HTTP response body contains an error message to display separately from the request.
func NewHTTPRequestErrorWithDetails(url string, status int, statusText string, requestBody any, details string) *HTTPRequestError {
	return &HTTPRequestError{
		URL:        url,
		Status:     status,
		StatusText: statusText,
		Body:       requestBody,
		Details:    details,
	}
}

// WebSocketRequestError represents a WebSocket request error.
type WebSocketRequestError struct {
	URL   string
	Body  any
	Cause error
}

func (e *WebSocketRequestError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("WebSocket request failed: %v (url: %s)", e.Cause, e.URL)
	}
	return fmt.Sprintf("WebSocket request failed (url: %s)", e.URL)
}

func (e *WebSocketRequestError) Unwrap() error {
	return e.Cause
}

// NewWebSocketRequestError creates a new WebSocketRequestError.
func NewWebSocketRequestError(url string, body any, cause error) *WebSocketRequestError {
	return &WebSocketRequestError{
		URL:   url,
		Body:  body,
		Cause: cause,
	}
}

// TimeoutError represents a request timeout error.
type TimeoutError struct {
	URL  string
	Body any
}

func (e *TimeoutError) Error() string {
	return fmt.Sprintf("request timed out (url: %s)", e.URL)
}

// NewTimeoutError creates a new TimeoutError.
func NewTimeoutError(url string, body any) *TimeoutError {
	return &TimeoutError{
		URL:  url,
		Body: body,
	}
}

// RPC error codes
const (
	// Standard JSON-RPC errors
	RPCErrorCodeParse          = -32700
	RPCErrorCodeInvalidRequest = -32600
	RPCErrorCodeMethodNotFound = -32601
	RPCErrorCodeInvalidParams  = -32602
	RPCErrorCodeInternal       = -32603

	// Server errors
	RPCErrorCodeInvalidInput        = -32000
	RPCErrorCodeResourceNotFound    = -32001
	RPCErrorCodeResourceUnavailable = -32002
	RPCErrorCodeTransactionRejected = -32003
	RPCErrorCodeMethodNotSupported  = -32004
	RPCErrorCodeLimitExceeded       = -32005
	RPCErrorCodeVersionUnsupported  = -32006
)

// IsRetryableError returns true if the error is retryable.
func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Check RPC error codes
	var rpcErr *RPCError
	if errors.As(err, &rpcErr) {
		switch rpcErr.Code {
		case -1,
			RPCErrorCodeLimitExceeded,
			RPCErrorCodeInternal:
			return true
		}
		return false
	}

	// Check HTTP errors
	var httpErr *HTTPRequestError
	if errors.As(err, &httpErr) {
		switch httpErr.Status {
		case 403, 408, 413, 429, 500, 502, 503, 504:
			return true
		}
		return false
	}

	return true
}

// IDGenerator generates unique request IDs.
type IDGenerator struct {
	counter uint64
}

// NewIDGenerator creates a new ID generator.
func NewIDGenerator() *IDGenerator {
	return &IDGenerator{}
}

// Next returns the next ID.
func (g *IDGenerator) Next() uint64 {
	return atomic.AddUint64(&g.counter, 1)
}

// Global ID generator instance.
var globalIDGenerator = NewIDGenerator()

// NextID returns the next global request ID.
func NextID() uint64 {
	return globalIDGenerator.Next()
}

// Subscription represents an active subscription.
type Subscription struct {
	// ID is the subscription ID.
	ID string
	// Unsubscribe cancels the subscription.
	Unsubscribe func() error
}
