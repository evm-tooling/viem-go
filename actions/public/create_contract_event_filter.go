package public

import (
	"context"
	"fmt"

	json "github.com/goccy/go-json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ChefBingbong/viem-go/abi"
)

// CreateContractEventFilterParameters contains the parameters for the CreateContractEventFilter action.
// This mirrors viem's CreateContractEventFilterParameters type.
type CreateContractEventFilterParameters struct {
	// Address is the contract address(es) to filter logs from.
	// Can be a single address or a slice of addresses.
	Address any // common.Address or []common.Address

	// ABI is the contract ABI as JSON bytes, string, or *abi.ABI.
	ABI any

	// EventName filters logs to a specific event. If omitted, all events
	// from the ABI are matched.
	EventName string

	// Args filters logs by indexed event parameters.
	// Provide indexed args in order (use nil for "any" match).
	Args []any

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

	// Strict mode filters out logs that don't match the indexed/non-indexed
	// arguments in the event ABI. Default is false.
	Strict bool
}

// CreateContractEventFilterReturnType is the return type for the CreateContractEventFilter action.
// It contains the filter ID and metadata needed for subsequent operations.
type CreateContractEventFilterReturnType struct {
	// ID is the filter identifier.
	ID FilterID

	// Type indicates this is an event filter.
	Type string

	// ABI is the parsed ABI for decoding logs later.
	ABI *abi.ABI

	// EventName is the filtered event name (empty if filtering all events).
	EventName string

	// Args are the indexed args used for filtering.
	Args []any

	// Strict indicates whether strict mode is enabled.
	Strict bool
}

// CreateContractEventFilter creates a filter to retrieve contract event logs that can be used
// with GetFilterChanges or GetFilterLogs.
//
// This is equivalent to viem's `createContractEventFilter` action.
//
// JSON-RPC Method: eth_newFilter
//
// Example:
//
//	filter, err := public.CreateContractEventFilter(ctx, client, public.CreateContractEventFilterParameters{
//	    Address:   contractAddress,
//	    ABI:       erc20ABI,
//	    EventName: "Transfer",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// Use filter.ID with GetFilterChanges or GetFilterLogs
//
//	// Filter with indexed args
//	filter, err := public.CreateContractEventFilter(ctx, client, public.CreateContractEventFilterParameters{
//	    Address:   contractAddress,
//	    ABI:       erc20ABI,
//	    EventName: "Transfer",
//	    Args:      []any{fromAddress, nil}, // from=specific, to=any
//	})
func CreateContractEventFilter(ctx context.Context, client Client, params CreateContractEventFilterParameters) (*CreateContractEventFilterReturnType, error) {
	// Parse the ABI
	parsedABI, err := parseABIParam(params.ABI)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	// Build topics from event name and args
	var topics []any

	if params.EventName != "" {
		// Encode topics for specific event
		eventTopics, topicsErr := encodeEventTopicsForFilter(parsedABI, params.EventName, params.Args)
		if topicsErr != nil {
			return nil, fmt.Errorf("failed to encode event topics: %w", topicsErr)
		}
		topics = eventTopics
	}

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
	if len(topics) > 0 {
		encodedTopics := make([]any, len(topics))
		for i, topic := range topics {
			encodedTopics[i] = encodeFilterTopic(topic)
		}
		filterParams.Topics = encodedTopics
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

	return &CreateContractEventFilterReturnType{
		ID:        FilterID(filterID),
		Type:      "event",
		ABI:       parsedABI,
		EventName: params.EventName,
		Args:      params.Args,
		Strict:    params.Strict,
	}, nil
}

// GetRequest returns a function that can be used to make requests with the filter.
// This is useful for compatibility with the viem pattern where filters have a request method.
func (f *CreateContractEventFilterReturnType) GetRequest(client Client) func(ctx context.Context, method string, params ...any) (json.RawMessage, error) {
	return func(ctx context.Context, method string, params ...any) (json.RawMessage, error) {
		resp, err := client.Request(ctx, method, params...)
		if err != nil {
			return nil, err
		}
		return resp.Result, nil
	}
}
