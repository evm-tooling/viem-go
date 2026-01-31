package rpc

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// HTTPClientOptions contains options for the HTTP RPC client.
type HTTPClientOptions struct {
	// Timeout is the request timeout.
	Timeout time.Duration
	// Headers are additional HTTP headers to send with each request.
	Headers map[string]string
	// HTTPClient allows providing a custom HTTP client.
	HTTPClient *http.Client
	// OnRequest is called before each request is sent.
	OnRequest func(req *http.Request) error
	// OnResponse is called after each response is received.
	OnResponse func(resp *http.Response) error
}

// DefaultHTTPClientOptions returns default options.
func DefaultHTTPClientOptions() HTTPClientOptions {
	return HTTPClientOptions{
		Timeout: 10 * time.Second,
	}
}

// HTTPClient is an HTTP JSON-RPC client.
type HTTPClient struct {
	url        string
	headers    map[string]string
	httpClient *http.Client
	onRequest  func(req *http.Request) error
	onResponse func(resp *http.Response) error
	idGen      *IDGenerator
}

// NewHTTPClient creates a new HTTP RPC client.
func NewHTTPClient(rawURL string, opts ...HTTPClientOptions) (*HTTPClient, error) {
	opt := DefaultHTTPClientOptions()
	if len(opts) > 0 {
		opt = opts[0]
	}

	// Parse URL to extract any embedded credentials
	parsedURL, headers := parseURL(rawURL)

	// Merge headers
	allHeaders := make(map[string]string)
	for k, v := range headers {
		allHeaders[k] = v
	}
	for k, v := range opt.Headers {
		allHeaders[k] = v
	}

	// Create HTTP client
	httpClient := opt.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: opt.Timeout,
		}
	}

	return &HTTPClient{
		url:        parsedURL,
		headers:    allHeaders,
		httpClient: httpClient,
		onRequest:  opt.OnRequest,
		onResponse: opt.OnResponse,
		idGen:      NewIDGenerator(),
	}, nil
}

// Request sends a single JSON-RPC request.
func (c *HTTPClient) Request(ctx context.Context, body RPCRequest) (*RPCResponse, error) {
	// Ensure request has an ID
	if body.ID == nil {
		body.ID = c.idGen.Next()
	}
	if body.JSONRPC == "" {
		body.JSONRPC = "2.0"
	}

	// Send request
	responses, err := c.doRequest(ctx, []RPCRequest{body})
	if err != nil {
		return nil, err
	}

	if len(responses) == 0 {
		return nil, fmt.Errorf("empty response")
	}

	return &responses[0], nil
}

// BatchRequest sends multiple JSON-RPC requests in a single HTTP call.
func (c *HTTPClient) BatchRequest(ctx context.Context, bodies []RPCRequest) ([]RPCResponse, error) {
	// Ensure all requests have IDs
	for i := range bodies {
		if bodies[i].ID == nil {
			bodies[i].ID = c.idGen.Next()
		}
		if bodies[i].JSONRPC == "" {
			bodies[i].JSONRPC = "2.0"
		}
	}

	return c.doRequest(ctx, bodies)
}

// doRequest performs the actual HTTP request.
func (c *HTTPClient) doRequest(ctx context.Context, bodies []RPCRequest) ([]RPCResponse, error) {
	// Marshal request body
	var reqBody []byte
	var bodyErr error

	if len(bodies) == 1 {
		reqBody, bodyErr = json.Marshal(bodies[0])
	} else {
		reqBody, bodyErr = json.Marshal(bodies)
	}
	if bodyErr != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", bodyErr)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, NewHTTPRequestError(c.url, 0, "", bodies, err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	// Call onRequest hook
	if c.onRequest != nil {
		if reqErr := c.onRequest(req); reqErr != nil {
			return nil, reqErr
		}
	}

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, NewHTTPRequestError(c.url, 0, "", bodies, err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Call onResponse hook
	if c.onResponse != nil {
		if respErr := c.onResponse(resp); respErr != nil {
			return nil, respErr
		}
	}

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewHTTPRequestError(c.url, resp.StatusCode, resp.Status, bodies, err)
	}

	// Check HTTP status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Try to parse error response
		var data any
		if json.Unmarshal(respBody, &data) != nil {
			data = string(respBody)
		}
		return nil, NewHTTPRequestError(c.url, resp.StatusCode, resp.Status, data, nil)
	}

	// Parse response
	var responses []RPCResponse

	// Try parsing as array first (batch response)
	if err := json.Unmarshal(respBody, &responses); err != nil {
		// Try parsing as single response
		var singleResp RPCResponse
		if err := json.Unmarshal(respBody, &singleResp); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}
		responses = []RPCResponse{singleResp}
	}

	return responses, nil
}

