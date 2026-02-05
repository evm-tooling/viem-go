package public

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/types"
)

// GetTransactionReceiptParameters contains the parameters for the GetTransactionReceipt action.
// This mirrors viem's GetTransactionReceiptParameters type.
type GetTransactionReceiptParameters struct {
	// Hash is the hash of the transaction to retrieve the receipt for. Required.
	Hash common.Hash
}

// GetTransactionReceiptReturnType is the return type for the GetTransactionReceipt action.
type GetTransactionReceiptReturnType = *types.Receipt

// TransactionReceiptNotFoundError is returned when a transaction receipt is not found.
type TransactionReceiptNotFoundError struct {
	Hash common.Hash
}

func (e *TransactionReceiptNotFoundError) Error() string {
	return fmt.Sprintf("transaction receipt not found: hash=%s", e.Hash.Hex())
}

// GetTransactionReceipt returns the transaction receipt given a transaction hash.
//
// This is equivalent to viem's `getTransactionReceipt` action.
//
// JSON-RPC Method: eth_getTransactionReceipt
//
// Example:
//
//	receipt, err := public.GetTransactionReceipt(ctx, client, public.GetTransactionReceiptParameters{
//	    Hash: txHash,
//	})
//	if receipt.IsSuccess() {
//	    fmt.Println("Transaction succeeded")
//	}
func GetTransactionReceipt(ctx context.Context, client Client, params GetTransactionReceiptParameters) (GetTransactionReceiptReturnType, error) {
	// Execute the request
	resp, err := client.Request(ctx, "eth_getTransactionReceipt", params.Hash.Hex())
	if err != nil {
		return nil, fmt.Errorf("eth_getTransactionReceipt failed: %w", err)
	}

	// Check for null result (receipt not found)
	if resp.Result == nil || string(resp.Result) == "null" {
		return nil, &TransactionReceiptNotFoundError{Hash: params.Hash}
	}

	// Parse the receipt
	var receipt types.Receipt
	if unmarshalErr := json.Unmarshal(resp.Result, &receipt); unmarshalErr != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction receipt: %w", unmarshalErr)
	}

	return &receipt, nil
}
