package public

import (
	"context"
	"fmt"

	json "github.com/goccy/go-json"
)

// CreatePendingTransactionFilterReturnType is the return type for the CreatePendingTransactionFilter action.
type CreatePendingTransactionFilterReturnType struct {
	// ID is the filter identifier.
	ID FilterID

	// Type indicates this is a pending transaction filter.
	Type string
}

// CreatePendingTransactionFilter creates a filter to listen for new pending transactions.
//
// This is equivalent to viem's `createPendingTransactionFilter` action.
//
// JSON-RPC Method: eth_newPendingTransactionFilter
//
// Example:
//
//	filter, err := public.CreatePendingTransactionFilter(ctx, client)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// Use filter.ID with GetFilterChanges to get new pending transaction hashes
func CreatePendingTransactionFilter(ctx context.Context, client Client) (*CreatePendingTransactionFilterReturnType, error) {
	// Execute the request
	resp, err := client.Request(ctx, "eth_newPendingTransactionFilter")
	if err != nil {
		return nil, fmt.Errorf("eth_newPendingTransactionFilter failed: %w", err)
	}

	// Parse the filter ID
	var filterID string
	if err := json.Unmarshal(resp.Result, &filterID); err != nil {
		return nil, fmt.Errorf("failed to unmarshal filter ID: %w", err)
	}

	return &CreatePendingTransactionFilterReturnType{
		ID:   FilterID(filterID),
		Type: "transaction",
	}, nil
}
