package wallet

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ChefBingbong/viem-go/utils/encoding"
	"github.com/ChefBingbong/viem-go/utils/formatters"
)

// SendRawTransactionSyncParameters contains the parameters for the SendRawTransactionSync action.
// This mirrors viem's SendRawTransactionSyncParameters type.
type SendRawTransactionSyncParameters struct {
	// SerializedTransaction is the signed serialized transaction hex string.
	SerializedTransaction string

	// ThrowOnReceiptRevert when true, throws an error if the transaction was detected as reverted.
	// Default: true
	ThrowOnReceiptRevert *bool

	// Timeout is the timeout (in ms) for the transaction.
	Timeout *int64
}

// SendRawTransactionSyncReturnType is the return type for the SendRawTransactionSync action.
// It is the formatted transaction receipt.
type SendRawTransactionSyncReturnType = *formatters.TransactionReceipt

// TransactionReceiptRevertedError is returned when a transaction receipt has a reverted status.
type TransactionReceiptRevertedError struct {
	Receipt *formatters.TransactionReceipt
}

func (e *TransactionReceiptRevertedError) Error() string {
	return fmt.Sprintf("transaction reverted (hash: %s)", e.Receipt.TransactionHash)
}

// SendRawTransactionSync sends a signed transaction to the network synchronously,
// and waits for the transaction to be included in a block.
//
// This is equivalent to viem's `sendRawTransactionSync` action.
//
// JSON-RPC Method: eth_sendRawTransactionSync (EIP-7966)
//
// Example:
//
//	receipt, err := wallet.SendRawTransactionSync(ctx, client, wallet.SendRawTransactionSyncParameters{
//	    SerializedTransaction: "0x02f850018203118080825208808080c080a04012522854168b27e5dc3d5839bab5e6b39e1a0ffd343901ce1622e3d64b48f1a04e00902ae0502c4728cbf12156290df99c3ed7de85b1dbfe20b5c36931733a33",
//	})
func SendRawTransactionSync(ctx context.Context, client Client, params SendRawTransactionSyncParameters) (SendRawTransactionSyncReturnType, error) {
	// Build RPC params - add timeout as hex if provided (mirrors viem's conditional params)
	var rpcParams []any
	if params.Timeout != nil {
		timeoutHex := encoding.NumberToHex(big.NewInt(*params.Timeout))
		rpcParams = []any{params.SerializedTransaction, timeoutHex}
	} else {
		rpcParams = []any{params.SerializedTransaction}
	}

	resp, err := client.Request(ctx, "eth_sendRawTransactionSync", rpcParams...)
	if err != nil {
		return nil, fmt.Errorf("eth_sendRawTransactionSync failed: %w", err)
	}

	// Parse the raw receipt from the RPC response
	var rpcReceipt formatters.RpcTransactionReceipt
	if unmarshalErr := json.Unmarshal(resp.Result, &rpcReceipt); unmarshalErr != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction receipt: %w", unmarshalErr)
	}

	// Format the receipt (mirrors viem's formatTransactionReceipt)
	receipt := formatters.FormatTransactionReceipt(rpcReceipt)

	// Check if receipt is reverted (mirrors viem's throwOnReceiptRevert check)
	throwOnRevert := true
	if params.ThrowOnReceiptRevert != nil {
		throwOnRevert = *params.ThrowOnReceiptRevert
	}
	if receipt.Status == formatters.ReceiptStatusReverted && throwOnRevert {
		return nil, &TransactionReceiptRevertedError{Receipt: &receipt}
	}

	return &receipt, nil
}
