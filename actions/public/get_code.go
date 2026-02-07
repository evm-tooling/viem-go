package public

import (
	"context"
	"fmt"

	json "github.com/goccy/go-json"

	"github.com/ethereum/go-ethereum/common"
)

// GetCodeParameters contains the parameters for the GetCode action.
// This mirrors viem's GetCodeParameters type.
type GetCodeParameters struct {
	// Address is the contract address to retrieve bytecode for.
	Address common.Address

	// BlockNumber is the block number to get the code at.
	// Mutually exclusive with BlockTag.
	BlockNumber *uint64

	// BlockTag is the block tag to get the code at (e.g., "latest", "pending").
	// Mutually exclusive with BlockNumber.
	// Default: "latest"
	BlockTag BlockTag
}

// GetCodeReturnType is the return type for the GetCode action.
// It represents the contract bytecode. Nil means no code (empty account).
type GetCodeReturnType = []byte

// GetCode retrieves the bytecode at an address.
//
// This is equivalent to viem's `getCode` action.
//
// JSON-RPC Method: eth_getCode
//
// Example:
//
//	code, err := public.GetCode(ctx, client, public.GetCodeParameters{
//	    Address: common.HexToAddress("0x..."),
//	})
//	if code == nil {
//	    // No contract deployed at this address
//	}
func GetCode(ctx context.Context, client Client, params GetCodeParameters) (GetCodeReturnType, error) {
	// Determine block tag/number
	blockTag := resolveBlockTag(client, params.BlockNumber, params.BlockTag)

	// Execute the request
	resp, err := client.Request(ctx, "eth_getCode", params.Address.Hex(), blockTag)
	if err != nil {
		return nil, fmt.Errorf("eth_getCode failed: %w", err)
	}

	var hexCode string
	if unmarshalErr := json.Unmarshal(resp.Result, &hexCode); unmarshalErr != nil {
		return nil, fmt.Errorf("failed to unmarshal code: %w", unmarshalErr)
	}

	// "0x" means no code (empty account)
	if hexCode == "" || hexCode == "0x" {
		return nil, nil
	}

	// Decode hex to bytes
	code := common.FromHex(hexCode)
	return code, nil
}