// Close closes the HTTP client.
func (c *HTTPClient) Close() error {
	c.httpClient.CloseIdleConnections()
	return nil
}

// URL returns the client URL.
func (c *HTTPClient) URL() string {
	return c.url
}

// parseURL extracts authentication from URL and returns clean URL with auth headers.
func parseURL(rawURL string) (string, map[string]string) {
	headers := make(map[string]string)

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL, headers
	}

	// Extract Basic auth credentials
	if parsed.User != nil {
		username := parsed.User.Username()
		password, _ := parsed.User.Password()
		credentials := username + ":" + password
		encoded := base64.StdEncoding.EncodeToString([]byte(credentials))
		headers["Authorization"] = "Basic " + encoded

		// Remove credentials from URL
		parsed.User = nil
	}

	return parsed.String(), headers
}

// WithTimeout returns a new HTTP client with the specified timeout.
func (c *HTTPClient) WithTimeout(timeout time.Duration) *HTTPClient {
	newClient := *c
	newClient.httpClient = &http.Client{
		Timeout:   timeout,
		Transport: c.httpClient.Transport,
	}
	return &newClient
}

// RequestWithRetry sends a request with retry logic.
func (c *HTTPClient) RequestWithRetry(
	ctx context.Context,
	body RPCRequest,
	retryCount int,
	retryDelay time.Duration,
) (*RPCResponse, error) {
	var lastErr error

	for attempt := 0; attempt <= retryCount; attempt++ {
		resp, err := c.Request(ctx, body)
		if err == nil {
			return resp, nil
		}

		lastErr = err

		// Check if error is retryable
		if !IsRetryableError(err) {
			return nil, err
		}

		// Check context
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// Calculate delay with exponential backoff
		if attempt < retryCount {
			delay := calculateRetryDelay(attempt, retryDelay, err)
			time.Sleep(delay)
		}
	}

	return nil, lastErr
}

// calculateRetryDelay calculates the delay before the next retry.
func calculateRetryDelay(attempt int, baseDelay time.Duration, err error) time.Duration {
	// Check for Retry-After header in HTTP errors
	var httpErr *HTTPRequestError
	if ok := isHTTPRequestError(err, &httpErr); ok {
		if retryAfter, exists := httpErr.Headers["Retry-After"]; exists {
			// Parse Retry-After value (assumes seconds)
			var seconds int
			if _, err := fmt.Sscanf(retryAfter, "%d", &seconds); err == nil {
				return time.Duration(seconds) * time.Second
			}
		}
	}

	// Exponential backoff: baseDelay * 2^attempt
	return baseDelay * time.Duration(1<<attempt)
}

// isHTTPRequestError checks if err is an HTTPRequestError and assigns it to target.
func isHTTPRequestError(err error, target **HTTPRequestError) bool {
	if httpErr, ok := err.(*HTTPRequestError); ok {
		*target = httpErr
		return true
	}
	return false
}

// GetHTTPRpcClient returns an HTTP RPC client (compatibility function).
func GetHTTPRpcClient(url string, opts ...HTTPClientOptions) (*HTTPClient, error) {
	return NewHTTPClient(url, opts...)
}

// ExtractURLHeaders extracts any headers embedded in the URL (like auth).
func ExtractURLHeaders(rawURL string) (string, map[string]string) {
	return parseURL(rawURL)
}

// SanitizeURL removes sensitive information from a URL for logging.
func SanitizeURL(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	// Remove password
	if parsed.User != nil {
		if _, hasPassword := parsed.User.Password(); hasPassword {
			parsed.User = url.UserPassword(parsed.User.Username(), "***")
		}
	}

	// Remove API keys from query params
	q := parsed.Query()
	sensitiveParams := []string{"apikey", "api_key", "key", "token", "secret"}
	for _, param := range sensitiveParams {
		if q.Has(param) {
			q.Set(param, "***")
		}
		// Also check case-insensitive
		for k := range q {
			if strings.EqualFold(k, param) {
				q.Set(k, "***")
			}
		}
	}
	parsed.RawQuery = q.Encode()

	return parsed.String()
}
