package abi

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// EncodeEventTopics encodes indexed event parameters as topics.
// Returns an array of topic hashes, where the first topic is the event signature
// (unless the event is anonymous).
//
// Example:
//
//	topics, err := abi.EncodeEventTopics("Transfer", from, to)
func (a *ABI) EncodeEventTopics(eventName string, indexedArgs ...any) ([][]byte, error) {
	e, ok := a.gethABI.Events[eventName]
	if !ok {
		return nil, fmt.Errorf("event %q not found on ABI", eventName)
	}

	// First topic is always the event signature (unless anonymous)
	var topics [][]byte
	if !e.Anonymous {
		topics = append(topics, e.ID.Bytes())
	}

	// Count indexed parameters
	indexedCount := 0
	for _, input := range e.Inputs {
		if input.Indexed {
			indexedCount++
		}
	}

	if len(indexedArgs) != indexedCount {
		return nil, fmt.Errorf("event %q expects %d indexed arguments, got %d", eventName, indexedCount, len(indexedArgs))
	}

	// Encode indexed parameters
	argIndex := 0
	for _, input := range e.Inputs {
		if input.Indexed {
			if indexedArgs[argIndex] != nil {
				// Encode the indexed argument as a topic
				topic, err := encodeEventTopic(input.Type.String(), indexedArgs[argIndex])
				if err != nil {
					return nil, fmt.Errorf("failed to encode indexed argument %q: %w", input.Name, err)
				}
				topics = append(topics, topic)
			} else {
				// nil means "match any" - add nil topic
				topics = append(topics, nil)
			}
			argIndex++
		}
	}

	return topics, nil
}

// encodeEventTopic encodes a value as an indexed topic (32 bytes).
func encodeEventTopic(typeStr string, value any) ([]byte, error) {
	topic := make([]byte, 32)

	switch v := value.(type) {
	case common.Address:
		copy(topic[12:], v.Bytes())
	case *common.Address:
		if v != nil {
			copy(topic[12:], v.Bytes())
		}
	case string:
		// For string types, check if it looks like an address
		if common.IsHexAddress(v) && (typeStr == "address" || typeStr == "") {
			addr := common.HexToAddress(v)
			copy(topic[12:], addr.Bytes())
		} else {
			// For string parameters, use keccak256 hash
			hash := crypto.Keccak256([]byte(v))
			copy(topic, hash)
		}
	case *big.Int:
		if v != nil {
			// Handle negative numbers (two's complement)
			if v.Sign() < 0 {
				// For negative numbers, we need to handle two's complement
				// Create a 256-bit representation
				twosComplement := new(big.Int).Add(new(big.Int).Lsh(big.NewInt(1), 256), v)
				b := twosComplement.Bytes()
				copy(topic[32-len(b):], b)
			} else {
				b := v.Bytes()
				copy(topic[32-len(b):], b)
			}
		}
	case int:
		bi := big.NewInt(int64(v))
		return encodeEventTopic(typeStr, bi)
	case int8:
		bi := big.NewInt(int64(v))
		return encodeEventTopic(typeStr, bi)
	case int16:
		bi := big.NewInt(int64(v))
		return encodeEventTopic(typeStr, bi)
	case int32:
		bi := big.NewInt(int64(v))
		return encodeEventTopic(typeStr, bi)
	case int64:
		bi := big.NewInt(v)
		return encodeEventTopic(typeStr, bi)
	case uint:
		bi := new(big.Int).SetUint64(uint64(v))
		return encodeEventTopic(typeStr, bi)
	case uint8:
		bi := new(big.Int).SetUint64(uint64(v))
		return encodeEventTopic(typeStr, bi)
	case uint16:
		bi := new(big.Int).SetUint64(uint64(v))
		return encodeEventTopic(typeStr, bi)
	case uint32:
		bi := new(big.Int).SetUint64(uint64(v))
		return encodeEventTopic(typeStr, bi)
	case uint64:
		bi := new(big.Int).SetUint64(v)
		return encodeEventTopic(typeStr, bi)
	case bool:
		if v {
			topic[31] = 1
		}
	case []byte:
		// For bytes/string, use keccak256 hash
		hash := crypto.Keccak256(v)
		copy(topic, hash)
	case common.Hash:
		copy(topic, v.Bytes())
	case [32]byte:
		copy(topic, v[:])
	default:
		return nil, fmt.Errorf("unsupported type for indexed argument: %T", value)
	}

	return topic, nil
}
