package wallet

import (
	"context"
	"fmt"

	json "github.com/goccy/go-json"
)

// SendRawTransactionParameters contains the parameters for the SendRawTransaction action.
// This mirrors viem's SendRawTransactionParameters type.
type SendRawTransactionParameters struct {
	// SerializedTransaction is the signed serialized transaction hex string.
	SerializedTransaction string
}

// SendRawTransactionReturnType is the return type for the SendRawTransaction action.
// It is the transaction hash as a hex string.
type SendRawTransactionReturnType = string

// SendRawTransaction sends a signed transaction to the network.
//
// This is equivalent to viem's `sendRawTransaction` action.
//
// JSON-RPC Method: eth_sendRawTransaction
//
// Example:
//
//	hash, err := wallet.SendRawTransaction(ctx, client, wallet.SendRawTransactionParameters{
//	    SerializedTransaction: "0x02f850018203118080825208808080c080a04012522854168b27e5dc3d5839bab5e6b39e1a0ffd343901ce1622e3d64b48f1a04e00902ae0502c4728cbf12156290df99c3ed7de85b1dbfe20b5c36931733a33",
//	})
func SendRawTransaction(ctx context.Context, client Client, params SendRawTransactionParameters) (SendRawTransactionReturnType, error) {
	resp, err := client.Request(ctx, "eth_sendRawTransaction", params.SerializedTransaction)
	if err != nil {
		return "", fmt.Errorf("eth_sendRawTransaction failed: %w", err)
	}

	var hash string
	if unmarshalErr := json.Unmarshal(resp.Result, &hash); unmarshalErr != nil {
		return "", fmt.Errorf("failed to unmarshal transaction hash: %w", unmarshalErr)
	}

	return hash, nil
}
