package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketClientOptions contains options for the WebSocket RPC client.
type WebSocketClientOptions struct {
	// KeepAlive enables keep-alive pings.
	KeepAlive *KeepAliveConfig
	// Reconnect enables automatic reconnection.
	Reconnect *ReconnectConfig
}

// KeepAliveConfig contains keep-alive configuration.
type KeepAliveConfig struct {
	// Enabled specifies whether keep-alive is enabled.
	Enabled bool
	// Interval is the interval between keep-alive pings.
	Interval time.Duration
}

// DefaultKeepAliveConfig returns default keep-alive configuration.
func DefaultKeepAliveConfig() *KeepAliveConfig {
	return &KeepAliveConfig{
		Enabled:  true,
		Interval: 30 * time.Second,
	}
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

// DefaultReconnectConfig returns default reconnection configuration.
func DefaultReconnectConfig() *ReconnectConfig {
	return &ReconnectConfig{
		Enabled:     true,
		MaxAttempts: 5,
		Delay:       2 * time.Second,
	}
}

// DefaultWebSocketClientOptions returns default options.
func DefaultWebSocketClientOptions() WebSocketClientOptions {
	return WebSocketClientOptions{
		KeepAlive: DefaultKeepAliveConfig(),
		Reconnect: DefaultReconnectConfig(),
	}
}

// callbackFn represents a callback for request/subscription responses.
type callbackFn struct {
	onResponse func(resp RPCResponse)
	onError    func(err error)
	body       *RPCRequest
}

// WebSocketClient is a WebSocket JSON-RPC client.
type WebSocketClient struct {
	url           string
	conn          *websocket.Conn
	dialer        *websocket.Dialer
	keepAlive     *KeepAliveConfig
	reconnect     *ReconnectConfig
	idGen         *IDGenerator
	requests      map[any]*callbackFn
	subscriptions map[string]*callbackFn
	mu            sync.RWMutex
	closed        bool
	closeCh       chan struct{}
	keepAliveTick *time.Ticker
	reconnecting  bool
	reconnectMu   sync.Mutex
}

// NewWebSocketClient creates a new WebSocket RPC client.
func NewWebSocketClient(url string, opts ...WebSocketClientOptions) (*WebSocketClient, error) {
	opt := DefaultWebSocketClientOptions()
	if len(opts) > 0 {
		opt = opts[0]
	}

	client := &WebSocketClient{
		url:           url,
		dialer:        websocket.DefaultDialer,
		keepAlive:     opt.KeepAlive,
		reconnect:     opt.Reconnect,
		idGen:         NewIDGenerator(),
		requests:      make(map[any]*callbackFn),
		subscriptions: make(map[string]*callbackFn),
		closeCh:       make(chan struct{}),
	}

	// Connect
	if err := client.connect(); err != nil {
		return nil, err
	}

	return client, nil
}

// connect establishes the WebSocket connection.
func (c *WebSocketClient) connect() error {
	conn, resp, err := c.dialer.Dial(c.url, nil)
	if resp != nil && resp.Body != nil {
		_ = resp.Body.Close()
	}
	if err != nil {
		return NewWebSocketRequestError(c.url, nil, err)
	}

	c.mu.Lock()
	c.conn = conn
	c.closed = false
	c.mu.Unlock()

	// Start message handler
	go c.handleMessages()

	// Start keep-alive
	if c.keepAlive != nil && c.keepAlive.Enabled {
		c.startKeepAlive()
	}

	return nil
}

// handleMessages reads and processes incoming messages.
func (c *WebSocketClient) handleMessages() {
	for {
		select {
		case <-c.closeCh:
			return
		default:
		}

		c.mu.RLock()
		conn := c.conn
		closed := c.closed
		c.mu.RUnlock()

		if closed || conn == nil {
			return
		}

		_, message, err := conn.ReadMessage()
		if err != nil {
			c.handleError(err)
			return
		}

		// Parse response
		var resp RPCResponse
		if err := json.Unmarshal(message, &resp); err != nil {
			continue // Ignore malformed messages
		}

		c.handleResponse(resp)
	}
}

// handleResponse processes a received response.
func (c *WebSocketClient) handleResponse(resp RPCResponse) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Check if it's a subscription notification
	if resp.Method == "eth_subscription" && resp.Params != nil {
		subID := resp.Params.Subscription
		if callback, ok := c.subscriptions[subID]; ok {
			callback.onResponse(resp)
		}
		return
	}

	// Regular request response
	if callback, ok := c.requests[resp.ID]; ok {
		callback.onResponse(resp)
		delete(c.requests, resp.ID)
	}
}

