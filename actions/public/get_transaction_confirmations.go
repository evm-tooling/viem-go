package public

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/types"
)

// GetTransactionConfirmationsParameters contains the parameters for the GetTransactionConfirmations action.
// This mirrors viem's GetTransactionConfirmationsParameters type.
//
// You must provide either:
//   - Hash: to get confirmations by transaction hash
//   - TransactionReceipt: to get confirmations using an existing receipt
type GetTransactionConfirmationsParameters struct {
	// Hash is the hash of the transaction. Optional if TransactionReceipt is provided.
	Hash *common.Hash

	// TransactionReceipt is an existing transaction receipt. Optional if Hash is provided.
	TransactionReceipt *types.Receipt
}

// GetTransactionConfirmationsReturnType is the return type for the GetTransactionConfirmations action.
// It represents the number of blocks passed since the transaction was processed.
type GetTransactionConfirmationsReturnType = uint64

// GetTransactionConfirmations returns the number of blocks passed (confirmations) since the transaction
// was processed on a block.
//
// This is equivalent to viem's `getTransactionConfirmations` action.
//
// Returns 0 if the transaction has not been confirmed & processed yet.
//
// Example:
//
//	// Get confirmations by hash
//	confirmations, err := public.GetTransactionConfirmations(ctx, client, public.GetTransactionConfirmationsParameters{
//	    Hash: &txHash,
//	})
//
//	// Get confirmations using an existing receipt
//	confirmations, err := public.GetTransactionConfirmations(ctx, client, public.GetTransactionConfirmationsParameters{
//	    TransactionReceipt: receipt,
//	})
func GetTransactionConfirmations(ctx context.Context, client Client, params GetTransactionConfirmationsParameters) (GetTransactionConfirmationsReturnType, error) {
	// Get current block number
	blockNumber, err := GetBlockNumber(ctx, client, GetBlockNumberParameters{})
	if err != nil {
		return 0, err
	}

	// Get transaction block number from receipt or by fetching the transaction
	var transactionBlockNumber *uint64

	if params.TransactionReceipt != nil {
		// Use the block number from the provided receipt
		bn := params.TransactionReceipt.BlockNumber
		transactionBlockNumber = &bn
	} else if params.Hash != nil {
		// Fetch the transaction to get its block number
		tx, err := GetTransaction(ctx, client, GetTransactionParameters{
			Hash: params.Hash,
		})
		if err != nil {
			return 0, err
		}
		transactionBlockNumber = tx.BlockNumber
	}

	// If the transaction hasn't been mined yet, return 0 confirmations
	if transactionBlockNumber == nil {
		return 0, nil
	}

	// Calculate confirmations: currentBlock - transactionBlock + 1
	return blockNumber - *transactionBlockNumber + 1, nil
}
