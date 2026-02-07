package public

import (
	"context"
	"fmt"

	json "github.com/goccy/go-json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// CreateEventFilterParameters contains the parameters for the CreateEventFilter action.
// This mirrors viem's CreateEventFilterParameters type.
type CreateEventFilterParameters struct {
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
}

// FilterID represents a filter identifier returned by eth_newFilter.
type FilterID string

// CreateEventFilterReturnType is the return type for the CreateEventFilter action.
type CreateEventFilterReturnType struct {
	// ID is the filter identifier.
	ID FilterID

	// Type indicates this is an event filter.
	Type string
}

// rpcFilterParams is the RPC format for filter parameters.
type rpcFilterParams struct {
	Address   any    `json:"address,omitempty"`
	Topics    []any  `json:"topics,omitempty"`
	FromBlock string `json:"fromBlock,omitempty"`
	ToBlock   string `json:"toBlock,omitempty"`
}

// CreateEventFilter creates a filter object to receive logs matching the specified criteria.
//
// This is equivalent to viem's `createEventFilter` action.
//
// JSON-RPC Method: eth_newFilter
//
// Example:
//
//	filter, err := public.CreateEventFilter(ctx, client, public.CreateEventFilterParameters{
//	    Address: contractAddress,
//	    Topics:  []any{transferEventTopic},
//	})
func CreateEventFilter(ctx context.Context, client Client, params CreateEventFilterParameters) (*CreateEventFilterReturnType, error) {
	// Build filter params
	filterParams := rpcFilterParams{}

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

	// Execute the request
	resp, err := client.Request(ctx, "eth_newFilter", filterParams)
	if err != nil {
		return nil, fmt.Errorf("eth_newFilter failed: %w", err)
	}

	// Parse the filter ID
	var filterID string
	if err := json.Unmarshal(resp.Result, &filterID); err != nil {
		return nil, fmt.Errorf("failed to unmarshal filter ID: %w", err)
	}

	return &CreateEventFilterReturnType{
		ID:   FilterID(filterID),
		Type: "event",
	}, nil
}

// encodeFilterTopic encodes a topic value for the filter.
func encodeFilterTopic(topic any) any {
	if topic == nil {
		return nil
	}

	switch t := topic.(type) {
	case common.Hash:
		return t.Hex()
	case *common.Hash:
		if t != nil {
			return t.Hex()
		}
		return nil
	case [32]byte:
		return common.BytesToHash(t[:]).Hex()
	case []byte:
		if len(t) == 32 {
			return common.BytesToHash(t).Hex()
		}
		return hexutil.Encode(t)
	case string:
		return t
	case []common.Hash:
		// Array of topics (OR condition)
		result := make([]string, len(t))
		for i, h := range t {
			result[i] = h.Hex()
		}
		return result
	case []string:
		return t
	case []any:
		// Array of mixed topics
		result := make([]any, len(t))
		for i, item := range t {
			result[i] = encodeFilterTopic(item)
		}
		return result
	default:
		return t
	}
}
