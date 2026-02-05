package public

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/types"
)

// ReplacementReason indicates why a transaction was replaced.
type ReplacementReason string

const (
	// ReplacementReasonCancelled indicates the transaction was canceled (value === 0, sent to self).
	ReplacementReasonCancelled ReplacementReason = "canceled"
	// ReplacementReasonReplaced indicates the transaction was replaced with a different transaction.
	ReplacementReasonReplaced ReplacementReason = "replaced"
	// ReplacementReasonRepriced indicates the transaction was repriced (same tx, different gas).
	ReplacementReasonRepriced ReplacementReason = "repriced"
)

// ReplacementInfo contains information about a replaced transaction.
type ReplacementInfo struct {
	Reason              ReplacementReason
	ReplacedTransaction *TransactionResponse
	Transaction         *TransactionResponse
	TransactionReceipt  *types.Receipt
}

// WaitForTransactionReceiptParameters contains the parameters for the WaitForTransactionReceipt action.
// This mirrors viem's WaitForTransactionReceiptParameters type.
type WaitForTransactionReceiptParameters struct {
	// Hash is the hash of the transaction to wait for. Required.
	Hash common.Hash

	// CheckReplacement indicates whether to check for transaction replacements.
	// Default: true
	CheckReplacement *bool

	// Confirmations is the number of confirmations (blocks that have passed) to wait before resolving.
	// Default: 1
	Confirmations uint64

	// OnReplaced is an optional callback to emit if the transaction has been replaced.
	OnReplaced func(info ReplacementInfo)

	// PollingInterval is the polling frequency (in duration).
	// Default: 4 seconds
	PollingInterval time.Duration

	// RetryCount is the number of times to retry if the transaction or block is not found.
	// Default: 6
	RetryCount int

	// RetryDelay is a function that returns the delay between retries.
	// Default: exponential backoff: (1 << count) * 200ms
	RetryDelay func(count int) time.Duration

	// Timeout is the maximum time to wait before stopping polling.
	// Default: 180 seconds
	Timeout time.Duration
}

// WaitForTransactionReceiptReturnType is the return type for the WaitForTransactionReceipt action.
type WaitForTransactionReceiptReturnType = *types.Receipt

// WaitForTransactionReceiptTimeoutError is returned when waiting for a transaction receipt times out.
type WaitForTransactionReceiptTimeoutError struct {
	Hash common.Hash
}

func (e *WaitForTransactionReceiptTimeoutError) Error() string {
	return fmt.Sprintf("timed out waiting for transaction receipt: hash=%s", e.Hash.Hex())
}

// WaitForTransactionReceipt waits for the transaction to be included on a block (one confirmation),
// and then returns the transaction receipt.
//
// This is equivalent to viem's `waitForTransactionReceipt` action.
//
// The function additionally supports replacement detection (e.g., sped up transactions).
// Transactions can be replaced when a user modifies their transaction in their wallet
// (to speed up or cancel). Transactions are replaced when they are sent from the same nonce.
//
// There are 3 types of transaction replacement reasons:
//   - repriced: The gas price has been modified (e.g., different maxFeePerGas)
//   - canceled: The transaction has been canceled (e.g., value === 0, sent to self)
//   - replaced: The transaction has been replaced (e.g., different value or data)
//
// JSON-RPC Methods:
//   - Polls eth_getTransactionReceipt on each block until it has been processed.
//   - If a transaction has been replaced, calls eth_getBlockByNumber to find the replacement.
//
// Example:
//
//	receipt, err := public.WaitForTransactionReceipt(ctx, client, public.WaitForTransactionReceiptParameters{
//	    Hash: txHash,
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Transaction mined in block %d\n", receipt.BlockNumber)
func WaitForTransactionReceipt(ctx context.Context, client Client, params WaitForTransactionReceiptParameters) (WaitForTransactionReceiptReturnType, error) {
	// Set defaults
	checkReplacement := true
	if params.CheckReplacement != nil {
		checkReplacement = *params.CheckReplacement
	}

	confirmations := params.Confirmations
	if confirmations == 0 {
		confirmations = 1
	}

	pollingInterval := params.PollingInterval
	if pollingInterval == 0 {
		pollingInterval = 4 * time.Second
	}

	retryCount := params.RetryCount
	if retryCount == 0 {
		retryCount = 6
	}

	retryDelay := params.RetryDelay
	if retryDelay == nil {
		retryDelay = func(count int) time.Duration {
			return time.Duration(1<<count) * 200 * time.Millisecond
		}
	}

	timeout := params.Timeout
	if timeout == 0 {
		timeout = 180 * time.Second
	}

	// Create timeout context
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var transaction *TransactionResponse
	var receipt *types.Receipt

	// Try to get the receipt immediately
	receipt, _ = GetTransactionReceipt(ctx, client, GetTransactionReceiptParameters{
		Hash: params.Hash,
	})

	if receipt != nil && confirmations <= 1 {
		return receipt, nil
	}

	// Poll for the receipt
	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-timeoutCtx.Done():
			if errors.Is(timeoutCtx.Err(), context.DeadlineExceeded) {
				return nil, &WaitForTransactionReceiptTimeoutError{Hash: params.Hash}
			}
			return nil, timeoutCtx.Err()

		case <-ticker.C:
			// Get current block number
			blockNumber, err := GetBlockNumber(ctx, client, GetBlockNumberParameters{})
			if err != nil {
				continue // Retry on next tick
			}

			// If we already have a valid receipt, check confirmations
			if receipt != nil {
				if confirmations > 1 {
					if blockNumber-receipt.BlockNumber+1 < confirmations {
						continue // Not enough confirmations yet
					}
				}
				return receipt, nil
			}

			// Try to get the transaction if we need to check for replacement
			if checkReplacement && transaction == nil {
				transaction, _ = getTransactionWithRetry(ctx, client, params.Hash, retryCount, retryDelay)
			}

			// Try to get the receipt
			receipt, err = GetTransactionReceipt(ctx, client, GetTransactionReceiptParameters{
				Hash: params.Hash,
			})

			if err == nil && receipt != nil {
				// Check confirmations
				if confirmations > 1 {
					if blockNumber-receipt.BlockNumber+1 < confirmations {
						continue // Not enough confirmations yet
					}
				}
				return receipt, nil
			}

			// Receipt not found - check for replacement
			var receiptNotFoundErr *TransactionReceiptNotFoundError
			var txNotFoundErr *TransactionNotFoundError
			if errors.As(err, &receiptNotFoundErr) || errors.As(err, &txNotFoundErr) {
				if transaction == nil {
					continue // No transaction to check for replacement
				}

				// Try to find a replacement transaction in the current block
				replacement, replacementReceipt, reason := findReplacementTransaction(ctx, client, transaction, blockNumber, retryCount, retryDelay)
				if replacement != nil && replacementReceipt != nil {
					// Check confirmations for replacement
					if confirmations > 1 {
						if blockNumber-replacementReceipt.BlockNumber+1 < confirmations {
							continue // Not enough confirmations yet
						}
					}

					// Call the onReplaced callback if provided
					if params.OnReplaced != nil {
						params.OnReplaced(ReplacementInfo{
							Reason:              reason,
							ReplacedTransaction: transaction,
							Transaction:         replacement,
							TransactionReceipt:  replacementReceipt,
						})
					}

					return replacementReceipt, nil
				}
			}
		}
	}
}

