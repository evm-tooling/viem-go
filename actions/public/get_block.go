package public

import (
	"context"
	"fmt"

	json "github.com/goccy/go-json"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/types"
)

// GetBlockParameters contains the parameters for the GetBlock action.
// This mirrors viem's GetBlockParameters type.
type GetBlockParameters struct {
	// BlockHash is the hash of the block to retrieve.
	// Mutually exclusive with BlockNumber and BlockTag.
	BlockHash *common.Hash

	// BlockNumber is the number of the block to retrieve.
	// Mutually exclusive with BlockHash and BlockTag.
	BlockNumber *uint64

	// BlockTag is the block tag to retrieve (e.g., "latest", "pending").
	// Mutually exclusive with BlockHash and BlockNumber.
	// Default: "latest"
	BlockTag BlockTag

	// IncludeTransactions indicates whether to include full transaction objects
	// in the response. If false, only transaction hashes are included.
	// Default: false
	IncludeTransactions bool
}

// GetBlockReturnType is the return type for the GetBlock action.
type GetBlockReturnType = *types.Block

// GetBlock returns information about a block at a block number, hash, or tag.
//
// This is equivalent to viem's `getBlock` action.
//
// JSON-RPC Methods:
//   - eth_getBlockByHash for blockHash
//   - eth_getBlockByNumber for blockNumber & blockTag
//
// Example:
//
//	// Get latest block
//	block, err := public.GetBlock(ctx, client, public.GetBlockParameters{})
//
//	// Get block by number
//	blockNum := uint64(12345)
//	block, err := public.GetBlock(ctx, client, public.GetBlockParameters{
//	    BlockNumber: &blockNum,
//	})
//
//	// Get block by hash with transactions
//	block, err := public.GetBlock(ctx, client, public.GetBlockParameters{
//	    BlockHash:           &hash,
//	    IncludeTransactions: true,
//	})
func GetBlock(ctx context.Context, client Client, params GetBlockParameters) (GetBlockReturnType, error) {
	var result json.RawMessage
	var err error

	if params.BlockHash != nil {
		// Get block by hash
		resp, reqErr := client.Request(ctx, "eth_getBlockByHash", params.BlockHash.Hex(), params.IncludeTransactions)
		if reqErr != nil {
			return nil, fmt.Errorf("eth_getBlockByHash failed: %w", reqErr)
		}
		result = resp.Result
	} else {
		// Get block by number or tag
		blockTag := resolveBlockTag(client, params.BlockNumber, params.BlockTag)
		resp, reqErr := client.Request(ctx, "eth_getBlockByNumber", blockTag, params.IncludeTransactions)
		if reqErr != nil {
			return nil, fmt.Errorf("eth_getBlockByNumber failed: %w", reqErr)
		}
		result = resp.Result
	}

	// Check for null result (block not found)
	if result == nil || string(result) == "null" {
		return nil, &BlockNotFoundError{
			BlockHash:   params.BlockHash,
			BlockNumber: params.BlockNumber,
		}
	}

	// Parse the block
	var block types.Block
	if err = json.Unmarshal(result, &block); err != nil {
		return nil, fmt.Errorf("failed to unmarshal block: %w", err)
	}

	return &block, nil
}
