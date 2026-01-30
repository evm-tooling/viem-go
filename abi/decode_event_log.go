package abi

import (
	"fmt"
	"math/big"

	gethABI "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// DecodedEventLog represents a decoded event log.
type DecodedEventLog struct {
	EventName string
	Args      map[string]any
	Topics    []common.Hash
	Data      []byte
}

// DecodeEventLog decodes event log data and topics into a structured result.
//
// Example:
//
//	result, err := abi.DecodeEventLog(log.Topics, log.Data)
//	fmt.Println(result.EventName, result.Args)
func (a *ABI) DecodeEventLog(topics []common.Hash, data []byte) (*DecodedEventLog, error) {
	if len(topics) == 0 {
		return nil, fmt.Errorf("topics cannot be empty")
	}

	// First topic is the event signature (for non-anonymous events)
	signature := topics[0]

	// Find the matching event
	var event *gethABI.Event
	var eventName string
	for name, e := range a.gethABI.Events {
		if e.ID == signature {
			event = &e
			eventName = name
			break
		}
	}

	if event == nil {
		return nil, fmt.Errorf("event with signature %s not found on ABI", signature.Hex())
	}

	return a.DecodeEventLogByName(eventName, topics, data)
}

// DecodeEventLogByName decodes event log data using a known event name.
func (a *ABI) DecodeEventLogByName(eventName string, topics []common.Hash, data []byte) (*DecodedEventLog, error) {
	e, ok := a.gethABI.Events[eventName]
	if !ok {
		return nil, fmt.Errorf("event %q not found on ABI", eventName)
	}

	result := make(map[string]any)

	// Separate indexed and non-indexed inputs
	var indexedInputs, nonIndexedInputs []int
	for i, input := range e.Inputs {
		if input.Indexed {
			indexedInputs = append(indexedInputs, i)
		} else {
			nonIndexedInputs = append(nonIndexedInputs, i)
		}
	}

	// Decode indexed topics
	// Skip the first topic if not anonymous (it's the event signature)
	topicOffset := 0
	if !e.Anonymous {
		topicOffset = 1
	}

	for i, idx := range indexedInputs {
		topicIdx := topicOffset + i
		if topicIdx >= len(topics) {
			return nil, fmt.Errorf("not enough topics for event %q: expected %d, got %d", eventName, len(indexedInputs)+topicOffset, len(topics))
		}

		input := e.Inputs[idx]
		// Indexed dynamic types (string, bytes, arrays) are hashed, so we can only return the hash
		typeStr := input.Type.String()
		if typeStr == "string" || typeStr == "bytes" || input.Type.T == gethABI.SliceTy {
			result[input.Name] = topics[topicIdx]
		} else {
			// For fixed-size types, decode the topic as the value
			result[input.Name] = decodeEventTopic(input.Type, topics[topicIdx])
		}
	}

	// Decode non-indexed data
	if len(nonIndexedInputs) > 0 && len(data) > 0 {
		unpacked, err := e.Inputs.UnpackValues(data)
		if err != nil {
			return nil, fmt.Errorf("failed to decode event data for %q: %w", eventName, err)
		}

		// Map unpacked values to non-indexed inputs
		unpackedIdx := 0
		for _, idx := range nonIndexedInputs {
			if unpackedIdx < len(unpacked) {
				result[e.Inputs[idx].Name] = unpacked[unpackedIdx]
				unpackedIdx++
			}
		}
	}

	return &DecodedEventLog{
		EventName: eventName,
		Args:      result,
		Topics:    topics,
		Data:      data,
	}, nil
}

// DecodeEventLogIntoStruct decodes event log data into the provided struct.
func (a *ABI) DecodeEventLogIntoStruct(eventName string, topics []common.Hash, data []byte, output any) error {
	_, ok := a.gethABI.Events[eventName]
	if !ok {
		return fmt.Errorf("event %q not found on ABI", eventName)
	}

	// Use go-ethereum's built-in unpacking
	return a.gethABI.UnpackIntoInterface(output, eventName, data)
}

// decodeEventTopic decodes an indexed topic value based on its type.
func decodeEventTopic(typ gethABI.Type, topic common.Hash) any {
	switch typ.T {
	case gethABI.AddressTy:
		return common.BytesToAddress(topic.Bytes())
	case gethABI.BoolTy:
		return topic[31] != 0
	case gethABI.IntTy:
		// Handle signed integers
		bi := new(big.Int).SetBytes(topic.Bytes())
		// Check if the high bit is set (negative number)
		if topic[0]&0x80 != 0 {
			// Two's complement for negative numbers
			bi.Sub(bi, new(big.Int).Lsh(big.NewInt(1), 256))
		}
		return bi
	case gethABI.UintTy:
		return new(big.Int).SetBytes(topic.Bytes())
	case gethABI.FixedBytesTy:
		// Return the appropriate fixed bytes type
		return topic
	default:
		// For other types, return the hash
		return topic
	}
}

// DecodeEvent is an alias for DecodeEventLogByName that returns a map.
// Deprecated: Use DecodeEventLog or DecodeEventLogByName instead.
func (a *ABI) DecodeEvent(name string, topics []common.Hash, data []byte) (map[string]any, error) {
	result, err := a.DecodeEventLogByName(name, topics, data)
	if err != nil {
		return nil, err
	}
	return result.Args, nil
}

// DecodeEventIntoStruct is an alias for DecodeEventLogIntoStruct.
// Deprecated: Use DecodeEventLogIntoStruct instead.
func (a *ABI) DecodeEventIntoStruct(name string, topics []common.Hash, data []byte, output any) error {
	return a.DecodeEventLogIntoStruct(name, topics, data, output)
}
