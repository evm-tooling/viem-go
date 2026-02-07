package types

import (
	"fmt"

	json "github.com/goccy/go-json"
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

// Error implements the error interface.
func (e *RPCError) Error() string {
	if e.Data != nil {
		return fmt.Sprintf("RPC error %d: %s (data: %v)", e.Code, e.Message, e.Data)
	}
	return fmt.Sprintf("RPC error %d: %s", e.Code, e.Message)
}

// RPC error codes
const (
	// Standard JSON-RPC errors
	RPCErrorCodeParse          = -32700 // Parse error
	RPCErrorCodeInvalidRequest = -32600 // Invalid request
	RPCErrorCodeMethodNotFound = -32601 // Method not found
	RPCErrorCodeInvalidParams  = -32602 // Invalid params
	RPCErrorCodeInternal       = -32603 // Internal error

	// Server errors
	RPCErrorCodeInvalidInput        = -32000 // Invalid input
	RPCErrorCodeResourceNotFound    = -32001 // Resource not found
	RPCErrorCodeResourceUnavailable = -32002 // Resource unavailable
	RPCErrorCodeTransactionRejected = -32003 // Transaction rejected
	RPCErrorCodeMethodNotSupported  = -32004 // Method not supported
	RPCErrorCodeLimitExceeded       = -32005 // Limit exceeded
	RPCErrorCodeVersionUnsupported  = -32006 // JSON-RPC version not supported

	// Provider errors (EIP-1193)
	RPCErrorCodeUserRejected      = 4001 // User rejected request
	RPCErrorCodeUnauthorized      = 4100 // Unauthorized
	RPCErrorCodeUnsupportedMethod = 4200 // Unsupported method
	RPCErrorCodeDisconnected      = 4900 // Provider disconnected
	RPCErrorCodeChainDisconnected = 4901 // Chain disconnected
	RPCErrorCodeSwitchChainError  = 4902 // Switch chain error
)

// IsRetryableError returns true if the RPC error code indicates a retryable error.
func (e *RPCError) IsRetryableError() bool {
	switch e.Code {
	case -1, RPCErrorCodeLimitExceeded, RPCErrorCodeInternal:
		return true
	default:
		return false
	}
}
