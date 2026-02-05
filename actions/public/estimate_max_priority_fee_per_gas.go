package public

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ChefBingbong/viem-go/types"
)

// EstimateMaxPriorityFeePerGasReturnType is the return type for the
// EstimateMaxPriorityFeePerGas action.
// It represents the max priority fee per gas in wei.
type EstimateMaxPriorityFeePerGasReturnType = *big.Int

// EstimateMaxPriorityFeePerGasParameters contains the parameters for the
// EstimateMaxPriorityFeePerGas action.
//
// This mirrors viem's EstimateMaxPriorityFeePerGas parameters shape, but is
// simplified for viem-go since chain fee configuration is not yet exposed on
// the Chain type. The action will:
//   - Prefer the `eth_maxPriorityFeePerGas` RPC method when available.
//   - Fallback to `gasPrice - baseFeePerGas` using `eth_getBlockByNumber`
//     and `eth_gasPrice` when the RPC method is not supported.
type EstimateMaxPriorityFeePerGasParameters struct {
	// Block is an optional pre-fetched block to use for fallback
	// calculations. If nil, the latest block will be fetched when needed.
	Block *types.Block
}

// EstimateMaxPriorityFeePerGas returns an estimate for the max priority fee
// per gas (in wei) for a transaction to be likely included in the next block.
//
// This is equivalent to viem's `estimateMaxPriorityFeePerGas` action.
//
// JSON-RPC Methods:
//   - eth_maxPriorityFeePerGas (preferred)
//   - eth_getBlockByNumber + eth_gasPrice (fallback)
func EstimateMaxPriorityFeePerGas(
	ctx context.Context,
	client Client,
	params EstimateMaxPriorityFeePerGasParameters,
) (EstimateMaxPriorityFeePerGasReturnType, error) {
	// First, try the direct RPC method.
	feeHex, err := estimateMaxPriorityFeePerGasViaRpc(ctx, client)
	if err == nil {
		priorityFee, decodeErr := hexutil.DecodeBig(feeHex)
		if decodeErr != nil {
			return nil, fmt.Errorf("failed to decode maxPriorityFeePerGas: %w", decodeErr)
		}
		return priorityFee, nil
	}

	// Fallback: compute maxPriorityFeePerGas as gasPrice - baseFeePerGas.
	block := params.Block
	if block == nil {
		blockResult, blockErr := GetBlock(ctx, client, GetBlockParameters{
			BlockTag: BlockTagLatest,
		})
		if blockErr != nil {
			return nil, fmt.Errorf("failed to fetch latest block: %w", blockErr)
		}
		block = blockResult
	}

	if block.BaseFeePerGas == nil {
		return nil, fmt.Errorf("EIP-1559 fees not supported: missing baseFeePerGas on block")
	}

	gasPrice, gasPriceErr := GetGasPrice(ctx, client)
	if gasPriceErr != nil {
		return nil, fmt.Errorf("failed to fetch gas price: %w", gasPriceErr)
	}

	priorityFee := new(big.Int).Sub(gasPrice, block.BaseFeePerGas)
	if priorityFee.Sign() < 0 {
		// Never return a negative fee.
		return big.NewInt(0), nil
	}

	return priorityFee, nil
}

// estimateMaxPriorityFeePerGasViaRpc attempts to fetch the priority fee using
// the eth_maxPriorityFeePerGas RPC method.
func estimateMaxPriorityFeePerGasViaRpc(
	ctx context.Context,
	client Client,
) (string, error) {
	resp, err := client.Request(ctx, "eth_maxPriorityFeePerGas")
	if err != nil {
		return "", err
	}

	var feeHex string
	if unmarshalErr := json.Unmarshal(resp.Result, &feeHex); unmarshalErr != nil {
		return "", fmt.Errorf("failed to unmarshal maxPriorityFeePerGas: %w", unmarshalErr)
	}

	return feeHex, nil
}