// getTransactionWithRetry attempts to get a transaction with retries.
func getTransactionWithRetry(ctx context.Context, client Client, hash common.Hash, retryCount int, retryDelay func(int) time.Duration) (*TransactionResponse, error) {
	var lastErr error
	for i := 0; i < retryCount; i++ {
		tx, err := GetTransaction(ctx, client, GetTransactionParameters{
			Hash: &hash,
		})
		if err == nil {
			return tx, nil
		}
		lastErr = err

		// Wait before retry
		if i < retryCount-1 {
			time.Sleep(retryDelay(i))
		}
	}
	return nil, lastErr
}

// findReplacementTransaction looks for a replacement transaction in the given block.
func findReplacementTransaction(
	ctx context.Context,
	client Client,
	originalTx *TransactionResponse,
	blockNumber uint64,
	retryCount int,
	retryDelay func(int) time.Duration,
) (*TransactionResponse, *types.Receipt, ReplacementReason) {
	// Get the block with full transactions
	var block *types.Block
	var lastErr error

	for i := 0; i < retryCount; i++ {
		block, lastErr = GetBlock(ctx, client, GetBlockParameters{
			BlockNumber:         &blockNumber,
			IncludeTransactions: true,
		})
		if lastErr == nil {
			break
		}

		// Check if it's a block not found error
		var blockNotFoundErr *BlockNotFoundError
		if !errors.As(lastErr, &blockNotFoundErr) {
			return nil, nil, ""
		}

		if i < retryCount-1 {
			time.Sleep(retryDelay(i))
		}
	}

	if block == nil {
		return nil, nil, ""
	}

	// Look for a transaction with the same from address and nonce
	// Note: When IncludeTransactions is true, we need to iterate differently
	// For now, we'll fetch transactions individually based on the transaction hashes
	for _, txHash := range block.Transactions {
		tx, err := GetTransaction(ctx, client, GetTransactionParameters{
			Hash: &txHash,
		})
		if err != nil {
			continue
		}

		// Check if this is a replacement (same from and nonce)
		if tx.From == originalTx.From && tx.Nonce == originalTx.Nonce && tx.Hash != originalTx.Hash {
			// Found a replacement - get its receipt
			receipt, err := GetTransactionReceipt(ctx, client, GetTransactionReceiptParameters{
				Hash: tx.Hash,
			})
			if err != nil {
				continue
			}

			// Determine the replacement reason
			reason := determineReplacementReason(originalTx, tx)
			return tx, receipt, reason
		}
	}

	return nil, nil, ""
}

// determineReplacementReason determines why a transaction was replaced.
func determineReplacementReason(original, replacement *TransactionResponse) ReplacementReason {
	// Same to, value, and input means repriced (only gas changed)
	sameValue := false
	if original.Value != nil && replacement.Value != nil {
		sameValue = original.Value.Cmp(replacement.Value) == 0
	} else if original.Value == nil && replacement.Value == nil {
		sameValue = true
	}

	sameInput := bytesEqual(original.Input, replacement.Input)

	toEqual := false
	if original.To != nil && replacement.To != nil {
		toEqual = *original.To == *replacement.To
	} else if original.To == nil && replacement.To == nil {
		toEqual = true
	}

	if toEqual && sameValue && sameInput {
		return ReplacementReasonRepriced
	}

	// Sent to self with zero value means canceled
	zeroValue := replacement.Value == nil || replacement.Value.Cmp(big.NewInt(0)) == 0
	sentToSelf := replacement.To != nil && replacement.From == *replacement.To

	if sentToSelf && zeroValue {
		return ReplacementReasonCancelled
	}

	// Otherwise it's a full replacement
	return ReplacementReasonReplaced
}

// bytesEqual compares two byte slices for equality.
func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
