package public

import (
	"context"
	"fmt"

	json "github.com/goccy/go-json"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/abi"
	"github.com/ChefBingbong/viem-go/utils/formatters"
)

// GetFilterLogsParameters contains the parameters for the GetFilterLogs action.
// This mirrors viem's GetFilterLogsParameters type.
type GetFilterLogsParameters struct {
	// Filter is the filter to get logs for.
	// This can be a basic event filter or a contract event filter with ABI.
	Filter any // FilterID, *CreateEventFilterReturnType, or *CreateContractEventFilterReturnType
}

// GetFilterLogsReturnType is the return type for GetFilterLogs.
// When using a contract event filter with ABI, logs will include decoded args.
type GetFilterLogsReturnType = []ContractEventLog

// GetFilterLogs returns a list of event logs since the filter was created.
//
// This is equivalent to viem's `getFilterLogs` action.
//
// JSON-RPC Method: eth_getFilterLogs
//
// Note: `getFilterLogs` is only compatible with **event filters**.
//
// Example:
//
//	// With a basic event filter
//	filter, _ := public.CreateEventFilter(ctx, client, params)
//	logs, err := public.GetFilterLogs(ctx, client, public.GetFilterLogsParameters{
//	    Filter: filter.ID,
//	})
//
//	// With a contract event filter (includes decoded args)
//	filter, _ := public.CreateContractEventFilter(ctx, client, params)
//	logs, err := public.GetFilterLogs(ctx, client, public.GetFilterLogsParameters{
//	    Filter: filter,
//	})
func GetFilterLogs(ctx context.Context, client Client, params GetFilterLogsParameters) (GetFilterLogsReturnType, error) {
	// Extract filter ID and optional ABI info
	filterID, parsedABI, eventName, strict, err := extractFilterInfo(params.Filter)
	if err != nil {
		return nil, err
	}

	// Execute the request
	resp, err := client.Request(ctx, "eth_getFilterLogs", string(filterID))
	if err != nil {
		return nil, fmt.Errorf("eth_getFilterLogs failed: %w", err)
	}

	// Parse the logs
	var rpcLogs []formatters.RpcLog
	if err := json.Unmarshal(resp.Result, &rpcLogs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal filter logs: %w", err)
	}

	// Format logs
	formattedLogs := formatters.FormatLogs(rpcLogs)

	// If no ABI, return basic logs wrapped in ContractEventLog
	if parsedABI == nil {
		results := make(GetFilterLogsReturnType, len(formattedLogs))
		for i, log := range formattedLogs {
			results[i] = ContractEventLog{Log: log}
		}
		return results, nil
	}

	// Parse and decode logs with ABI
	return parseFilterLogs(formattedLogs, parsedABI, eventName, strict)
}

// GetFilterLogsWithFilter is a convenience function that takes a contract event filter directly.
func GetFilterLogsWithFilter(ctx context.Context, client Client, filter *CreateContractEventFilterReturnType) (GetFilterLogsReturnType, error) {
	return GetFilterLogs(ctx, client, GetFilterLogsParameters{
		Filter: filter,
	})
}

// GetFilterLogsRaw returns raw logs without decoding.
// Useful when you don't need decoded event args.
func GetFilterLogsRaw(ctx context.Context, client Client, filterID FilterID) ([]formatters.Log, error) {
	// Execute the request
	resp, err := client.Request(ctx, "eth_getFilterLogs", string(filterID))
	if err != nil {
		return nil, fmt.Errorf("eth_getFilterLogs failed: %w", err)
	}

	// Parse the logs
	var rpcLogs []formatters.RpcLog
	if err := json.Unmarshal(resp.Result, &rpcLogs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal filter logs: %w", err)
	}

	return formatters.FormatLogs(rpcLogs), nil
}

// extractFilterInfo extracts the filter ID and optional ABI info from the filter parameter.
func extractFilterInfo(filter any) (FilterID, *abi.ABI, string, bool, error) {
	switch f := filter.(type) {
	case FilterID:
		return f, nil, "", false, nil

	case string:
		return FilterID(f), nil, "", false, nil

	case *CreateEventFilterReturnType:
		if f == nil {
			return "", nil, "", false, fmt.Errorf("filter cannot be nil")
		}
		return f.ID, nil, "", false, nil

	case *CreateContractEventFilterReturnType:
		if f == nil {
			return "", nil, "", false, fmt.Errorf("filter cannot be nil")
		}
		return f.ID, f.ABI, f.EventName, f.Strict, nil

	case CreateEventFilterReturnType:
		return f.ID, nil, "", false, nil

	case CreateContractEventFilterReturnType:
		return f.ID, f.ABI, f.EventName, f.Strict, nil

	default:
		return "", nil, "", false, fmt.Errorf("unsupported filter type: %T", filter)
	}
}

// parseFilterLogs parses and decodes logs using the provided ABI.
func parseFilterLogs(logs []formatters.Log, parsedABI *abi.ABI, eventName string, strict bool) (GetFilterLogsReturnType, error) {
	results := make(GetFilterLogsReturnType, 0, len(logs))

	// Build event lookup map
	var targetEvent *abi.Event
	var allEvents []abi.Event

	if eventName != "" {
		event, err := parsedABI.GetEvent(eventName)
		if err != nil {
			return nil, fmt.Errorf("event %q not found in ABI: %w", eventName, err)
		}
		targetEvent = event
	} else {
		for _, event := range parsedABI.Events {
			allEvents = append(allEvents, event)
		}
	}

	for _, log := range logs {
		if len(log.Topics) == 0 {
			if !strict {
				results = append(results, ContractEventLog{Log: log})
			}
			continue
		}

		// Find matching event
		var matchedEvent *abi.Event
		topicHash := common.HexToHash(log.Topics[0])

		if targetEvent != nil {
			if targetEvent.Topic == topicHash {
				matchedEvent = targetEvent
			}
		} else {
			for i := range allEvents {
				if allEvents[i].Topic == topicHash {
					matchedEvent = &allEvents[i]
					break
				}
			}
		}

		if matchedEvent == nil {
			if strict {
				continue
			}
			results = append(results, ContractEventLog{Log: log})
			continue
		}

		// Decode the event log
		decodedArgs, decodeErr := decodeEventLog(parsedABI, matchedEvent.Name, log)
		if decodeErr != nil {
			if strict {
				continue
			}
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
