package errors

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

// FeeConflictError is returned when both gasPrice and maxFeePerGas are specified.
type FeeConflictError struct{}

func (e *FeeConflictError) Error() string {
	return "Cannot specify both a `gasPrice` and a `maxFeePerGas`/`maxPriorityFeePerGas`.\nUse `maxFeePerGas`/`maxPriorityFeePerGas` for EIP-1559 compatible networks, and `gasPrice` for others."
}

// InvalidLegacyVError is returned when v value is invalid (expected 27 or 28).
type InvalidLegacyVError struct {
	V uint64
}

func (e *InvalidLegacyVError) Error() string {
	return fmt.Sprintf("Invalid `v` value \"%d\". Expected 27 or 28.", e.V)
}

// InvalidSerializableTransactionError is returned when transaction type cannot be inferred.
type InvalidSerializableTransactionError struct {
	Transaction string
}

func (e *InvalidSerializableTransactionError) Error() string {
	msg := "Cannot infer a transaction type from provided transaction."
	if e.Transaction != "" {
		msg += "\nProvided Transaction:\n{\n" + e.Transaction + "\n}"
	}
	msg += "\n\nTo infer the type, either provide:"
	msg += "\n- a `type` to the Transaction, or"
	msg += "\n- an EIP-1559 Transaction with `maxFeePerGas`, or"
	msg += "\n- an EIP-2930 Transaction with `gasPrice` & `accessList`, or"
	msg += "\n- an EIP-4844 Transaction with `blobs`, `blobVersionedHashes`, `sidecars`, or"
	msg += "\n- an EIP-7702 Transaction with `authorizationList`, or"
	msg += "\n- a Legacy Transaction with `gasPrice`"
	return msg
}

// InvalidSerializedTransactionTypeError is returned when serialized transaction type is invalid.
type InvalidSerializedTransactionTypeError struct {
	SerializedType string
}

func (e *InvalidSerializedTransactionTypeError) Error() string {
	return fmt.Sprintf("Serialized transaction type \"%s\" is invalid.", e.SerializedType)
}

// InvalidSerializedTransactionError is returned when serialized transaction is invalid.
type InvalidSerializedTransactionError struct {
	SerializedTransaction string
	Type                  string
	MissingAttributes     []string
}

func (e *InvalidSerializedTransactionError) Error() string {
	msg := fmt.Sprintf("Invalid serialized transaction of type \"%s\" was provided.", e.Type)
	msg += "\nSerialized Transaction: \"" + e.SerializedTransaction + "\""
	if len(e.MissingAttributes) > 0 {
		msg += "\nMissing Attributes: " + strings.Join(e.MissingAttributes, ", ")
	}
	return msg
}

// InvalidStorageKeySizeError is returned when storage key size is invalid.
type InvalidStorageKeySizeError struct {
	StorageKey   string
	ExpectedSize int
	ActualSize   int
}

func (e *InvalidStorageKeySizeError) Error() string {
	exp := 32
	if e.ExpectedSize > 0 {
		exp = e.ExpectedSize
	}
	act := e.ActualSize
	if act == 0 && e.StorageKey != "" {
		act = (len(e.StorageKey) - 2) / 2
	}
	return fmt.Sprintf("Size for storage key \"%s\" is invalid. Expected %d bytes. Got %d bytes.", e.StorageKey, exp, act)
}

// TransactionExecutionError is returned when transaction execution fails.
type TransactionExecutionError struct {
	Cause        error
	MetaMessages []string
}

func (e *TransactionExecutionError) Error() string {
	base := "Transaction execution failed."
	if e.Cause != nil {
		base = e.Cause.Error()
	}
	if len(e.MetaMessages) > 0 {
		base += "\n\n" + strings.Join(e.MetaMessages, "\n")
	}
	return "TransactionExecutionError: " + base
}

func (e *TransactionExecutionError) Unwrap() error {
	return e.Cause
}

// TransactionNotFoundError is returned when a transaction is not found.
type TransactionNotFoundError struct {
	Hash        *common.Hash
	BlockHash   *common.Hash
	BlockNumber *uint64
	Index       *int
}

func (e *TransactionNotFoundError) Error() string {
	if e.Hash != nil {
		return fmt.Sprintf("Transaction with hash \"%s\" could not be found.", e.Hash.Hex())
	}
	if e.BlockHash != nil && e.Index != nil {
		return fmt.Sprintf("Transaction at block hash \"%s\" at index \"%d\" could not be found.", e.BlockHash.Hex(), *e.Index)
	}
	if e.BlockNumber != nil && e.Index != nil {
		return fmt.Sprintf("Transaction at block number \"%d\" at index \"%d\" could not be found.", *e.BlockNumber, *e.Index)
	}
	return "Transaction could not be found."
}

// TransactionReceiptNotFoundError is returned when a transaction receipt is not found.
type TransactionReceiptNotFoundError struct {
	Hash *common.Hash
}

func (e *TransactionReceiptNotFoundError) Error() string {
	if e.Hash != nil {
		return fmt.Sprintf("Transaction receipt with hash \"%s\" could not be found. The Transaction may not be processed on a block yet.", e.Hash.Hex())
	}
	return "Transaction receipt could not be found."
}

// TransactionReceiptRevertedError is returned when a transaction receipt indicates revert.
type TransactionReceiptRevertedError struct {
	TransactionHash *common.Hash
}

func (e *TransactionReceiptRevertedError) Error() string {
	if e.TransactionHash != nil {
		return fmt.Sprintf("Transaction with hash \"%s\" reverted.", e.TransactionHash.Hex())
	}
	return "Transaction reverted."
}

// WaitForTransactionReceiptTimeoutError is returned when waiting for receipt times out.
type WaitForTransactionReceiptTimeoutError struct {
	Hash *common.Hash
}

func (e *WaitForTransactionReceiptTimeoutError) Error() string {
	if e.Hash != nil {
		return fmt.Sprintf("Timed out while waiting for transaction with hash \"%s\" to be confirmed.", e.Hash.Hex())
	}
	return "Timed out while waiting for transaction to be confirmed."
}
