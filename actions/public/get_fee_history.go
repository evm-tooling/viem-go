package public

import (
	"context"
	"fmt"

	json "github.com/goccy/go-json"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ChefBingbong/viem-go/utils/formatters"
)

// GetFeeHistoryParameters contains the parameters for the GetFeeHistory action.
// This mirrors viem's GetFeeHistoryParameters type.
type GetFeeHistoryParameters struct {
	// BlockCount is the number of blocks in the requested range.
	// Between 1 and 1024 blocks can be requested in a single query.
	BlockCount uint64

	// RewardPercentiles is a monotonically increasing list of percentile values
	// to sample from each block's effective priority fees per gas in ascending
	// order, weighted by gas used.
	RewardPercentiles []float64

	// BlockNumber is the highest block number of the requested range.
	// Mutually exclusive with BlockTag.
	BlockNumber *uint64

	// BlockTag is the highest block tag of the requested range (e.g., "latest").
	// Mutually exclusive with BlockNumber.
	// Default: "latest"
	BlockTag BlockTag
}

// GetFeeHistoryReturnType is the return type for the GetFeeHistory action.
// It represents a formatted fee history.
type GetFeeHistoryReturnType = formatters.FeeHistory

// GetFeeHistory returns a collection of historical gas information.
//
// This is equivalent to viem's `getFeeHistory` action.
//
// JSON-RPC Method: eth_feeHistory
//
// Example:
//
//	history, err := public.GetFeeHistory(ctx, client, public.GetFeeHistoryParameters{
//	    BlockCount:        4,
//	    RewardPercentiles: []float64{25, 75},
//	})
func GetFeeHistory(ctx context.Context, client Client, params GetFeeHistoryParameters) (GetFeeHistoryReturnType, error) {
	if params.BlockCount == 0 {
		return formatters.FeeHistory{}, fmt.Errorf("blockCount must be greater than 0")
	}

	// Encode blockCount as hex quantity
	blockCountHex := hexutil.EncodeUint64(params.BlockCount)

	// Determine newest block (hex number or tag)
	newestBlock := ""
	if params.BlockNumber != nil {
		newestBlock = hexutil.EncodeUint64(*params.BlockNumber)
	} else {
		newestBlock = resolveBlockTag(client, params.BlockNumber, params.BlockTag)
	}

	// Execute the request
	resp, err := client.Request(ctx, "eth_feeHistory", blockCountHex, newestBlock, params.RewardPercentiles)
	if err != nil {
		return formatters.FeeHistory{}, fmt.Errorf("eth_feeHistory failed: %w", err)
	}

	var rpcHistory formatters.RpcFeeHistory
	if unmarshalErr := json.Unmarshal(resp.Result, &rpcHistory); unmarshalErr != nil {
		return formatters.FeeHistory{}, fmt.Errorf("failed to unmarshal fee history: %w", unmarshalErr)
	}

	// Format RPC fee history into typed FeeHistory
	history := formatters.FormatFeeHistory(rpcHistory)
	return history, nil
}
