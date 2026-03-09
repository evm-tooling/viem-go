package errors

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestFeeConflictError_Error(t *testing.T) {
	err := &FeeConflictError{}
	testErrorContains(t, err, []string{"gasPrice", "maxFeePerGas"})
}

func TestInvalidLegacyVError_Error(t *testing.T) {
	err := &InvalidLegacyVError{V: 30}
	testErrorContains(t, err, []string{"30", "27 or 28"})
}

func TestInvalidSerializableTransactionError_Error(t *testing.T) {
	err := &InvalidSerializableTransactionError{}
	testErrorContains(t, err, []string{"Cannot infer", "transaction type"})
}

func TestInvalidSerializedTransactionTypeError_Error(t *testing.T) {
	err := &InvalidSerializedTransactionTypeError{SerializedType: "0x99"}
	testErrorContains(t, err, []string{"0x99", "invalid"})
}

func TestInvalidSerializedTransactionError_Error(t *testing.T) {
	err := &InvalidSerializedTransactionError{Type: "legacy", SerializedTransaction: "0x01"}
	testErrorContains(t, err, []string{"legacy", "0x01"})
}

func TestInvalidStorageKeySizeError_Error(t *testing.T) {
	err := &InvalidStorageKeySizeError{StorageKey: "0x1234", ActualSize: 2}
	testErrorContains(t, err, []string{"storage", "0x1234", "32"})
}

func TestTransactionExecutionError_Error(t *testing.T) {
	err := &TransactionExecutionError{Cause: &BaseError{ShortMessage: "revert"}}
	testErrorContains(t, err, []string{"revert"})
}

func TestTransactionNotFoundError_Error(t *testing.T) {
	h := common.HexToHash("0xabcd")
	err := &TransactionNotFoundError{Hash: &h}
	testErrorContains(t, err, []string{"could not be found", "Transaction"})
}

func TestTransactionReceiptNotFoundError_Error(t *testing.T) {
	h := common.HexToHash("0x1234")
	err := &TransactionReceiptNotFoundError{Hash: &h}
	testErrorContains(t, err, []string{"could not be found", "receipt"})
}

func TestTransactionReceiptRevertedError_Error(t *testing.T) {
	h := common.HexToHash("0xdead")
	err := &TransactionReceiptRevertedError{TransactionHash: &h}
	testErrorContains(t, err, []string{"reverted", "Transaction"})
}

func TestWaitForTransactionReceiptTimeoutError_Error(t *testing.T) {
	h := common.HexToHash("0xbeef")
	err := &WaitForTransactionReceiptTimeoutError{Hash: &h}
	got := err.Error()
	if !strings.Contains(got, "Timed out") || !strings.Contains(got, "transaction") {
		t.Errorf("Error() = %q, want to contain 'Timed out' and 'transaction'", got)
	}
}
