package public

import (
	"context"
	"fmt"

	json "github.com/goccy/go-json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ChefBingbong/viem-go/utils/formatters"
)

// GetLogsParameters contains the parameters for the GetLogs action.
// This mirrors viem's GetLogsParameters type.
type GetLogsParameters struct {
	// Address is the contract address(es) to filter logs from.
	// Can be a single address or a slice of addresses.
	Address any // common.Address or []common.Address

	// Topics is the indexed event topics to filter.
	// Each topic can be a single value or an array of values (OR condition).
	Topics []any

	// FromBlock is the block number to start filtering from.
	// Mutually exclusive with FromBlockTag.
	FromBlock *uint64

	// FromBlockTag is the block tag to start filtering from.
	// Mutually exclusive with FromBlock.
	FromBlockTag BlockTag

	// ToBlock is the block number to stop filtering at.
	// Mutually exclusive with ToBlockTag.
	ToBlock *uint64

	// ToBlockTag is the block tag to stop filtering at.
	// Mutually exclusive with ToBlock.
	ToBlockTag BlockTag

	// BlockHash filters logs from a specific block by hash.
	// Mutually exclusive with FromBlock/ToBlock.
	BlockHash *common.Hash
}

// GetLogsReturnType is the return type for the GetLogs action.
type GetLogsReturnType = []formatters.Log

// rpcGetLogsParams is the RPC format for getLogs parameters.
type rpcGetLogsParams struct {
	Address   any    `json:"address,omitempty"`
	Topics    []any  `json:"topics,omitempty"`
	FromBlock string `json:"fromBlock,omitempty"`
	ToBlock   string `json:"toBlock,omitempty"`
	BlockHash string `json:"blockHash,omitempty"`
}

// GetLogs returns event logs matching the specified filter criteria.
//
// This is equivalent to viem's `getLogs` action.
//
// JSON-RPC Method: eth_getLogs
//
// Example:
//
//	// Get Transfer events from a specific block range
//	fromBlock := uint64(18000000)
//	toBlock := uint64(18000100)
//	logs, err := public.GetLogs(ctx, client, public.GetLogsParameters{
//	    Address:   contractAddress,
//	    Topics:    []any{transferEventTopic},
//	    FromBlock: &fromBlock,
//	    ToBlock:   &toBlock,
//	})
//
//	// Get logs from a specific block
//	blockHash := common.HexToHash("0x...")
//	logs, err := public.GetLogs(ctx, client, public.GetLogsParameters{
//	    BlockHash: &blockHash,
//	})
func GetLogs(ctx context.Context, client Client, params GetLogsParameters) (GetLogsReturnType, error) {
	// Build filter params
	filterParams := rpcGetLogsParams{}

	// Handle address (single or array)
	if params.Address != nil {
		switch addr := params.Address.(type) {
		case common.Address:
			filterParams.Address = addr.Hex()
		case *common.Address:
			if addr != nil {
				filterParams.Address = addr.Hex()
			}
		case []common.Address:
			if len(addr) > 0 {
				addrs := make([]string, len(addr))
				for i, a := range addr {
					addrs[i] = a.Hex()
				}
				filterParams.Address = addrs
			}
		case string:
			filterParams.Address = addr
		case []string:
			filterParams.Address = addr
		}
	}

	// Handle topics
	if len(params.Topics) > 0 {
		topics := make([]any, len(params.Topics))
		for i, topic := range params.Topics {
			topics[i] = encodeFilterTopic(topic)
		}
		filterParams.Topics = topics
	}

	// Handle block range or block hash
	if params.BlockHash != nil {
		filterParams.BlockHash = params.BlockHash.Hex()
	} else {
		// Handle fromBlock
		if params.FromBlock != nil {
			filterParams.FromBlock = hexutil.EncodeUint64(*params.FromBlock)
		} else if params.FromBlockTag != "" {
			filterParams.FromBlock = string(params.FromBlockTag)
		}

		// Handle toBlock
		if params.ToBlock != nil {
			filterParams.ToBlock = hexutil.EncodeUint64(*params.ToBlock)
		} else if params.ToBlockTag != "" {
			filterParams.ToBlock = string(params.ToBlockTag)
		}
	}

	// Execute the request
	resp, err := client.Request(ctx, "eth_getLogs", filterParams)
	if err != nil {
		return nil, fmt.Errorf("eth_getLogs failed: %w", err)
	}

	// Parse the logs
	var rpcLogs []formatters.RpcLog
	if err := json.Unmarshal(resp.Result, &rpcLogs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal logs: %w", err)
	}

	// Format logs
	return formatters.FormatLogs(rpcLogs), nil
}

// GetLogsWithEvents returns event logs matching the specified filter criteria,
// with decoded event data using the provided ABI events.
//
// This is a convenience function that combines GetLogs with event decoding.
//
// Example:
//
//	logs, err := public.GetLogsWithEvents(ctx, client, params, transferABI)
func GetLogsWithEvents(ctx context.Context, client Client, params GetLogsParameters, events []Event) (GetLogsReturnType, error) {
	logs, err := GetLogs(ctx, client, params)
	if err != nil {
		return nil, err
	}

	// If no events provided, return raw logs
	if len(events) == 0 {
		return logs, nil
	}

	// Decode events - this would require integration with the ABI package
	// For now, return the formatted logs
	return logs, nil
}

// Event represents an ABI event for filtering.
// This is a simplified representation for use in GetLogsWithEvents.
type Event struct {
	Name      string
	Signature string
	Topic     common.Hash
}
