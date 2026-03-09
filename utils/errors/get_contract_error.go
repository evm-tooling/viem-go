// Package errors provides error handling utilities.
package errors

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/abi"
	"github.com/ChefBingbong/viem-go/client/transport"
	pkgerrors "github.com/ChefBingbong/viem-go/errors"
	"github.com/ChefBingbong/viem-go/utils/rpc"
)

// RPC error codes for execution reverted (from viem TS).
const (
	executionRevertedErrorCode = 3
)

// GetContractErrorParams contains parameters for wrapping contract errors.
type GetContractErrorParams struct {
	ABI          *abi.ABI
	FunctionName string
	Address      *common.Address
	Args         []any
	Sender       *common.Address
}

// GetContractError wraps an error with contract execution context and decodes
// revert data when available (Error(string), Panic(uint256), or custom errors).
func GetContractError(err error, params GetContractErrorParams) error {
	if err == nil {
		return nil
	}

	// AbiDecodingZeroDataError -> ContractFunctionZeroDataError (mirrors viem getContractError L59-60)
	var abiZeroErr *pkgerrors.AbiDecodingZeroDataError
	if errors.As(err, &abiZeroErr) {
		return &pkgerrors.ContractFunctionZeroDataError{FunctionName: params.FunctionName}
	}

	// RawContractError pass-through (mirrors viem getContractError L48-51)
	var rawErr *pkgerrors.RawContractError
	if errors.As(err, &rawErr) {
		return &pkgerrors.ContractFunctionExecutionError{
			Cause:           rawErr,
			FunctionName:    params.FunctionName,
			ContractAddress: params.Address,
			Sender:          params.Sender,
		}
	}

	revertData := GetRevertErrorData(err)
	if len(revertData) == 0 {
		return &pkgerrors.ContractFunctionExecutionError{
			Cause:           err,
			FunctionName:    params.FunctionName,
			ContractAddress: params.Address,
			Sender:          params.Sender,
		}
	}

	// Check if this is an execution reverted error
	if !isExecutionRevertedError(err) {
		return &pkgerrors.ContractFunctionExecutionError{
			Cause:           err,
			FunctionName:    params.FunctionName,
			ContractAddress: params.Address,
			Sender:          params.Sender,
		}
	}

	// Decode revert data
	var cause error
	if params.ABI != nil {
		decoded, decodeErr := params.ABI.DecodeErrorResult(revertData)
		if decodeErr == nil {
			reason := ""
			signature := ""
			var decodedData interface{} = decoded.Args

			switch decoded.ErrorName {
			case "Error":
				if len(decoded.Args) > 0 {
					if s, ok := decoded.Args[0].(string); ok {
						reason = s
					}
				}
			case "Panic":
				reason = fmt.Sprintf("panic (code: %v)", decoded.Args)
			default:
				// Custom error
				if decoded.AbiItem != nil {
					sig, _ := abi.FormatAbiItem(decoded.AbiItem)
					signature = sig
				}
			}

			cause = &pkgerrors.ContractFunctionRevertedError{
				FunctionName: params.FunctionName,
				Reason:       reason,
				Signature:    signature,
				Data:         decodedData,
				Raw:          revertData,
			}
		} else {
			cause = err
		}
	} else {
		// No ABI - try DecodeErrorResultWithoutABI for Error(string) and Panic(uint256)
		decoded, decodeErr := abi.DecodeErrorResultWithoutABI(revertData)
		if decodeErr == nil {
			reason := ""
			switch decoded.ErrorName {
			case "Error":
				if len(decoded.Args) > 0 {
					if s, ok := decoded.Args[0].(string); ok {
						reason = s
					}
				}
			case "Panic":
				reason = fmt.Sprintf("panic (code: %v)", decoded.Args)
			}
			cause = &pkgerrors.ContractFunctionRevertedError{
				FunctionName: params.FunctionName,
				Reason:       reason,
				Raw:          revertData,
			}
		} else {
			cause = err
		}
	}

	return &pkgerrors.ContractFunctionExecutionError{
		Cause:           cause,
		FunctionName:    params.FunctionName,
		ContractAddress: params.Address,
		Sender:          params.Sender,
	}
}

// isExecutionRevertedError returns true if the error indicates an execution revert.
func isExecutionRevertedError(err error) bool {
	var rpcErr *transport.RPCRequestError
	if !errors.As(err, &rpcErr) || rpcErr.RPCError == nil {
		return false
	}

	code := rpcErr.RPCError.Code
	msg := strings.ToLower(rpcErr.RPCError.Message)

	// Code 3 (execution reverted) or -32603 (internal)
	if code == executionRevertedErrorCode || code == rpc.RPCErrorCodeInternal {
		return true
	}

	// Code -32000 (invalid input) with "execution reverted" in message
	if code == rpc.RPCErrorCodeInvalidInput && strings.Contains(msg, "execution reverted") {
		return true
	}

	return false
}
