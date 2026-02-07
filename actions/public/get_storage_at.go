package public

import (
	"context"
	"fmt"

	json "github.com/goccy/go-json"

	"github.com/ethereum/go-ethereum/common"
)

// GetStorageAtParameters contains the parameters for the GetStorageAt action.
// This mirrors viem's GetStorageAtParameters type.
type GetStorageAtParameters struct {
	// Address is the contract/account address to read storage from.
	Address common.Address

	// Slot is the 32-byte storage slot key.
	Slot common.Hash

	// BlockNumber is the block number to read the storage at.
	// Mutually exclusive with BlockTag.
	BlockNumber *uint64

	// BlockTag is the block tag to read the storage at (e.g., "latest", "pending").
	// Mutually exclusive with BlockNumber.
	// Default: "latest"
	BlockTag BlockTag
}

// GetStorageAtReturnType is the return type for the GetStorageAt action.
// It represents the raw 32-byte storage value. Nil means zero/empty.
type GetStorageAtReturnType = []byte

// GetStorageAt returns the value from a storage slot at a given address.
//
// This is equivalent to viem's `getStorageAt` action.
//
// JSON-RPC Method: eth_getStorageAt
//
// Example:
//
//	value, err := public.GetStorageAt(ctx, client, public.GetStorageAtParameters{
//	    Address: common.HexToAddress("0x..."),
//	    Slot:    common.HexToHash("0x0"),
//	})
func GetStorageAt(ctx context.Context, client Client, params GetStorageAtParameters) (GetStorageAtReturnType, error) {
	// Determine block tag/number
	blockTag := resolveBlockTag(client, params.BlockNumber, params.BlockTag)

	// Execute the request
	resp, err := client.Request(ctx, "eth_getStorageAt", params.Address.Hex(), params.Slot.Hex(), blockTag)
	if err != nil {
		return nil, fmt.Errorf("eth_getStorageAt failed: %w", err)
	}

	var hexValue string
	if unmarshalErr := json.Unmarshal(resp.Result, &hexValue); unmarshalErr != nil {
		return nil, fmt.Errorf("failed to unmarshal storage value: %w", unmarshalErr)
	}

	// RPC returns 0x for zero value. Normalize to nil vs non-nil caller decision.
	if hexValue == "" || hexValue == "0x" {
		return nil, nil
	}

	value := common.FromHex(hexValue)
	return value, nil
}
