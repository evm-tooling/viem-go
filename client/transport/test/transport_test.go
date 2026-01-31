package transport_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ChefBingbong/viem-go/client/transport"
)

func TestHTTPTransport_BasicRequest(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req transport.RPCRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)

		resp := map[string]any{
			"jsonrpc": "2.0",
			"id":      req.ID,
			"result":  "0x1",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create transport
	factory := transport.HTTP(server.URL)
	tr, err := factory(transport.TransportParams{})
	require.NoError(t, err)
	defer tr.Close()

	// Send request
	ctx := context.Background()
	resp, err := tr.Request(ctx, transport.RPCRequest{
		Method: "eth_chainId",
		Params: []any{},
	})

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, `"0x1"`, string(resp.Result))
}

func TestHTTPTransport_BatchRequest(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var reqs []transport.RPCRequest
		err := json.NewDecoder(r.Body).Decode(&reqs)
		if err != nil {
			// Not a batch request
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var responses []map[string]any
		for _, req := range reqs {
			responses = append(responses, map[string]any{
				"jsonrpc": "2.0",
				"id":      req.ID,
				"result":  "0x" + req.Method,
			})
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(responses)
	}))
	defer server.Close()

	// Create transport with batching
	config := transport.HTTPTransportConfig{
		URL:  server.URL,
		Key:  "http",
		Name: "HTTP Test",
		Batch: &transport.BatchConfig{
			Enabled:   true,
			BatchSize: 10,
			Wait:      50 * time.Millisecond,
		},
		Timeout: 5 * time.Second,
	}

	tr, err := transport.NewHTTPTransport(config)
	require.NoError(t, err)
	defer tr.Close()

	// Verify config
	cfg := tr.Config()
	assert.Equal(t, "http", cfg.Type)
	assert.Equal(t, "HTTP Test", cfg.Name)
}

func TestHTTPTransport_Error(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req transport.RPCRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)

		resp := map[string]any{
			"jsonrpc": "2.0",
			"id":      req.ID,
			"error": map[string]any{
				"code":    -32601,
				"message": "Method not found",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create transport
	factory := transport.HTTP(server.URL)
	tr, err := factory(transport.TransportParams{})
	require.NoError(t, err)
	defer tr.Close()

	// Send request
	ctx := context.Background()
	_, err = tr.Request(ctx, transport.RPCRequest{
		Method: "eth_unknownMethod",
		Params: []any{},
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "Method not found")
}

func TestHTTPTransport_Timeout(t *testing.T) {
	// Create a test server that delays
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create transport with short timeout
	config := transport.HTTPTransportConfig{
		URL:        server.URL,
		Timeout:    100 * time.Millisecond,
		RetryCount: 0, // No retries for this test
	}

	tr, err := transport.NewHTTPTransport(config)
	require.NoError(t, err)
	defer tr.Close()

	// Send request
	ctx := context.Background()
	_, err = tr.Request(ctx, transport.RPCRequest{
		Method: "eth_chainId",
		Params: []any{},
	})

	require.Error(t, err)
}

func TestHTTPTransport_MethodFilter(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req transport.RPCRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)

		resp := map[string]any{
			"jsonrpc": "2.0",
			"id":      req.ID,
			"result":  "0x1",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create transport with method filter
	config := transport.HTTPTransportConfig{
		URL: server.URL,
		Methods: &transport.MethodFilter{
			Include: []string{"eth_chainId", "eth_blockNumber"},
		},
		Timeout: 5 * time.Second,
	}

	tr, err := transport.NewHTTPTransport(config)
	require.NoError(t, err)
	defer tr.Close()

	ctx := context.Background()

	// Allowed method
	resp, err := tr.Request(ctx, transport.RPCRequest{
		Method: "eth_chainId",
		Params: []any{},
	})
	require.NoError(t, err)
	assert.NotNil(t, resp)

	// Blocked method
	_, err = tr.Request(ctx, transport.RPCRequest{
		Method: "eth_sendTransaction",
		Params: []any{},
	})
	require.Error(t, err)
	assert.Equal(t, transport.ErrMethodNotSupported, err)
}

func TestCustomTransport(t *testing.T) {
	callCount := 0

	// Create custom transport
	factory := transport.Custom(transport.CustomTransportConfig{
		Key:  "custom",
		Name: "Custom Test",
		Request: func(ctx context.Context, req transport.RPCRequest) (*transport.RPCResponse, error) {
			callCount++
			return &transport.RPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`"0x123"`),
			}, nil
		},
	})

	tr, err := factory(transport.TransportParams{})
	require.NoError(t, err)
	defer tr.Close()

	// Send request
	ctx := context.Background()
	resp, err := tr.Request(ctx, transport.RPCRequest{
		Method: "test_method",
		Params: []any{},
	})

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, callCount)
	assert.Equal(t, `"0x123"`, string(resp.Result))
}