// handleError processes connection errors.
func (c *WebSocketClient) handleError(err error) {
	c.mu.Lock()

	// Notify all pending requests
	for id, callback := range c.requests {
		if callback.onError != nil {
			callback.onError(err)
		}
		delete(c.requests, id)
	}

	// Notify all subscriptions
	for subID, callback := range c.subscriptions {
		if callback.onError != nil {
			callback.onError(ErrSocketClosed)
		}
		delete(c.subscriptions, subID)
	}

	c.mu.Unlock()

	// Attempt reconnect
	c.attemptReconnect()
}

// attemptReconnect tries to reconnect to the server.
func (c *WebSocketClient) attemptReconnect() {
	if c.reconnect == nil || !c.reconnect.Enabled {
		return
	}

	c.reconnectMu.Lock()
	if c.reconnecting {
		c.reconnectMu.Unlock()
		return
	}
	c.reconnecting = true
	c.reconnectMu.Unlock()

	defer func() {
		c.reconnectMu.Lock()
		c.reconnecting = false
		c.reconnectMu.Unlock()
	}()

	// Close existing connection
	c.mu.Lock()
	if c.conn != nil {
		_ = c.conn.Close()
	}
	c.mu.Unlock()

	// Try to reconnect
	for attempt := 0; attempt < c.reconnect.MaxAttempts; attempt++ {
		select {
		case <-c.closeCh:
			return
		default:
		}

		time.Sleep(c.reconnect.Delay)

		if err := c.connect(); err == nil {
			// Reconnected successfully
			// Re-establish subscriptions
			c.resubscribe()
			return
		}
	}
}

// resubscribe re-establishes subscriptions after reconnect.
func (c *WebSocketClient) resubscribe() {
	c.mu.RLock()
	subs := make(map[string]*callbackFn)
	for k, v := range c.subscriptions {
		if v.body != nil {
			subs[k] = v
		}
	}
	c.mu.RUnlock()

	for _, callback := range subs {
		if callback.body == nil {
			continue
		}

		// Re-send subscription request
		if err := c.Request(*callback.body, callback.onResponse, callback.onError); err != nil {
			fmt.Println("Error resubscrbing request", err)
		}
	}
}

// startKeepAlive starts the keep-alive ping routine.
func (c *WebSocketClient) startKeepAlive() {
	if c.keepAliveTick != nil {
		c.keepAliveTick.Stop()
	}

	c.keepAliveTick = time.NewTicker(c.keepAlive.Interval)

	go func() {
		for {
			select {
			case <-c.closeCh:
				return
			case <-c.keepAliveTick.C:
				c.ping()
			}
		}
	}()
}

// ping sends a keep-alive ping.
func (c *WebSocketClient) ping() {
	c.mu.RLock()
	conn := c.conn
	closed := c.closed
	c.mu.RUnlock()

	if closed || conn == nil {
		return
	}

	// Send a simple request as ping
	body := RPCRequest{
		JSONRPC: "2.0",
		ID:      nil, // No ID for ping
		Method:  "net_version",
		Params:  []any{},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return
	}

	c.mu.Lock()
	err = c.conn.WriteMessage(websocket.TextMessage, data)
	c.mu.Unlock()

	if err != nil {
		c.handleError(err)
	}
}

// Request sends a JSON-RPC request.
func (c *WebSocketClient) Request(
	body RPCRequest,
	onResponse func(resp RPCResponse),
	onError func(err error),
) error {
	c.mu.RLock()
	if c.closed {
		c.mu.RUnlock()
		return ErrSocketClosed
	}
	conn := c.conn
	c.mu.RUnlock()

	// Ensure request has an ID
	if body.ID == nil {
		body.ID = c.idGen.Next()
	}
	if body.JSONRPC == "" {
		body.JSONRPC = "2.0"
	}

	// Register callback
	callback := &callbackFn{
		onResponse: onResponse,
		onError:    onError,
		body:       &body,
	}

	c.mu.Lock()
	c.requests[body.ID] = callback
	c.mu.Unlock()

	// Marshal and send
	data, err := json.Marshal(body)
	if err != nil {
		c.mu.Lock()
		delete(c.requests, body.ID)
		c.mu.Unlock()
		return err
	}

	c.mu.Lock()
	err = conn.WriteMessage(websocket.TextMessage, data)
	c.mu.Unlock()

	if err != nil {
		c.mu.Lock()
		delete(c.requests, body.ID)
		c.mu.Unlock()
		return NewWebSocketRequestError(c.url, body, err)
	}

	return nil
}

