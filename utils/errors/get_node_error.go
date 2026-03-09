// Package errors provides error handling utilities.
package errors

import (
	"errors"
	"math/big"
	"regexp"
	"strings"

	"github.com/ChefBingbong/viem-go/client/transport"
	pkgerrors "github.com/ChefBingbong/viem-go/errors"
)

// GetNodeErrorParams contains parameters for node error classification.
// Mirrors viem's GetNodeErrorParameters (partial SendTransactionParameters).
type GetNodeErrorParams struct {
	Gas                  *uint64
	GasPrice             *big.Int
	MaxFeePerGas         *big.Int
	MaxPriorityFeePerGas *big.Int
	Nonce                *uint64
	To                   *string
	Value                *big.Int
}

// nodeMessagePatterns are regex patterns for matching node error messages (mirrors viem node.ts).
var nodeMessagePatterns = struct {
	executionReverted     *regexp.Regexp
	feeCapTooHigh         *regexp.Regexp
	feeCapTooLow          *regexp.Regexp
	nonceTooHigh          *regexp.Regexp
	nonceTooLow           *regexp.Regexp
	nonceMaxValue         *regexp.Regexp
	insufficientFunds     *regexp.Regexp
	intrinsicGasTooHigh   *regexp.Regexp
	intrinsicGasTooLow    *regexp.Regexp
	transactionTypeNotSup *regexp.Regexp
	tipAboveFeeCap        *regexp.Regexp
}{
	executionReverted:     regexp.MustCompile(`(?i)execution reverted|gas required exceeds allowance`),
	feeCapTooHigh:         regexp.MustCompile(`(?i)max fee per gas higher than 2\^256-1|fee cap higher than 2\^256-1`),
	feeCapTooLow:          regexp.MustCompile(`(?i)max fee per gas less than block base fee|fee cap less than block base fee|transaction is outdated`),
	nonceTooHigh:          regexp.MustCompile(`(?i)nonce too high`),
	nonceTooLow:           regexp.MustCompile(`(?i)nonce too low|transaction already imported|already known`),
	nonceMaxValue:         regexp.MustCompile(`(?i)nonce has max value`),
	insufficientFunds:     regexp.MustCompile(`(?i)insufficient funds|exceeds transaction sender account balance`),
	intrinsicGasTooHigh:   regexp.MustCompile(`(?i)intrinsic gas too high|gas limit reached`),
	intrinsicGasTooLow:    regexp.MustCompile(`(?i)intrinsic gas too low`),
	transactionTypeNotSup: regexp.MustCompile(`(?i)transaction type not valid`),
	tipAboveFeeCap:        regexp.MustCompile(`(?i)max priority fee per gas higher than max fee per gas|tip higher than fee cap`),
}

// extractErrorMessage extracts the error message from the error chain.
func extractErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	// Check RPCRequestError for RPCError.Message
	var rpcErr *transport.RPCRequestError
	if errors.As(err, &rpcErr) && rpcErr.RPCError != nil {
		return rpcErr.RPCError.Message
	}
	// Fallback to Error()
	return err.Error()
}

// GetNodeError classifies an error into a typed node error (ExecutionRevertedError,
// FeeCapTooHighError, NonceTooHighError, etc.) by matching the message and RPC code.
// Mirrors viem's getNodeError.
func GetNodeError(err error, params GetNodeErrorParams) error {
	if err == nil {
		return nil
	}

	message := extractErrorMessage(err)
	messageLower := strings.ToLower(message)

	// Check for execution reverted by RPC code (code 3)
	var rpcErr *transport.RPCRequestError
	if errors.As(err, &rpcErr) && rpcErr.RPCError != nil && rpcErr.RPCError.Code == pkgerrors.ExecutionRevertedErrorCode {
		return &pkgerrors.ExecutionRevertedError{
			Cause:   err,
			Message: rpcErr.RPCError.Message,
		}
	}

	// Match by message pattern
	if nodeMessagePatterns.executionReverted.MatchString(messageLower) {
		return &pkgerrors.ExecutionRevertedError{Cause: err, Message: message}
	}
	if nodeMessagePatterns.feeCapTooHigh.MatchString(messageLower) {
		return &pkgerrors.FeeCapTooHighError{Cause: err, MaxFeePerGas: params.MaxFeePerGas}
	}
	if nodeMessagePatterns.feeCapTooLow.MatchString(messageLower) {
		return &pkgerrors.FeeCapTooLowError{Cause: err, MaxFeePerGas: params.MaxFeePerGas}
	}
	if nodeMessagePatterns.nonceTooHigh.MatchString(messageLower) {
		return &pkgerrors.NonceTooHighError{Cause: err, Nonce: params.Nonce}
	}
	if nodeMessagePatterns.nonceTooLow.MatchString(messageLower) {
		return &pkgerrors.NonceTooLowError{Cause: err, Nonce: params.Nonce}
	}
	if nodeMessagePatterns.nonceMaxValue.MatchString(messageLower) {
		return &pkgerrors.NonceMaxValueError{Cause: err, Nonce: params.Nonce}
	}
	if nodeMessagePatterns.insufficientFunds.MatchString(messageLower) {
		return &pkgerrors.InsufficientFundsError{Cause: err}
	}
	if nodeMessagePatterns.intrinsicGasTooHigh.MatchString(messageLower) {
		gas := (*big.Int)(nil)
		if params.Gas != nil {
			gas = new(big.Int).SetUint64(*params.Gas)
		}
		return &pkgerrors.IntrinsicGasTooHighError{Cause: err, Gas: gas}
	}
	if nodeMessagePatterns.intrinsicGasTooLow.MatchString(messageLower) {
		gas := (*big.Int)(nil)
		if params.Gas != nil {
			gas = new(big.Int).SetUint64(*params.Gas)
		}
		return &pkgerrors.IntrinsicGasTooLowError{Cause: err, Gas: gas}
	}
	if nodeMessagePatterns.transactionTypeNotSup.MatchString(messageLower) {
		return &pkgerrors.TransactionTypeNotSupportedError{Cause: err}
	}
	if nodeMessagePatterns.tipAboveFeeCap.MatchString(messageLower) {
		return &pkgerrors.TipAboveFeeCapError{
			Cause:                err,
			MaxPriorityFeePerGas: params.MaxPriorityFeePerGas,
			MaxFeePerGas:         params.MaxFeePerGas,
		}
	}

	return &pkgerrors.UnknownNodeError{Cause: err}
}

// IsUnknownNodeError returns true if err is UnknownNodeError.
func IsUnknownNodeError(err error) bool {
	var unknownErr *pkgerrors.UnknownNodeError
	return errors.As(err, &unknownErr)
}

// ContainsNodeError returns true if the error is a transaction-related node error
// (TransactionRejectedRpcError, InvalidInputRpcError, or execution reverted).
func ContainsNodeError(err error) bool {
	if err == nil {
		return false
	}
	// Check for execution reverted by walking chain
	var execErr *pkgerrors.ExecutionRevertedError
	if errors.As(err, &execErr) {
		return true
	}
	// RPCRequestError with code 3 (execution reverted)
	var rpcErr *transport.RPCRequestError
	if errors.As(err, &rpcErr) && rpcErr.RPCError != nil && rpcErr.RPCError.Code == pkgerrors.ExecutionRevertedErrorCode {
		return true
	}
	return false
}
