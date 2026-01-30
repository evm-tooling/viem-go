package abi

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

// AbiItem represents a generic ABI item (function, event, or error).
type AbiItem interface{}

// GetAbiItemOptions configures how to find an ABI item.
type GetAbiItemOptions struct {
	// Args are used to disambiguate overloaded functions/events.
	Args []any
}

// GetAbiItem finds an ABI item by name or selector.
// For overloaded functions, pass Args to disambiguate.
//
// Example:
//
//	item, err := abi.GetAbiItem("transfer", nil) // by name
//	item, err := abi.GetAbiItem("0xa9059cbb", nil) // by selector
func (a *ABI) GetAbiItem(nameOrSelector string, opts *GetAbiItemOptions) (AbiItem, error) {
	if opts == nil {
		opts = &GetAbiItemOptions{}
	}

	isSelector := isHexSelector(nameOrSelector)

	// Collect all matching items
	var matches []AbiItem

	if isSelector {
		selector := normalizeSelector(nameOrSelector)

		// Check functions
		for _, fn := range a.Functions {
			if matchesSelector(fn.Selector[:], selector) {
				matches = append(matches, fn)
			}
		}

		// Check events (32-byte topic)
		if len(selector) == 32 {
			for _, ev := range a.Events {
				if ev.Topic.Hex() == "0x"+hex.EncodeToString(selector) {
					matches = append(matches, ev)
				}
			}
		}

		// Check errors
		for _, err := range a.Errors {
			if matchesSelector(err.Selector[:], selector) {
				matches = append(matches, err)
			}
		}
	} else {
		// Match by name
		for _, fn := range a.Functions {
			if fn.Name == nameOrSelector {
				matches = append(matches, fn)
			}
		}
		for _, ev := range a.Events {
			if ev.Name == nameOrSelector {
				matches = append(matches, ev)
			}
		}
		for _, err := range a.Errors {
			if err.Name == nameOrSelector {
				matches = append(matches, err)
			}
		}
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("ABI item %q not found", nameOrSelector)
	}

	if len(matches) == 1 {
		return matches[0], nil
	}

	// Multiple matches - need to disambiguate with args
	if len(opts.Args) == 0 {
		// No args provided, check if any item has no inputs
		for _, item := range matches {
			inputs := getItemInputs(item)
			if len(inputs) == 0 {
				return item, nil
			}
		}
		// Return first match (like viem does)
		return matches[0], nil
	}

	// Try to match based on args
	var matchedItem AbiItem
	for _, item := range matches {
		inputs := getItemInputs(item)
		if len(inputs) != len(opts.Args) {
			continue
		}

		// Check if all args match the expected types
		allMatch := true
		for i, arg := range opts.Args {
			if !isArgOfType(arg, inputs[i]) {
				allMatch = false
				break
			}
		}

		if allMatch {
			if matchedItem != nil {
				// Ambiguous - check for type ambiguity
				prevInputs := getItemInputs(matchedItem)
				ambiguous := checkAmbiguousTypes(inputs, prevInputs, opts.Args)
				if ambiguous != "" {
					return nil, fmt.Errorf("ambiguous ABI item: %s", ambiguous)
				}
			}
			matchedItem = item
		}
	}

	if matchedItem != nil {
		return matchedItem, nil
	}

	// Fall back to first match
	return matches[0], nil
}

// GetFunction finds a function by name or selector.
func (a *ABI) GetFunction(name string) (*Function, error) {
	fn, ok := a.Functions[name]
	if !ok {
		return nil, fmt.Errorf("function %q not found on ABI", name)
	}
	return &fn, nil
}

// GetEvent finds an event by name.
func (a *ABI) GetEvent(name string) (*Event, error) {
	ev, ok := a.Events[name]
	if !ok {
		return nil, fmt.Errorf("event %q not found on ABI", name)
	}
	return &ev, nil
}

// GetError finds an error by name.
func (a *ABI) GetError(name string) (*Error, error) {
	err, ok := a.Errors[name]
	if !ok {
		return nil, fmt.Errorf("error %q not found on ABI", name)
	}
	return &err, nil
}

// GetFunctionBySelector finds a function by its 4-byte selector.
func (a *ABI) GetFunctionBySelector(selector [4]byte) (*Function, error) {
	for _, fn := range a.Functions {
		if fn.Selector == selector {
			return &fn, nil
		}
	}
	return nil, fmt.Errorf("function with selector 0x%x not found", selector)
}

// GetEventByTopic finds an event by its topic hash.
func (a *ABI) GetEventByTopic(topic common.Hash) (*Event, error) {
	for _, ev := range a.Events {
		if ev.Topic == topic {
			return &ev, nil
		}
	}
	return nil, fmt.Errorf("event with topic %s not found", topic.Hex())
}

