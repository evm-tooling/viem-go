package public

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/abi"
	"github.com/ChefBingbong/viem-go/utils/formatters"
)

// GetContractEventsParameters contains the parameters for the GetContractEvents action.
// This mirrors viem's GetContractEventsParameters type.
type GetContractEventsParameters struct {
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

	// BlockHash filters logs from a specific block by hash.
	// Mutually exclusive with FromBlock/ToBlock.
	BlockHash *common.Hash

	// Strict mode filters out logs that don't match the indexed/non-indexed
	// arguments in the event ABI. Default is false.
	Strict bool
}

// ContractEventLog represents a formatted and decoded event log.
type ContractEventLog struct {
	formatters.Log

	// EventName is the decoded event name.
	EventName string `json:"eventName,omitempty"`

	// DecodedArgs contains the decoded event arguments as a map.
	DecodedArgs map[string]any `json:"args,omitempty"`
}

// GetContractEventsReturnType is the return type for the GetContractEvents action.
type GetContractEventsReturnType = []ContractEventLog

// GetContractEvents returns a list of event logs emitted by a contract.
//
// This is equivalent to viem's `getContractEvents` action.
//
// JSON-RPC Method: eth_getLogs
//
// Example:
//
//	logs, err := public.GetContractEvents(ctx, client, public.GetContractEventsParameters{
//	    Address:   contractAddress,
//	    ABI:       erc20ABI,
//	    EventName: "Transfer",
//	})
//
//	// Get all events from a contract
//	logs, err := public.GetContractEvents(ctx, client, public.GetContractEventsParameters{
//	    Address: contractAddress,
//	    ABI:     contractABI,
//	})
//
//	// Filter by indexed args
//	logs, err := public.GetContractEvents(ctx, client, public.GetContractEventsParameters{
//	    Address:   contractAddress,
//	    ABI:       erc20ABI,
//	    EventName: "Transfer",
//	    Args:      []any{fromAddress, nil}, // from=specific, to=any
//	})
func GetContractEvents(ctx context.Context, client Client, params GetContractEventsParameters) (GetContractEventsReturnType, error) {
	// Parse the ABI
	parsedABI, err := parseABIParam(params.ABI)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	// Build topics from event name and args
	var topics []any
	var targetEvent *abi.Event
	var allEvents []abi.Event

	if params.EventName != "" {
		// Get specific event
		event, eventErr := parsedABI.GetEvent(params.EventName)
		if eventErr != nil {
			return nil, fmt.Errorf("event %q not found in ABI: %w", params.EventName, eventErr)
		}
		targetEvent = event

		// Encode topics for this event
		eventTopics, topicsErr := encodeEventTopicsForFilter(parsedABI, params.EventName, params.Args)
		if topicsErr != nil {
			return nil, fmt.Errorf("failed to encode event topics: %w", topicsErr)
		}
		topics = eventTopics
	} else {
		// Get all events from ABI
		for _, event := range parsedABI.Events {
			allEvents = append(allEvents, event)
		}

		// If no specific event, we can optionally filter by all event signatures
		// For simplicity, we don't add topic filters when getting all events
	}

	// Build GetLogs parameters
	logsParams := GetLogsParameters{
		Address:      params.Address,
		Topics:       topics,
		FromBlock:    params.FromBlock,
		FromBlockTag: params.FromBlockTag,
		ToBlock:      params.ToBlock,
		ToBlockTag:   params.ToBlockTag,
		BlockHash:    params.BlockHash,
	}

	// Get logs
	logs, err := GetLogs(ctx, client, logsParams)
	if err != nil {
		return nil, err
	}

	// Parse and decode logs
	results := make(GetContractEventsReturnType, 0, len(logs))

	for _, log := range logs {
		if len(log.Topics) == 0 {
			continue
		}

		// Find matching event
		var matchedEvent *abi.Event
		topicHash := common.HexToHash(log.Topics[0])

		if targetEvent != nil {
			// We're filtering for a specific event
			if targetEvent.Topic == topicHash {
				matchedEvent = targetEvent
			}
		} else {
			// Search all events
			for i := range allEvents {
				if allEvents[i].Topic == topicHash {
					matchedEvent = &allEvents[i]
					break
				}
			}
		}

		if matchedEvent == nil {
			// Log doesn't match any known event
			if params.Strict {
				continue
			}
			// In non-strict mode, include without decoding
			results = append(results, ContractEventLog{
				Log: log,
			})
			continue
		}

		// Decode the event log
		decodedArgs, decodeErr := decodeEventLog(parsedABI, matchedEvent.Name, log)
		if decodeErr != nil {
			if params.Strict {
				continue
			}
			// In non-strict mode, include with just the event name
			results = append(results, ContractEventLog{
				Log:       log,
				EventName: matchedEvent.Name,
			})
			continue
		}

		results = append(results, ContractEventLog{
			Log:         log,
			EventName:   matchedEvent.Name,
			DecodedArgs: decodedArgs,
		})
	}

	return results, nil
}

