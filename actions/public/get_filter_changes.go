package public

import (
	"context"
	"fmt"

	json "github.com/goccy/go-json"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/utils/formatters"
)

// GetFilterChangesParameters contains the parameters for the GetFilterChanges action.
type GetFilterChangesParameters struct {
	// Filter is the filter to get changes for.
	// This should be a filter returned by CreateEventFilter, CreateBlockFilter,
	// or CreatePendingTransactionFilter.
	Filter FilterID
}

// GetFilterChangesLogsReturnType is the return type when getting changes for an event filter.
type GetFilterChangesLogsReturnType = []formatters.Log

// GetFilterChangesBlocksReturnType is the return type when getting changes for a block filter.
type GetFilterChangesBlocksReturnType = []common.Hash

// GetFilterChangesTransactionsReturnType is the return type when getting changes for a pending transaction filter.
type GetFilterChangesTransactionsReturnType = []common.Hash

// GetFilterChangesLogs returns logs that have occurred since the last poll for an event filter.
//
// This is equivalent to viem's `getFilterChanges` action for event filters.
//
// JSON-RPC Method: eth_getFilterChanges
//
// Example:
//
//	logs, err := public.GetFilterChangesLogs(ctx, client, filter.ID)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, log := range logs {
//	    fmt.Printf("Log from %s at block %d\n", log.Address, log.BlockNumber)
//	}
func GetFilterChangesLogs(ctx context.Context, client Client, filterID FilterID) (GetFilterChangesLogsReturnType, error) {
	// Execute the request
	resp, err := client.Request(ctx, "eth_getFilterChanges", string(filterID))
	if err != nil {
		return nil, fmt.Errorf("eth_getFilterChanges failed: %w", err)
	}

	// Parse the logs
	var rpcLogs []formatters.RpcLog
	if err := json.Unmarshal(resp.Result, &rpcLogs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal filter changes (logs): %w", err)
	}

	// Format logs
	return formatters.FormatLogs(rpcLogs), nil
}

// GetFilterChangesBlocks returns block hashes that have occurred since the last poll for a block filter.
//
// This is equivalent to viem's `getFilterChanges` action for block filters.
//
// JSON-RPC Method: eth_getFilterChanges
//
// Example:
//
//	blockHashes, err := public.GetFilterChangesBlocks(ctx, client, filter.ID)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, hash := range blockHashes {
//	    fmt.Printf("New block: %s\n", hash.Hex())
//	}
func GetFilterChangesBlocks(ctx context.Context, client Client, filterID FilterID) (GetFilterChangesBlocksReturnType, error) {
	// Execute the request
	resp, err := client.Request(ctx, "eth_getFilterChanges", string(filterID))
	if err != nil {
		return nil, fmt.Errorf("eth_getFilterChanges failed: %w", err)
	}

	// Parse the block hashes
	var hexHashes []string
	if err := json.Unmarshal(resp.Result, &hexHashes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal filter changes (blocks): %w", err)
	}

	// Convert to common.Hash
	hashes := make([]common.Hash, len(hexHashes))
	for i, hexHash := range hexHashes {
		hashes[i] = common.HexToHash(hexHash)
	}

	return hashes, nil
}

// GetFilterChangesTransactions returns transaction hashes that have occurred since the last poll
// for a pending transaction filter.
//
// This is equivalent to viem's `getFilterChanges` action for pending transaction filters.
//
// JSON-RPC Method: eth_getFilterChanges
//
// Example:
//
//	txHashes, err := public.GetFilterChangesTransactions(ctx, client, filter.ID)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, hash := range txHashes {
//	    fmt.Printf("New pending tx: %s\n", hash.Hex())
//	}
func GetFilterChangesTransactions(ctx context.Context, client Client, filterID FilterID) (GetFilterChangesTransactionsReturnType, error) {
	// Execute the request
	resp, err := client.Request(ctx, "eth_getFilterChanges", string(filterID))
	if err != nil {
		return nil, fmt.Errorf("eth_getFilterChanges failed: %w", err)
	}

	// Parse the transaction hashes
	var hexHashes []string
	if err := json.Unmarshal(resp.Result, &hexHashes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal filter changes (transactions): %w", err)
	}

	// Convert to common.Hash
	hashes := make([]common.Hash, len(hexHashes))
	for i, hexHash := range hexHashes {
		hashes[i] = common.HexToHash(hexHash)
	}

	return hashes, nil
}

// GetFilterChangesRaw returns raw JSON results from a filter.
// This is useful when you need to handle the result type dynamically.
//
// JSON-RPC Method: eth_getFilterChanges
func GetFilterChangesRaw(ctx context.Context, client Client, filterID FilterID) (json.RawMessage, error) {
	resp, err := client.Request(ctx, "eth_getFilterChanges", string(filterID))
	if err != nil {
		return nil, fmt.Errorf("eth_getFilterChanges failed: %w", err)
	}
	return resp.Result, nil
}