// GetErrorBySelector finds an error by its 4-byte selector.
func (a *ABI) GetErrorBySelector(selector [4]byte) (*Error, error) {
	for _, err := range a.Errors {
		if err.Selector == selector {
			return &err, nil
		}
	}
	return nil, fmt.Errorf("error with selector 0x%x not found", selector)
}

// Helper functions

func isHexSelector(s string) bool {
	if !strings.HasPrefix(s, "0x") && !strings.HasPrefix(s, "0X") {
		return false
	}
	s = s[2:]
	// Must be 8 chars (4 bytes) for function/error or 64 chars (32 bytes) for event
	if len(s) != 8 && len(s) != 64 {
		return false
	}
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

func normalizeSelector(s string) []byte {
	s = strings.TrimPrefix(s, "0x")
	s = strings.TrimPrefix(s, "0X")
	b, _ := hex.DecodeString(s)
	return b
}

func matchesSelector(itemSelector []byte, querySelector []byte) bool {
	if len(querySelector) > len(itemSelector) {
		return false
	}
	for i, b := range querySelector {
		if itemSelector[i] != b {
			return false
		}
	}
	return true
}

func getItemInputs(item AbiItem) []Parameter {
	switch v := item.(type) {
	case Function:
		return v.Inputs
	case *Function:
		return v.Inputs
	case Event:
		return v.Inputs
	case *Event:
		return v.Inputs
	case Error:
		return v.Inputs
	case *Error:
		return v.Inputs
	default:
		return nil
	}
}

// isArgOfType checks if an argument matches the expected ABI parameter type.
func isArgOfType(arg any, param Parameter) bool {
	if arg == nil {
		return true
	}

	argType := reflect.TypeOf(arg)
	paramType := param.Type

	switch paramType {
	case "address":
		// Accept common.Address or string that looks like an address
		if _, ok := arg.(common.Address); ok {
			return true
		}
		if s, ok := arg.(string); ok {
			return common.IsHexAddress(s)
		}
		return false

	case "bool":
		_, ok := arg.(bool)
		return ok

	case "string":
		_, ok := arg.(string)
		return ok

	case "bytes":
		_, ok := arg.([]byte)
		if ok {
			return true
		}
		_, ok = arg.(string)
		return ok

	default:
		// Handle integers
		if intRegex.MatchString(paramType) {
			return isIntegerType(arg)
		}

		// Handle fixed bytes (bytes1 to bytes32)
		if fixedBytesRegex.MatchString(paramType) {
			return isBytesType(arg)
		}

		// Handle arrays
		if arrayTypeRegex.MatchString(paramType) {
			return reflect.TypeOf(arg).Kind() == reflect.Slice || reflect.TypeOf(arg).Kind() == reflect.Array
		}

		// Handle tuple
		if paramType == "tuple" {
			return argType.Kind() == reflect.Struct || argType.Kind() == reflect.Map
		}

		return false
	}
}

var (
	intRegex        = regexp.MustCompile(`^u?int(8|16|24|32|40|48|56|64|72|80|88|96|104|112|120|128|136|144|152|160|168|176|184|192|200|208|216|224|232|240|248|256)?$`)
	fixedBytesRegex = regexp.MustCompile(`^bytes([1-9]|1[0-9]|2[0-9]|3[0-2])$`)
	arrayTypeRegex  = regexp.MustCompile(`^.+\[\d*\]$`)
)

func isIntegerType(arg any) bool {
	switch arg.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		*big.Int:
		return true
	default:
		return false
	}
}

func isBytesType(arg any) bool {
	switch arg.(type) {
	case []byte, string,
		[1]byte, [2]byte, [3]byte, [4]byte, [5]byte, [6]byte, [7]byte, [8]byte,
		[9]byte, [10]byte, [11]byte, [12]byte, [13]byte, [14]byte, [15]byte, [16]byte,
		[17]byte, [18]byte, [19]byte, [20]byte, [21]byte, [22]byte, [23]byte, [24]byte,
		[25]byte, [26]byte, [27]byte, [28]byte, [29]byte, [30]byte, [31]byte, [32]byte:
		return true
	default:
		return false
	}
}

func checkAmbiguousTypes(inputs1, inputs2 []Parameter, args []any) string {
	for i := range inputs1 {
		type1 := inputs1[i].Type
		type2 := inputs2[i].Type

		// Check for address vs bytes20 ambiguity
		if (type1 == "address" && type2 == "bytes20") || (type1 == "bytes20" && type2 == "address") {
			return fmt.Sprintf("ambiguous types: %s vs %s", type1, type2)
		}

		// Check for address vs string ambiguity
		if (type1 == "address" && type2 == "string") || (type1 == "string" && type2 == "address") {
			if s, ok := args[i].(string); ok && common.IsHexAddress(s) {
				return fmt.Sprintf("ambiguous types: %s vs %s", type1, type2)
			}
		}
	}
	return ""
}
