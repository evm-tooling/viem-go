package transport

import (
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/utils/rpc"

	pkgerrors "github.com/ChefBingbong/viem-go/errors"
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

// ErrorData returns the revert data from the RPC error for eth_call reverts.
// Implements the interface used by utils/errors.GetRevertErrorData.
func (e *RPCRequestError) ErrorData() interface{} {
	if e.RPCError == nil {
		return nil
	}
	d := e.RPCError.Data
	if d == nil {
		return nil
	}
	// Handle string "0x..."
	if s, ok := d.(string); ok {
		s = strings.TrimPrefix(s, "0x")
		if len(s) > 0 {
			return common.FromHex(s)
		}
		return nil
	}
	// Handle nested {"data": "0x..."} or {"data": {"data": "0x..."}} (some nodes)
	if m, ok := d.(map[string]any); ok {
		v := m["data"]
		if v == nil {
			return nil
		}
		if s, ok := v.(string); ok {
			s = strings.TrimPrefix(s, "0x")
			if len(s) > 0 {
				return common.FromHex(s)
			}
		}
		// Nested: {"data": {"data": "0x..."}}
		if nested, ok := v.(map[string]any); ok {
			if s, ok := nested["data"].(string); ok {
				s = strings.TrimPrefix(s, "0x")
				if len(s) > 0 {
					return common.FromHex(s)
				}
			}
		}
	}
	return nil
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

// MapRpcErrorToTyped maps an RPCRequestError to a typed RPC error based on the error code.
// Mirrors viem's buildRequest error code mapping. Returns the original error if RPCError is nil
// or code is unknown. Covers core JSON-RPC codes (-32700 to -32006) and UserRejectedRequestError (4001).
func MapRpcErrorToTyped(rpcErr *RPCRequestError) error {
	if rpcErr == nil {
		return nil
	}
	if rpcErr.RPCError == nil {
		return rpcErr
	}
	code := rpcErr.RPCError.Code
	method := extractMethodFromBody(rpcErr.Body)

	switch code {
	case RPCErrorCodeParse:
		return pkgerrors.NewParseRpcError(rpcErr)
	case RPCErrorCodeInvalidRequest:
		return pkgerrors.NewInvalidRequestRpcError(rpcErr)
	case RPCErrorCodeMethodNotFound:
		return pkgerrors.NewMethodNotFoundRpcError(rpcErr, method)
	case RPCErrorCodeInvalidParams:
		return pkgerrors.NewInvalidParamsRpcError(rpcErr)
	case RPCErrorCodeInternal:
		return pkgerrors.NewInternalRpcError(rpcErr)
	case RPCErrorCodeInvalidInput:
		return pkgerrors.NewInvalidInputRpcError(rpcErr)
	case RPCErrorCodeResourceNotFound:
		return pkgerrors.NewResourceNotFoundRpcError(rpcErr)
	case RPCErrorCodeResourceUnavailable:
		return pkgerrors.NewResourceUnavailableRpcError(rpcErr)
	case RPCErrorCodeTransactionRejected:
		return pkgerrors.NewTransactionRejectedRpcError(rpcErr)
	case RPCErrorCodeMethodNotSupported:
		return pkgerrors.NewMethodNotSupportedRpcError(rpcErr, method)
	case RPCErrorCodeLimitExceeded:
		return pkgerrors.NewLimitExceededRpcError(rpcErr)
	case RPCErrorCodeVersionUnsupported:
		return pkgerrors.NewJsonRpcVersionUnsupportedError(rpcErr)
	case pkgerrors.RpcCodeUserRejected:
		return pkgerrors.NewUserRejectedRequestError(rpcErr)
	case 5000: // CAIP-25: User Rejected Error
		return pkgerrors.NewUserRejectedRequestError(rpcErr)
	default:
		return rpcErr
	}
}

func extractMethodFromBody(body any) string {
	if body == nil {
		return ""
	}
	if req, ok := body.(rpc.RPCRequest); ok {
		return req.Method
	}
	return ""
}
