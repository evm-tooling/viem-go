package abi

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// RawLog represents a raw Ethereum log with topics and data.
type RawLog struct {
	Address common.Address
	Topics  []common.Hash
	Data    []byte
	// Optional fields
	BlockNumber     uint64
	TransactionHash common.Hash
	LogIndex        uint
}

// ParsedEventLog represents a parsed and decoded event log.
type ParsedEventLog struct {
	EventName string
	Args      map[string]any
	Address   common.Address
	Topics    []common.Hash
	Data      []byte
	// Optional fields from RawLog
	BlockNumber     uint64
	TransactionHash common.Hash
	LogIndex        uint
}

// ParseEventLogsOptions configures how logs are parsed.
type ParseEventLogsOptions struct {
	// EventName filters logs to only the specified event name(s).
	// If nil, all events in the ABI are matched.
	EventName []string
	// Args filters logs to only those matching the specified indexed arguments.
	// The map keys are parameter names, values are the expected values.
	Args map[string]any
	// Strict mode fails on decode errors instead of skipping. Default is true.
	Strict bool
}

// ParseEventLogs extracts and decodes logs matching the ABI from a set of raw logs.
// This is useful for parsing multiple logs at once from a transaction receipt.
//
// Example:
//
//	logs := abi.ParseEventLogs(rawLogs, nil) // parse all matching events
//	logs := abi.ParseEventLogs(rawLogs, &ParseEventLogsOptions{
//	    EventName: []string{"Transfer"},
//	})
func (a *ABI) ParseEventLogs(logs []RawLog, opts *ParseEventLogsOptions) []ParsedEventLog {
	if opts == nil {
		opts = &ParseEventLogsOptions{Strict: true}
	}

	var results []ParsedEventLog

	for _, log := range logs {
		if len(log.Topics) == 0 {
			continue
		}

		// Find matching event by topic signature
		var matchedEvent *Event
		for _, event := range a.Events {
			if event.Topic == log.Topics[0] {
				matchedEvent = &event
				break
			}
		}

		if matchedEvent == nil {
			continue
		}

		// Filter by event name if specified
		if len(opts.EventName) > 0 {
			matched := false
			for _, name := range opts.EventName {
				if matchedEvent.Name == name {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}

		// Decode the event
		decoded, err := a.DecodeEventLogByName(matchedEvent.Name, log.Topics, log.Data)
		if err != nil {
			if opts.Strict {
				continue
			}
			// In non-strict mode, add with empty args
			results = append(results, ParsedEventLog{
				EventName:       matchedEvent.Name,
				Args:            make(map[string]any),
				Address:         log.Address,
				Topics:          log.Topics,
				Data:            log.Data,
				BlockNumber:     log.BlockNumber,
				TransactionHash: log.TransactionHash,
				LogIndex:        log.LogIndex,
			})
			continue
		}

		// Filter by args if specified
		if len(opts.Args) > 0 && !matchesArgs(decoded.Args, opts.Args, matchedEvent) {
			continue
		}

		results = append(results, ParsedEventLog{
			EventName:       decoded.EventName,
			Args:            decoded.Args,
			Address:         log.Address,
			Topics:          log.Topics,
			Data:            log.Data,
			BlockNumber:     log.BlockNumber,
			TransactionHash: log.TransactionHash,
			LogIndex:        log.LogIndex,
		})
	}

	return results
}

// matchesArgs checks if decoded args match the filter args.
func matchesArgs(decodedArgs map[string]any, filterArgs map[string]any, event *Event) bool {
	for name, filterValue := range filterArgs {
		if filterValue == nil {
			continue
		}

		decodedValue, exists := decodedArgs[name]
		if !exists {
			return false
		}

		// Find the parameter type
		var paramType string
		for _, input := range event.Inputs {
			if input.Name == name {
				paramType = input.Type
				break
			}
		}

		if !matchesValue(decodedValue, filterValue, paramType) {
			return false
		}
	}
	return true
}

// matchesValue compares a decoded value with a filter value.
func matchesValue(decoded, filter any, paramType string) bool {
	// Handle address comparison (case-insensitive)
	if paramType == "address" {
		decodedAddr, ok1 := decoded.(common.Address)
		filterAddr, ok2 := filter.(common.Address)
		if ok1 && ok2 {
			return bytes.Equal(decodedAddr.Bytes(), filterAddr.Bytes())
		}
		// Also handle string addresses
		if filterStr, ok := filter.(string); ok && ok1 {
			return bytes.Equal(decodedAddr.Bytes(), common.HexToAddress(filterStr).Bytes())
		}
	}

	// Handle string/bytes comparison (indexed strings/bytes are hashed)
	if paramType == "string" || paramType == "bytes" {
		// If decoded is a hash (indexed), hash the filter value for comparison
		if decodedHash, ok := decoded.(common.Hash); ok {
			if filterStr, ok := filter.(string); ok {
				filterHash := crypto.Keccak256Hash([]byte(filterStr))
				return decodedHash == filterHash
			}
			if filterBytes, ok := filter.([]byte); ok {
				filterHash := crypto.Keccak256Hash(filterBytes)
				return decodedHash == filterHash
			}
		}
	}

	// Direct comparison for other types
	return decoded == filter
}
