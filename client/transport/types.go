package transport

import (
	"context"
	"time"

	"github.com/ChefBingbong/viem-go/utils/rpc"
)

// Re-export types from rpc package for convenience
type (
	RPCRequest        = rpc.RPCRequest
	RPCResponse       = rpc.RPCResponse
	RPCError          = rpc.RPCError
	SubscriptionParam = rpc.SubscriptionParam
	Subscription      = rpc.Subscription
)

// Re-export error types from rpc package
type (
	HTTPRequestError      = rpc.HTTPRequestError
	WebSocketRequestError = rpc.WebSocketRequestError
	TimeoutError          = rpc.TimeoutError
)

// Re-export error constructors
var (
	NewHTTPRequestError      = rpc.NewHTTPRequestError
	NewWebSocketRequestError = rpc.NewWebSocketRequestError
	NewTimeoutError          = rpc.NewTimeoutError
)

// Re-export utility functions
var (
	IsRetryableError = rpc.IsRetryableError
	NextID           = rpc.NextID
)

// MethodFilter specifies which methods to include or exclude.
type MethodFilter struct {
	// Include specifies methods to allow. If set, only these methods are allowed.
	Include []string
	// Exclude specifies methods to block. If set, these methods are blocked.
	Exclude []string
}

// IsAllowed checks if a method is allowed by the filter.
func (f *MethodFilter) IsAllowed(method string) bool {
	if f == nil {
		return true
	}

	// Check exclusions first
	for _, m := range f.Exclude {
		if m == method {
			return false
		}
	}

	// If include list is empty, all methods are allowed (except excluded)
	if len(f.Include) == 0 {
		return true
	}

	// Check inclusions
	for _, m := range f.Include {
		if m == method {
			return true
		}
	}

	return false
}

// TransportConfig contains the configuration for a transport.
type TransportConfig struct {
	// Name is a human-readable name for the transport.
	Name string
	// Key is a unique identifier for the transport type.
	Key string
	// Type is the transport type (e.g., "http", "webSocket", "custom").
	Type string
	// Methods specifies which RPC methods to allow/block.
	Methods *MethodFilter
	// RetryCount is the maximum number of retry attempts.
	RetryCount int
	// RetryDelay is the base delay between retries in milliseconds.
	RetryDelay time.Duration
	// Timeout is the request timeout.
	Timeout time.Duration
}

// DefaultTransportConfig returns default transport configuration.
func DefaultTransportConfig() TransportConfig {
	return TransportConfig{
		Name:       "Transport",
		Key:        "transport",
		Type:       "custom",
		RetryCount: 3,
		RetryDelay: 150 * time.Millisecond,
		Timeout:    10 * time.Second,
	}
}

// TransportValue contains transport-specific attributes.
type TransportValue struct {
	// URL is the RPC endpoint URL.
	URL string
	// Custom attributes for specific transports.
	Attributes map[string]any
}

// Transport represents a JSON-RPC transport.
type Transport interface {
	// Config returns the transport configuration.
	Config() TransportConfig
	// Request sends a JSON-RPC request.
	Request(ctx context.Context, req RPCRequest) (*RPCResponse, error)
	// Value returns transport-specific attributes.
	Value() *TransportValue
	// Close closes the transport.
	Close() error
}

// TransportParams contains parameters passed when creating a transport instance.
type TransportParams struct {
	// Chain is the chain configuration (optional).
	Chain *Chain
	// PollingInterval is the interval for polling operations.
	PollingInterval time.Duration
	// RetryCount overrides the default retry count.
	RetryCount *int
	// Timeout overrides the default timeout.
	Timeout *time.Duration
}

// Chain represents minimal chain information needed by transports.
type Chain struct {
	// ID is the chain ID.
	ID int
	// Name is the chain name.
	Name string
	// RPCUrls contains the RPC URLs for the chain.
	RPCUrls ChainRPCUrls
	// BlockTime is the expected block time in milliseconds.
	BlockTime int
}

// ChainRPCUrls contains RPC URLs for different transport types.
type ChainRPCUrls struct {
	// Default contains the default RPC URLs.
	Default ChainRPCEndpoints
	// Public contains public RPC URLs.
	Public ChainRPCEndpoints
}

// ChainRPCEndpoints contains HTTP and WebSocket endpoints.
type ChainRPCEndpoints struct {
	// HTTP contains HTTP RPC URLs.
	HTTP []string
	// WebSocket contains WebSocket RPC URLs.
	WebSocket []string
}

// TransportFactory is a function that creates a transport instance.
type TransportFactory func(params TransportParams) (Transport, error)

// EIP1193RequestFn is a function that sends an EIP-1193 request.
type EIP1193RequestFn func(ctx context.Context, method string, params ...any) (*RPCResponse, error)

// RequestOptions contains options for a single request.
type RequestOptions struct {
	// RetryCount overrides retry count for this request.
	RetryCount *int
	// RetryDelay overrides retry delay for this request.
	RetryDelay *time.Duration
	// Timeout overrides timeout for this request.
	Timeout *time.Duration
	// Dedupe enables request deduplication.
	Dedupe bool
}

// BatchSchedulerConfig contains configuration for batch scheduling.
type BatchSchedulerConfig struct {
	// BatchSize is the maximum number of requests per batch.
	BatchSize int
	// Wait is the maximum time to wait before sending a batch.
	Wait time.Duration
}

// DefaultBatchSchedulerConfig returns default batch scheduler configuration.
func DefaultBatchSchedulerConfig() BatchSchedulerConfig {
	return BatchSchedulerConfig{
		BatchSize: 1000,
		Wait:      0,
	}
}

// SubscriptionHandler handles subscription events.
type SubscriptionHandler struct {
	// OnData is called when subscription data is received.
	OnData func(data []byte)
	// OnError is called when an error occurs.
	OnError func(err error)
}
