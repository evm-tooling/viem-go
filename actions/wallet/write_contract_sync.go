package wallet

import (
	"context"
	"fmt"
	"time"

	"github.com/ChefBingbong/viem-go/utils/formatters"
)

// WriteContractSyncParameters contains the parameters for the WriteContractSync action.
// This mirrors viem's WriteContractSyncParameters type which extends WriteContractParameters
// with additional sync-specific fields (pollingInterval, throwOnReceiptRevert, timeout).
type WriteContractSyncParameters struct {
	WriteContractParameters

	// PollingInterval is the polling interval to poll for the transaction receipt.
	// Defaults to client.PollingInterval().
	PollingInterval time.Duration

	// ThrowOnReceiptRevert when true, throws an error if the transaction was detected as reverted.
	// Default: true.
	ThrowOnReceiptRevert *bool

	// Timeout is the timeout to wait for a response.
	// Default: max(chain.blockTime * 3, 5000)ms.
	Timeout *time.Duration
}

// WriteContractSyncReturnType is the return type for the WriteContractSync action.
// It is the formatted transaction receipt.
type WriteContractSyncReturnType = *formatters.TransactionReceipt

// WriteContractSync executes a write function on a contract synchronously.
// Returns the transaction receipt after the transaction is included in a block.
//
// A "write" function on a Solidity contract modifies the state of the blockchain.
// These types of functions require gas to be executed, and hence a Transaction is
// needed to be broadcast in order to change the state.
//
// Internally, encodes the function call using the ABI and delegates to SendTransactionSync
// with the ABI-encoded data.
//
// Warning: This internally sends a transaction – it does not validate if the contract
// write will succeed. It is highly recommended to simulate the contract write first.
//
// This is equivalent to viem's `writeContractSync` action.
//
// Example:
//
//	receipt, err := wallet.WriteContractSync(ctx, client, wallet.WriteContractSyncParameters{
//	    WriteContractParameters: wallet.WriteContractParameters{
//	        Address:      "0xFBA3912Ca04dd458c843e2EE08967fC04f3579c2",
//	        ABI:          mintABI,
//	        FunctionName: "mint",
//	        Args:         []any{uint32(69420)},
//	    },
//	})
func WriteContractSync(ctx context.Context, client Client, params WriteContractSyncParameters) (WriteContractSyncReturnType, error) {
	// Resolve account: param > client
	account := params.Account
	if account == nil {
		account = client.Account()
	}
	if account == nil {
		return nil, &AccountNotFoundError{DocsPath: "/docs/contract/writeContractSync"}
	}

	// Parse the ABI
	parsedABI, err := parseABIParam(params.ABI)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	// Encode function data (mirrors viem's encodeFunctionData({ abi, args, functionName }))
	calldata, encodeErr := parsedABI.EncodeFunctionData(params.FunctionName, params.Args...)
	if encodeErr != nil {
		return nil, wrapContractError(encodeErr, params.WriteContractParameters)
	}

	// Convert encoded calldata to hex string
	calldataHex := "0x" + fmt.Sprintf("%x", calldata)

	// Delegate to SendTransactionSync
	// This mirrors viem's writeContract.internal(client, sendTransactionSync, 'sendTransactionSync', parameters)
	receipt, txErr := SendTransactionSync(ctx, client, SendTransactionSyncParameters{
		Account:              account,
		Chain:                params.Chain,
		AssertChainID:        params.AssertChainID,
		DataSuffix:           params.DataSuffix,
		Data:                 calldataHex,
		To:                   params.Address,
		Value:                params.Value,
		AccessList:           params.AccessList,
		AuthorizationList:    params.AuthorizationList,
		BlobVersionedHashes:  params.BlobVersionedHashes,
		Blobs:                params.Blobs,
		Gas:                  params.Gas,
		GasPrice:             params.GasPrice,
		MaxFeePerBlobGas:     params.MaxFeePerBlobGas,
		MaxFeePerGas:         params.MaxFeePerGas,
		MaxPriorityFeePerGas: params.MaxPriorityFeePerGas,
		Nonce:                params.Nonce,
		Type:                 params.Type,
		PollingInterval:      params.PollingInterval,
		ThrowOnReceiptRevert: params.ThrowOnReceiptRevert,
		Timeout:              params.Timeout,
	})
	if txErr != nil {
		return nil, wrapContractError(txErr, params.WriteContractParameters)
	}

	return receipt, nil
}