func TestFallbackTransport(t *testing.T) {
	primaryCalled := false
	secondaryCalled := false

	// Create primary transport that fails
	primaryFactory := transport.Custom(transport.CustomTransportConfig{
		Key:  "primary",
		Name: "Primary",
		Request: func(ctx context.Context, req transport.RPCRequest) (*transport.RPCResponse, error) {
			primaryCalled = true
			return &transport.RPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error: &transport.RPCError{
					Code:    -32000,
					Message: "Primary failed",
				},
			}, nil
		},
		RetryCount: 0,
	})

	// Create secondary transport that succeeds
	secondaryFactory := transport.Custom(transport.CustomTransportConfig{
		Key:  "secondary",
		Name: "Secondary",
		Request: func(ctx context.Context, req transport.RPCRequest) (*transport.RPCResponse, error) {
			secondaryCalled = true
			return &transport.RPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`"0x456"`),
			}, nil
		},
	})

	// Create fallback transport
	factory := transport.Fallback(primaryFactory, secondaryFactory)
	tr, err := factory(transport.TransportParams{})
	require.NoError(t, err)
	defer tr.Close()

	// Send request
	ctx := context.Background()
	resp, err := tr.Request(ctx, transport.RPCRequest{
		Method: "eth_chainId",
		Params: []any{},
	})

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, primaryCalled)
	assert.True(t, secondaryCalled)
	assert.Equal(t, `"0x456"`, string(resp.Result))
}

func TestTransportConfig(t *testing.T) {
	// Test default config
	cfg := transport.DefaultTransportConfig()
	assert.Equal(t, "Transport", cfg.Name)
	assert.Equal(t, "transport", cfg.Key)
	assert.Equal(t, "custom", cfg.Type)
	assert.Equal(t, 3, cfg.RetryCount)
	assert.Equal(t, 150*time.Millisecond, cfg.RetryDelay)
	assert.Equal(t, 10*time.Second, cfg.Timeout)
}

func TestMethodFilter(t *testing.T) {
	// Test with include list
	filter := &transport.MethodFilter{
		Include: []string{"eth_chainId", "eth_blockNumber"},
	}
	assert.True(t, filter.IsAllowed("eth_chainId"))
	assert.True(t, filter.IsAllowed("eth_blockNumber"))
	assert.False(t, filter.IsAllowed("eth_sendTransaction"))

	// Test with exclude list
	filter = &transport.MethodFilter{
		Exclude: []string{"eth_sendTransaction"},
	}
	assert.True(t, filter.IsAllowed("eth_chainId"))
	assert.False(t, filter.IsAllowed("eth_sendTransaction"))

	// Test nil filter
	assert.True(t, (*transport.MethodFilter)(nil).IsAllowed("anything"))
}

func TestTransportValue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	factory := transport.HTTP(server.URL)
	tr, err := factory(transport.TransportParams{})
	require.NoError(t, err)
	defer tr.Close()

	value := tr.Value()
	assert.NotNil(t, value)
	assert.Equal(t, server.URL, value.URL)
}
