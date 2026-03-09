package errors

import (
	"fmt"
	"strings"
)

// HttpRequestError is returned when an HTTP request fails.
type HttpRequestError struct {
	URL        string
	Status     int
	StatusText string
	Body       interface{}
	Cause      error
}

func (e *HttpRequestError) Error() string {
	msg := "HTTP request failed."
	var parts []string
	if e.Status > 0 {
		parts = append(parts, fmt.Sprintf("Status: %d", e.Status))
		if e.StatusText != "" {
			parts[len(parts)-1] += " " + e.StatusText
		}
	}
	if e.URL != "" {
		parts = append(parts, fmt.Sprintf("URL: %s", e.URL))
	}
	if e.Body != nil {
		parts = append(parts, fmt.Sprintf("Request body: %v", e.Body))
	}
	if len(parts) > 0 {
		msg += "\n" + strings.Join(parts, "\n")
	}
	if e.Cause != nil {
		msg += "\nDetails: " + e.Cause.Error()
	}
	return msg
}

func (e *HttpRequestError) Unwrap() error {
	return e.Cause
}

// WebSocketRequestError is returned when a WebSocket request fails.
type WebSocketRequestError struct {
	URL   string
	Body  interface{}
	Cause error
}

func (e *WebSocketRequestError) Error() string {
	msg := "WebSocket request failed."
	if e.URL != "" {
		msg += "\nURL: " + e.URL
	}
	if e.Body != nil {
		msg += "\nRequest body: " + fmt.Sprintf("%v", e.Body)
	}
	if e.Cause != nil {
		msg += "\nDetails: " + e.Cause.Error()
	}
	return msg
}

func (e *WebSocketRequestError) Unwrap() error {
	return e.Cause
}

// RpcRequestError is returned when an RPC request fails.
type RpcRequestError struct {
	URL     string
	Code    int
	Body    interface{}
	Message string
	Cause   error
}

func (e *RpcRequestError) Error() string {
	msg := "RPC Request failed."
	if e.URL != "" {
		msg += "\nURL: " + e.URL
	}
	if e.Body != nil {
		msg += "\nRequest body: " + fmt.Sprintf("%v", e.Body)
	}
	if e.Message != "" {
		msg += "\nDetails: " + e.Message
	} else if e.Cause != nil {
		msg += "\nDetails: " + e.Cause.Error()
	}
	return msg
}

func (e *RpcRequestError) Unwrap() error {
	return e.Cause
}

// SocketClosedError is returned when the socket has been closed.
type SocketClosedError struct {
	URL string
}

func (e *SocketClosedError) Error() string {
	msg := "The socket has been closed."
	if e.URL != "" {
		msg += "\nURL: " + e.URL
	}
	return msg
}

// TimeoutError is returned when a request times out.
type TimeoutError struct {
	URL  string
	Body interface{}
}

func (e *TimeoutError) Error() string {
	msg := "The request took too long to respond."
	msg += "\nDetails: The request timed out."
	if e.URL != "" {
		msg += "\nURL: " + e.URL
	}
	if e.Body != nil {
		msg += "\nRequest body: " + fmt.Sprintf("%v", e.Body)
	}
	return msg
}