// RequestAsync sends a request and waits for the response.
func (c *WebSocketClient) RequestAsync(ctx context.Context, body RPCRequest, timeout time.Duration) (*RPCResponse, error) {
	respCh := make(chan RPCResponse, 1)
	errCh := make(chan error, 1)

	err := c.Request(body, func(resp RPCResponse) {
		respCh <- resp
	}, func(err error) {
		errCh <- err
	})

	if err != nil {
		return nil, err
	}

	// Create timeout context
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	select {
	case resp := <-respCh:
		return &resp, nil
	case err := <-errCh:
		return nil, err
	case <-timeoutCtx.Done():
		return nil, NewTimeoutError(c.url, body)
	}
}

// Subscribe creates a subscription.
func (c *WebSocketClient) Subscribe(
	params []any,
	onData func(data json.RawMessage),
	onError func(err error),
) (*Subscription, error) {
	body := RPCRequest{
		JSONRPC: "2.0",
		ID:      c.idGen.Next(),
		Method:  "eth_subscribe",
		Params:  params,
	}

	respCh := make(chan RPCResponse, 1)
	errCh := make(chan error, 1)

	err := c.Request(body, func(resp RPCResponse) {
		respCh <- resp
	}, func(err error) {
		errCh <- err
	})

	if err != nil {
		return nil, err
	}

	// Wait for subscription confirmation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	select {
	case resp := <-respCh:
		if resp.Error != nil {
			return nil, resp.Error
		}

		// Extract subscription ID
		var subID string
		if err := json.Unmarshal(resp.Result, &subID); err != nil {
			return nil, fmt.Errorf("failed to parse subscription ID: %w", err)
		}

		// Register subscription callback
		callback := &callbackFn{
			onResponse: func(r RPCResponse) {
				if r.Params != nil {
					onData(r.Params.Result)
				}
			},
			onError: onError,
			body:    &body,
		}

		c.mu.Lock()
		c.subscriptions[subID] = callback
		c.mu.Unlock()

		return &Subscription{
			ID: subID,
			Unsubscribe: func() error {
				return c.Unsubscribe(subID)
			},
		}, nil

	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, NewTimeoutError(c.url, body)
	}
}

// Unsubscribe cancels a subscription.
func (c *WebSocketClient) Unsubscribe(subscriptionID string) error {
	c.mu.Lock()
	delete(c.subscriptions, subscriptionID)
	c.mu.Unlock()

	body := RPCRequest{
		JSONRPC: "2.0",
		ID:      c.idGen.Next(),
		Method:  "eth_unsubscribe",
		Params:  []any{subscriptionID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := c.RequestAsync(ctx, body, 10*time.Second)
	return err
}

// Close closes the WebSocket connection.
func (c *WebSocketClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true
	close(c.closeCh)

	if c.keepAliveTick != nil {
		c.keepAliveTick.Stop()
	}

	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}

// URL returns the client URL.
func (c *WebSocketClient) URL() string {
	return c.url
}

// IsConnected returns true if the client is connected.
func (c *WebSocketClient) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return !c.closed && c.conn != nil
}

// GetWebSocketRpcClient returns a WebSocket RPC client (compatibility function).
func GetWebSocketRpcClient(url string, opts ...WebSocketClientOptions) (*WebSocketClient, error) {
	return NewWebSocketClient(url, opts...)
}

// WebSocket RPC client cache
var (
	wsClientCache   = make(map[string]*WebSocketClient)
	wsClientCacheMu sync.RWMutex
)

// GetCachedWebSocketClient returns a cached WebSocket client or creates a new one.
func GetCachedWebSocketClient(url string, opts ...WebSocketClientOptions) (*WebSocketClient, error) {
	wsClientCacheMu.RLock()
	if client, ok := wsClientCache[url]; ok {
		wsClientCacheMu.RUnlock()
		return client, nil
	}
	wsClientCacheMu.RUnlock()

	wsClientCacheMu.Lock()
	defer wsClientCacheMu.Unlock()

	// Double-check after acquiring write lock
	if client, ok := wsClientCache[url]; ok {
		return client, nil
	}

	client, err := NewWebSocketClient(url, opts...)
	if err != nil {
		return nil, err
	}

	wsClientCache[url] = client
	return client, nil
}

// CloseAllCachedWebSocketClients closes all cached WebSocket clients.
func CloseAllCachedWebSocketClients() {
	wsClientCacheMu.Lock()
	defer wsClientCacheMu.Unlock()

	for url, client := range wsClientCache {
		_ = client.Close()
		delete(wsClientCache, url)
	}
}
