package errors

import (
	"fmt"
	"math/big"
)

// ExecutionRevertedError is returned when execution reverts.
type ExecutionRevertedError struct {
	Cause   error
	Message string
}

// ExecutionRevertedErrorCode is the geth error code for execution reverted.
const ExecutionRevertedErrorCode = 3

func (e *ExecutionRevertedError) Error() string {
	reason := e.Message
	if reason != "" {
		// Strip common prefixes
		if len(reason) > 18 && reason[:18] == "execution reverted: " {
			reason = reason[18:]
		} else if reason == "execution reverted" {
			reason = ""
		}
	}
	if reason != "" {
		return fmt.Sprintf("Execution reverted with reason: %s.", reason)
	}
	return "Execution reverted for an unknown reason."
}

func (e *ExecutionRevertedError) Unwrap() error {
	return e.Cause
}

// FeeCapTooHighError is returned when maxFeePerGas exceeds maximum.
type FeeCapTooHighError struct {
	Cause       error
	MaxFeePerGas *big.Int
}

func (e *FeeCapTooHighError) Error() string {
	msg := "The fee cap (`maxFeePerGas`"
	if e.MaxFeePerGas != nil {
		msg += fmt.Sprintf(" = %s gwei", formatGwei(e.MaxFeePerGas))
	}
	msg += ") cannot be higher than the maximum allowed value (2^256-1)."
	return msg
}

func (e *FeeCapTooHighError) Unwrap() error {
	return e.Cause
}

// FeeCapTooLowError is returned when maxFeePerGas is below block base fee.
type FeeCapTooLowError struct {
	Cause       error
	MaxFeePerGas *big.Int
}

func (e *FeeCapTooLowError) Error() string {
	msg := "The fee cap (`maxFeePerGas`"
	if e.MaxFeePerGas != nil {
		msg += fmt.Sprintf(" = %s", formatGwei(e.MaxFeePerGas))
	}
	msg += " gwei) cannot be lower than the block base fee."
	return msg
}

func (e *FeeCapTooLowError) Unwrap() error {
	return e.Cause
}

// NonceTooHighError is returned when nonce is higher than expected.
type NonceTooHighError struct {
	Cause error
	Nonce *uint64
}

func (e *NonceTooHighError) Error() string {
	msg := "Nonce provided for the transaction "
	if e.Nonce != nil {
		msg += fmt.Sprintf("(%d) ", *e.Nonce)
	}
	msg += "is higher than the next one expected."
	return msg
}

func (e *NonceTooHighError) Unwrap() error {
	return e.Cause
}

// NonceTooLowError is returned when nonce is lower than current.
type NonceTooLowError struct {
	Cause error
	Nonce *uint64
}

func (e *NonceTooLowError) Error() string {
	msg := "Nonce provided for the transaction "
	if e.Nonce != nil {
		msg += fmt.Sprintf("(%d) ", *e.Nonce)
	}
	msg += "is lower than the current nonce of the account."
	msg += "\nTry increasing the nonce or find the latest nonce with `getTransactionCount`."
	return msg
}

func (e *NonceTooLowError) Unwrap() error {
	return e.Cause
}

// NonceMaxValueError is returned when nonce exceeds maximum.
type NonceMaxValueError struct {
	Cause error
	Nonce *uint64
}

func (e *NonceMaxValueError) Error() string {
	msg := "Nonce provided for the transaction "
	if e.Nonce != nil {
		msg += fmt.Sprintf("(%d) ", *e.Nonce)
	}
	msg += "exceeds the maximum allowed nonce."
	return msg
}

func (e *NonceMaxValueError) Unwrap() error {
	return e.Cause
}

// InsufficientFundsError is returned when account has insufficient funds.
type InsufficientFundsError struct {
	Cause error
}

func (e *InsufficientFundsError) Error() string {
	msg := "The total cost (gas * gas fee + value) of executing this transaction exceeds the balance of the account."
	msg += "\nThis error could arise when the account does not have enough funds to:"
	msg += "\n - pay for the total gas fee,"
	msg += "\n - pay for the value to send."
	msg += "\n "
	msg += "\nThe cost of the transaction is calculated as `gas * gas fee + value`, where:"
	msg += "\n - `gas` is the amount of gas needed for transaction to execute,"
	msg += "\n - `gas fee` is the gas fee,"
	msg += "\n - `value` is the amount of ether to send to the recipient."
	return msg
}

func (e *InsufficientFundsError) Unwrap() error {
	return e.Cause
}

// IntrinsicGasTooHighError is returned when gas exceeds block limit.
type IntrinsicGasTooHighError struct {
	Cause error
	Gas   *big.Int
}

func (e *IntrinsicGasTooHighError) Error() string {
	msg := "The amount of gas "
	if e.Gas != nil {
		msg += fmt.Sprintf("(%s) ", e.Gas.String())
	}
	msg += "provided for the transaction exceeds the limit allowed for the block."
	return msg
}

func (e *IntrinsicGasTooHighError) Unwrap() error {
	return e.Cause
}

// IntrinsicGasTooLowError is returned when gas is too low.
type IntrinsicGasTooLowError struct {
	Cause error
	Gas   *big.Int
}

func (e *IntrinsicGasTooLowError) Error() string {
	msg := "The amount of gas "
	if e.Gas != nil {
		msg += fmt.Sprintf("(%s) ", e.Gas.String())
	}
	msg += "provided for the transaction is too low."
	return msg
}

func (e *IntrinsicGasTooLowError) Unwrap() error {
	return e.Cause
}

// TransactionTypeNotSupportedError is returned when transaction type is not supported.
type TransactionTypeNotSupportedError struct {
	Cause error
}

func (e *TransactionTypeNotSupportedError) Error() string {
	return "The transaction type is not supported for this chain."
}

func (e *TransactionTypeNotSupportedError) Unwrap() error {
	return e.Cause
}

// TipAboveFeeCapError is returned when tip exceeds fee cap.
type TipAboveFeeCapError struct {
	Cause                error
	MaxPriorityFeePerGas *big.Int
	MaxFeePerGas         *big.Int
}

func (e *TipAboveFeeCapError) Error() string {
	msg := "The provided tip (`maxPriorityFeePerGas`"
	if e.MaxPriorityFeePerGas != nil {
		msg += fmt.Sprintf(" = %s gwei", formatGwei(e.MaxPriorityFeePerGas))
	}
	msg += ") cannot be higher than the fee cap (`maxFeePerGas`"
	if e.MaxFeePerGas != nil {
		msg += fmt.Sprintf(" = %s gwei", formatGwei(e.MaxFeePerGas))
	}
	msg += ")."
	return msg
}

func (e *TipAboveFeeCapError) Unwrap() error {
	return e.Cause
}

// UnknownNodeError is returned when an unknown node error occurs.
type UnknownNodeError struct {
	Cause error
}

func (e *UnknownNodeError) Error() string {
	shortMsg := "unknown"
	if e.Cause != nil {
		shortMsg = e.Cause.Error()
	}
	return fmt.Sprintf("An error occurred while executing: %s", shortMsg)
}

func (e *UnknownNodeError) Unwrap() error {
	return e.Cause
}

func formatGwei(wei *big.Int) string {
	if wei == nil {
		return ""
	}
	gwei := new(big.Int).Div(wei, big.NewInt(1e9))
	return gwei.String()
}
