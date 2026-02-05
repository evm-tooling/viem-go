package public

import (
	"context"
	"encoding/json"
	"fmt"
)

// CreateBlockFilterReturnType is the return type for the CreateBlockFilter action.
type CreateBlockFilterReturnType struct {
	// ID is the filter identifier.
	ID FilterID

	// Type indicates this is a block filter.
	Type string
}

// CreateBlockFilter creates a filter to listen for new blocks.
//
// This is equivalent to viem's `createBlockFilter` action.
//
// JSON-RPC Method: eth_newBlockFilter
//
// Example:
//
//	filter, err := public.CreateBlockFilter(ctx, client)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// Use filter.ID with GetFilterChanges to get new block hashes
func CreateBlockFilter(ctx context.Context, client Client) (*CreateBlockFilterReturnType, error) {
	// Execute the request
	resp, err := client.Request(ctx, "eth_newBlockFilter")
	if err != nil {
		return nil, fmt.Errorf("eth_newBlockFilter failed: %w", err)
	}

	// Parse the filter ID
	var filterID string
	if err := json.Unmarshal(resp.Result, &filterID); err != nil {
		return nil, fmt.Errorf("failed to unmarshal filter ID: %w", err)
	}

	return &CreateBlockFilterReturnType{
		ID:   FilterID(filterID),
		Type: "block",
	}, nil
}
