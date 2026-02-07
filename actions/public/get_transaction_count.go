package public

import (
	"context"
	"fmt"

	json "github.com/goccy/go-json"

	"github.com/ethereum/go-ethereum/common"
)

type GetTransactionCountParameters struct {
	// Address is the address to get the tx count of. Required.
	Address common.Address

	// BlockNumber is the block number to get the tx count at.
	// Mutually exclusive with BlockTag.
	BlockNumber *uint64

	// BlockTag is the block tag to get the tx count at (e.g., "latest", "pending").
	// Mutually exclusive with BlockNumber.
	BlockTag BlockTag
}

// GetTransactionCountReturnType is the return type for the GetTransactionCount action.
// It represents the tx count in wei.
type GetTransactionCountReturnType = uint64

func GetTransactionCount(ctx context.Context, client Client, params GetTransactionCountParameters) (GetTransactionCountReturnType, error) {
	// Determine block tag
	blockTag := resolveBlockTag(client, params.BlockNumber, params.BlockTag)

	// Execute the request
	resp, err := client.Request(ctx, "eth_getTransactionCount", params.Address.Hex(), blockTag)
	if err != nil {
		return 0, fmt.Errorf("eth_getTransactionCount failed: %w", err)
	}

	var txCountHex string
	if unmarshalErr := json.Unmarshal(resp.Result, &txCountHex); unmarshalErr != nil {
		return 0, fmt.Errorf("failed to unmarshal tx count: %w", unmarshalErr)
	}

	// Parse the tx count
	txCount, err := parseHexUint64(txCountHex)
	if err != nil {
		return 0, fmt.Errorf("failed to parse tx count: %w", err)
	}

	return txCount, nil
}
