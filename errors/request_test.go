package errors

import (
	"strings"
	"testing"
)

func TestHttpRequestError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *HttpRequestError
		contains []string
	}{
		{"basic", &HttpRequestError{URL: "https://rpc.example.com"}, []string{"HTTP request failed", "https://rpc.example.com"}},
		{"with status", &HttpRequestError{URL: "https://x.com", Status: 500, StatusText: "Internal Server Error"}, []string{"HTTP request failed", "500", "https://x.com"}},
		{"with cause", &HttpRequestError{URL: "https://y.com", Cause: &BaseError{ShortMessage: "connection refused"}}, []string{"HTTP request failed", "connection refused"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { testErrorContains(t, tt.err, tt.contains) })
	}
}

func TestWebSocketRequestError_Error(t *testing.T) {
	err := &WebSocketRequestError{URL: "wss://example.com", Cause: &BaseError{ShortMessage: "dial failed"}}
	testErrorContains(t, err, []string{"WebSocket request failed", "wss://example.com", "dial failed"})
}

func TestRpcRequestError_Error(t *testing.T) {
	err := &RpcRequestError{URL: "https://rpc.io", Code: -32603, Message: "internal error"}
	got := err.Error()
	if !strings.Contains(got, "RPC Request failed") {
		t.Errorf("Error() = %q, want to contain 'RPC Request failed'", got)
	}
	if !strings.Contains(got, "internal error") {
		t.Errorf("Error() = %q, want to contain 'internal error'", got)
	}
}

func TestSocketClosedError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *SocketClosedError
		contains []string
	}{
		{"with url", &SocketClosedError{URL: "wss://x.com"}, []string{"socket has been closed", "wss://x.com"}},
		{"without url", &SocketClosedError{}, []string{"socket has been closed"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { testErrorContains(t, tt.err, tt.contains) })
	}
}

func TestTimeoutError_Error(t *testing.T) {
	err := &TimeoutError{URL: "https://slow.com"}
	testErrorContains(t, err, []string{"took too long", "timed out", "https://slow.com"})
}
