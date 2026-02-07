package public

import (
	"context"
	"fmt"

	json "github.com/goccy/go-json"
)

// GetChainIDReturnType is the return type for the GetChainID action.
// It represents the chain ID as a uint64.
type GetChainIDReturnType = uint64

// GetChainID returns the chain ID associated with the current network.
//
// This is equivalent to viem's `getChainId` action.
//
// JSON-RPC Method: eth_chainId
//
// Example:
//
//	chainID, err := public.GetChainID(ctx, client)
//	// chainID == 1 for Ethereum mainnet
func GetChainID(ctx context.Context, client Client) (GetChainIDReturnType, error) {
	resp, err := client.Request(ctx, "eth_chainId")
	if err != nil {
		return 0, fmt.Errorf("eth_chainId failed: %w", err)
	}

	var hexChainID string
	if unmarshalErr := json.Unmarshal(resp.Result, &hexChainID); unmarshalErr != nil {
		return 0, fmt.Errorf("failed to unmarshal chain id: %w", unmarshalErr)
	}

	chainID, parseErr := parseHexUint64(hexChainID)
	if parseErr != nil {
		return 0, fmt.Errorf("failed to parse chain id: %w", parseErr)
	}

	return chainID, nil
}