// encodeEventTopicsForFilter encodes event topics for use in a log filter.
func encodeEventTopicsForFilter(parsedABI *abi.ABI, eventName string, args []any) ([]any, error) {
	event, err := parsedABI.GetEvent(eventName)
	if err != nil {
		return nil, err
	}

	// First topic is always the event signature (for non-anonymous events)
	topics := []any{event.Topic.Hex()}

	// No args to encode
	if len(args) == 0 {
		return topics, nil
	}

	// Encode indexed args
	indexedCount := 0
	for _, input := range event.Inputs {
		if input.Indexed {
			indexedCount++
		}
	}

	if len(args) > indexedCount {
		return nil, fmt.Errorf("event %q has %d indexed parameters but %d args provided",
			eventName, indexedCount, len(args))
	}

	// Encode each indexed arg
	argIdx := 0
	for _, input := range event.Inputs {
		if !input.Indexed {
			continue
		}
		if argIdx >= len(args) {
			break
		}

		arg := args[argIdx]
		argIdx++

		if arg == nil {
			// nil means "match any"
			topics = append(topics, nil)
			continue
		}

		// Encode the topic
		topic, encodeErr := encodeTopicValue(arg, input.Type)
		if encodeErr != nil {
			return nil, fmt.Errorf("failed to encode arg %d: %w", argIdx-1, encodeErr)
		}
		topics = append(topics, topic)
	}

	return topics, nil
}

// encodeTopicValue encodes a single value as a topic.
func encodeTopicValue(value any, typeName string) (string, error) {
	switch v := value.(type) {
	case common.Address:
		// Pad address to 32 bytes
		var topic [32]byte
		copy(topic[12:], v.Bytes())
		return common.BytesToHash(topic[:]).Hex(), nil

	case *common.Address:
		if v == nil {
			return "", fmt.Errorf("nil address pointer")
		}
		var topic [32]byte
		copy(topic[12:], v.Bytes())
		return common.BytesToHash(topic[:]).Hex(), nil

	case string:
		// Check if it's an address string
		if common.IsHexAddress(v) && (typeName == "address" || typeName == "") {
			addr := common.HexToAddress(v)
			var topic [32]byte
			copy(topic[12:], addr.Bytes())
			return common.BytesToHash(topic[:]).Hex(), nil
		}
		// For string type, would need to hash - not typically used in filters
		return "", fmt.Errorf("string topics require hashing, use hash directly")

	case common.Hash:
		return v.Hex(), nil

	case [32]byte:
		return common.BytesToHash(v[:]).Hex(), nil

	case []byte:
		if len(v) == 32 {
			return common.BytesToHash(v).Hex(), nil
		}
		return "", fmt.Errorf("bytes topic must be 32 bytes, got %d", len(v))

	default:
		return "", fmt.Errorf("unsupported topic type: %T", value)
	}
}

// decodeEventLog decodes a formatted log using the ABI.
func decodeEventLog(parsedABI *abi.ABI, eventName string, log formatters.Log) (map[string]any, error) {
	// Convert topics from strings to common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.HexToHash(t)
	}

	// Decode data
	data := common.FromHex(log.Data)

	// Use the ABI's DecodeEventLogByName
	decoded, err := parsedABI.DecodeEventLogByName(eventName, topics, data)
	if err != nil {
		return nil, err
	}

	return decoded.Args, nil
}
