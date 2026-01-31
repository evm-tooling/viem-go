package transport

import (
	"errors"

	"github.com/ChefBingbong/viem-go/utils/rpc"
)

// Common transport errors
var (
	// ErrURLRequired is returned when a URL is required but not provided.
	ErrURLRequired = rpc.ErrURLRequired
	// ErrSocketClosed is returned when attempting to use a closed socket.
	ErrSocketClosed = rpc.ErrSocketClosed
	// ErrTimeout is returned when a request times out.
	ErrTimeout = rpc.ErrTimeout
	// ErrMethodNotSupported is returned when a method is not allowed.
	ErrMethodNotSupported = errors.New("method not supported")
)

// RPCRequestError wraps an RPC error response.
type RPCRequestError struct {
	URL      string
	Body     any
	RPCError *RPCError
}

func (e *RPCRequestError) Error() string {
	if e.RPCError != nil {
		return e.RPCError.Error()
	}
	return "RPC request error"
}

func (e *RPCRequestError) Unwrap() error {
	return e.RPCError
}

// RPC error codes - re-export from rpc package
const (
	// Standard JSON-RPC errors
	RPCErrorCodeParse          = rpc.RPCErrorCodeParse
	RPCErrorCodeInvalidRequest = rpc.RPCErrorCodeInvalidRequest
	RPCErrorCodeMethodNotFound = rpc.RPCErrorCodeMethodNotFound
	RPCErrorCodeInvalidParams  = rpc.RPCErrorCodeInvalidParams
	RPCErrorCodeInternal       = rpc.RPCErrorCodeInternal

	// Server errors
	RPCErrorCodeInvalidInput        = rpc.RPCErrorCodeInvalidInput
	RPCErrorCodeResourceNotFound    = rpc.RPCErrorCodeResourceNotFound
	RPCErrorCodeResourceUnavailable = rpc.RPCErrorCodeResourceUnavailable
	RPCErrorCodeTransactionRejected = rpc.RPCErrorCodeTransactionRejected
	RPCErrorCodeMethodNotSupported  = rpc.RPCErrorCodeMethodNotSupported
	RPCErrorCodeLimitExceeded       = rpc.RPCErrorCodeLimitExceeded
	RPCErrorCodeVersionUnsupported  = rpc.RPCErrorCodeVersionUnsupported
)
