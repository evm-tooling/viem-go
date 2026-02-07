package public

import (
	"time"

	json "github.com/goccy/go-json"

	"github.com/ChefBingbong/viem-go/client/transport"
)

// WatchClient extends the Client interface with methods required for watch actions.
// This interface adds transport detection and subscription capabilities needed
// for polling vs WebSocket subscription decisions.
//
// Watch actions like WatchBlocks, WatchBlockNumber, and WatchEvent require
// additional capabilities beyond the base Client interface:
//   - Transport type detection (to choose polling vs subscription)
//   - Polling interval configuration
//   - WebSocket subscription support
//
// Example implementation:
//
//	type PublicClient struct {
//	    // ... base client fields
//	    transportType   string
//	    pollingInterval time.Duration
//	    wsTransport     *transport.WebSocketTransport
//	}
//
//	func (c *PublicClient) TransportType() string {
//	    return c.transportType
//	}
//
//	func (c *PublicClient) PollingInterval() time.Duration {
//	    return c.pollingInterval
//	}
//
//	func (c *PublicClient) Subscribe(...) (*transport.Subscription, error) {
//	    return c.wsTransport.Subscribe(...)
//	}
type WatchClient interface {
	// Embed the base Client interface
	Client

	// TransportType returns the type of transport being used.
	// Returns one of: "http", "webSocket", "ipc", "fallback", "custom"
	//
	// This is used to determine whether to use polling or subscriptions:
	//   - "http": Use polling
	//   - "webSocket" or "ipc": Use subscriptions
	//   - "fallback": Check first transport in the fallback chain
	//   - "custom": Defaults to polling unless explicitly configured
	TransportType() string

	// PollingInterval returns the default polling interval for this client.
	// This is used by watch actions when no explicit polling interval is provided.
	//
	// Typical values:
	//   - 4 seconds for mainnet
	//   - 1-2 seconds for L2s with faster block times
	PollingInterval() time.Duration

	// Subscribe creates a WebSocket subscription.
	// This is only available when TransportType() returns "webSocket" or "ipc".
	//
	// Parameters:
	//   - params: Subscription parameters (type and optional filter)
	//   - onData: Callback for subscription data
	//   - onError: Callback for errors
	//
	// Returns the subscription (for unsubscribing) or an error.
	//
	// Returns ErrSubscriptionNotSupported if the transport doesn't support subscriptions.
	Subscribe(
		params transport.SubscribeParams,
		onData func(data json.RawMessage),
		onError func(err error),
	) (*transport.Subscription, error)
}

// TransportType constants for identifying the transport type.
const (
	TransportTypeHTTP      = "http"
	TransportTypeWebSocket = "webSocket"
	TransportTypeIPC       = "ipc"
	TransportTypeFallback  = "fallback"
	TransportTypeCustom    = "custom"
)

// DefaultPollingInterval is the default polling interval for watch actions.
// This matches viem's default of 4 seconds for mainnet.
const DefaultPollingInterval = 4 * time.Second

// ShouldPoll determines whether to use polling or subscriptions for a watch action.
//
// Parameters:
//   - client: The watch client
//   - poll: Optional explicit poll preference from parameters
//
// Returns true if polling should be used, false for subscriptions.
//
// Logic:
//   - If poll is explicitly set, use that value
//   - If transport is WebSocket or IPC, use subscriptions (return false)
//   - If transport is fallback and first transport is WebSocket/IPC, use subscriptions
//   - Otherwise, use polling (return true)
func ShouldPoll(client WatchClient, poll *bool) bool {
	// If explicitly configured, use that
	if poll != nil {
		return *poll
	}

	// Check transport type
	transportType := client.TransportType()
	switch transportType {
	case TransportTypeWebSocket, TransportTypeIPC:
		return false
	case TransportTypeFallback:
		// For fallback, check if it's WebSocket capable
		// In viem, it checks the first transport
		// Here we'll default to polling for safety
		return true
	default:
		// HTTP and custom transports use polling
		return true
	}
}

// GetPollingInterval returns the polling interval to use for a watch action.
//
// Parameters:
//   - client: The watch client
//   - configuredInterval: The interval from parameters (may be 0)
//
// Returns the interval to use, defaulting to the client's polling interval
// or the global default.
func GetPollingInterval(client WatchClient, configuredInterval time.Duration) time.Duration {
	if configuredInterval > 0 {
		return configuredInterval
	}

	clientInterval := client.PollingInterval()
	if clientInterval > 0 {
		return clientInterval
	}

	return DefaultPollingInterval
}

// WatchClientAdapter wraps a basic Client to implement WatchClient.
// This is useful when you want to use watch actions with a client that
// doesn't natively support the WatchClient interface.
//
// The adapter defaults to polling mode since it doesn't have access
// to WebSocket subscription capabilities.
//
// Example:
//
//	basicClient := // ... your Client implementation
//	watchClient := public.NewWatchClientAdapter(basicClient, public.WatchClientAdapterOptions{
//	    TransportType:   public.TransportTypeHTTP,
//	    PollingInterval: 4 * time.Second,
//	})
//
//	// Now you can use watch actions
//	events := public.WatchBlockNumber(ctx, watchClient, params)
type WatchClientAdapter struct {
	Client
	transportType   string
	pollingInterval time.Duration
}

// WatchClientAdapterOptions configures the WatchClientAdapter.
type WatchClientAdapterOptions struct {
	// TransportType is the transport type to report.
	// Defaults to "http" (polling mode).
	TransportType string

	// PollingInterval is the polling interval to use.
	// Defaults to 4 seconds.
	PollingInterval time.Duration
}

// NewWatchClientAdapter creates a new WatchClientAdapter.
func NewWatchClientAdapter(client Client, opts ...WatchClientAdapterOptions) *WatchClientAdapter {
	opt := WatchClientAdapterOptions{
		TransportType:   TransportTypeHTTP,
		PollingInterval: DefaultPollingInterval,
	}
	if len(opts) > 0 {
		opt = opts[0]
	}

	return &WatchClientAdapter{
		Client:          client,
		transportType:   opt.TransportType,
		pollingInterval: opt.PollingInterval,
	}
}

// TransportType implements WatchClient.
func (a *WatchClientAdapter) TransportType() string {
	return a.transportType
}

// PollingInterval implements WatchClient.
func (a *WatchClientAdapter) PollingInterval() time.Duration {
	return a.pollingInterval
}

// Subscribe implements WatchClient.
// The adapter doesn't support subscriptions, so this always returns an error.
func (a *WatchClientAdapter) Subscribe(
	params transport.SubscribeParams,
	onData func(data json.RawMessage),
	onError func(err error),
) (*transport.Subscription, error) {
	return nil, ErrSubscriptionNotSupported
}

// Errors for watch client operations.
var (
	// ErrSubscriptionNotSupported is returned when trying to subscribe
	// on a client that doesn't support subscriptions.
	ErrSubscriptionNotSupported = &WatchError{
		Message: "subscriptions are not supported on this transport; use polling mode instead",
	}
)

// WatchError represents an error in watch operations.
type WatchError struct {
	Message string
	Cause   error
}

func (e *WatchError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

func (e *WatchError) Unwrap() error {
	return e.Cause
}
