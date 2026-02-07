package public

import (
	"context"
	"fmt"

	json "github.com/goccy/go-json"
)

// UninstallFilterParameters contains the parameters for the UninstallFilter action.
type UninstallFilterParameters struct {
	// Filter is the filter ID to uninstall.
	Filter FilterID
}

// UninstallFilterReturnType is the return type for the UninstallFilter action.
type UninstallFilterReturnType = bool

// UninstallFilter destroys a filter and frees server resources.
//
// This is equivalent to viem's `uninstallFilter` action.
//
// JSON-RPC Method: eth_uninstallFilter
//
// Note: Filters timeout when they aren't requested with eth_getFilterChanges
// for a period of time, so this is mainly useful for cleanup when you're
// done with a filter early.
//
// Example:
//
//	success, err := public.UninstallFilter(ctx, client, filter.ID)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	if success {
//	    fmt.Println("Filter successfully uninstalled")
//	}
func UninstallFilter(ctx context.Context, client Client, filterID FilterID) (UninstallFilterReturnType, error) {
	// Execute the request
	resp, err := client.Request(ctx, "eth_uninstallFilter", string(filterID))
	if err != nil {
		return false, fmt.Errorf("eth_uninstallFilter failed: %w", err)
	}

	// Parse the result
	var success bool
	if err := json.Unmarshal(resp.Result, &success); err != nil {
		return false, fmt.Errorf("failed to unmarshal uninstall result: %w", err)
	}

	return success, nil
}
