package abi

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
)

// Parse parses a JSON ABI and returns an ABI instance.
func Parse(jsonABI []byte) (*ABI, error) {
	gethABI, err := abi.JSON(strings.NewReader(string(jsonABI)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	a := &ABI{
		gethABI:   gethABI,
		raw:       jsonABI,
		Functions: make(map[string]Function),
		Events:    make(map[string]Event),
		Errors:    make(map[string]Error),
	}

	// Convert functions
	for name, method := range gethABI.Methods {
		var selector [4]byte
		copy(selector[:], method.ID)

		a.Functions[name] = Function{
			Name:            method.Name,
			Inputs:          convertGethArgumentsToParameters(method.Inputs),
			Outputs:         convertGethArgumentsToParameters(method.Outputs),
			StateMutability: parseGethStateMutability(method.StateMutability),
			Selector:        selector,
			Signature:       method.Sig,
		}
	}

	// Convert events
	for name, event := range gethABI.Events {
		a.Events[name] = Event{
			Name:      event.Name,
			Inputs:    convertGethArgumentsToParameters(event.Inputs),
			Anonymous: event.Anonymous,
			Topic:     event.ID,
			Signature: event.Sig,
		}
	}

	// Convert errors
	for name, abiErr := range gethABI.Errors {
		var selector [4]byte
		copy(selector[:], abiErr.ID[:4])

		a.Errors[name] = Error{
			Name:      abiErr.Name,
			Inputs:    convertGethArgumentsToParameters(abiErr.Inputs),
			Selector:  selector,
			Signature: abiErr.Sig,
		}
	}

	return a, nil
}

// MustParse parses a JSON ABI and panics on error.
func MustParse(jsonABI []byte) *ABI {
	a, err := Parse(jsonABI)
	if err != nil {
		panic(err)
	}
	return a
}

// ParseFromString parses a JSON ABI string.
func ParseFromString(jsonABI string) (*ABI, error) {
	return Parse([]byte(jsonABI))
}

// GetFunction returns the function with the given name, or an error if not found.
func (a *ABI) GetFunction(name string) (*Function, error) {
	fn, ok := a.Functions[name]
	if !ok {
		return nil, fmt.Errorf("function %q not found in ABI", name)
	}
	return &fn, nil
}

// GetFunctionBySelector returns the function with the given selector.
func (a *ABI) GetFunctionBySelector(selector [4]byte) (*Function, error) {
	for _, fn := range a.Functions {
		if fn.Selector == selector {
			return &fn, nil
		}
	}
	return nil, fmt.Errorf("function with selector %x not found in ABI", selector)
}

// GetEvent returns the event with the given name, or an error if not found.
func (a *ABI) GetEvent(name string) (*Event, error) {
	ev, ok := a.Events[name]
	if !ok {
		return nil, fmt.Errorf("event %q not found in ABI", name)
	}
	return &ev, nil
}

// GetEventByTopic returns the event with the given topic.
func (a *ABI) GetEventByTopic(topic [32]byte) (*Event, error) {
	for _, ev := range a.Events {
		if ev.Topic == topic {
			return &ev, nil
		}
	}
	return nil, fmt.Errorf("event with topic %x not found in ABI", topic)
}

// GetError returns the error with the given name, or an error if not found.
func (a *ABI) GetError(name string) (*Error, error) {
	e, ok := a.Errors[name]
	if !ok {
		return nil, fmt.Errorf("error %q not found in ABI", name)
	}
	return &e, nil
}

// HasFunction returns true if the ABI contains a function with the given name.
func (a *ABI) HasFunction(name string) bool {
	_, ok := a.Functions[name]
	return ok
}

// HasEvent returns true if the ABI contains an event with the given name.
func (a *ABI) HasEvent(name string) bool {
	_, ok := a.Events[name]
	return ok
}

// FunctionNames returns the names of all functions in the ABI.
func (a *ABI) FunctionNames() []string {
	names := make([]string, 0, len(a.Functions))
	for name := range a.Functions {
		names = append(names, name)
	}
	return names
}

// EventNames returns the names of all events in the ABI.
func (a *ABI) EventNames() []string {
	names := make([]string, 0, len(a.Events))
	for name := range a.Events {
		names = append(names, name)
	}
	return names
}

// parseGethStateMutability converts geth's state mutability string to our type.
func parseGethStateMutability(s string) StateMutability {
	return ParseStateMutability(s)
}

// ComputeFunctionSignature computes the function signature from name and input types.
func ComputeFunctionSignature(name string, inputs []Parameter) string {
	var types []string
	for _, input := range inputs {
		types = append(types, input.Type)
	}
	return fmt.Sprintf("%s(%s)", name, strings.Join(types, ","))
}

// ComputeFunctionSelector computes the 4-byte function selector from a signature.
func ComputeFunctionSelector(signature string) [4]byte {
	hash := crypto.Keccak256([]byte(signature))
	var selector [4]byte
	copy(selector[:], hash[:4])
	return selector
}

// MarshalJSON implements json.Marshaler for ABI.
func (a *ABI) MarshalJSON() ([]byte, error) {
	return a.raw, nil
}

// UnmarshalJSON implements json.Unmarshaler for ABI.
func (a *ABI) UnmarshalJSON(data []byte) error {
	parsed, err := Parse(data)
	if err != nil {
		return err
	}
	*a = *parsed
	return nil
}

// ABIItem represents a single ABI item for JSON parsing.
type ABIItem struct {
	Type            string     `json:"type"`
	Name            string     `json:"name,omitempty"`
	Inputs          []ABIInput `json:"inputs,omitempty"`
	Outputs         []ABIInput `json:"outputs,omitempty"`
	StateMutability string     `json:"stateMutability,omitempty"`
	Anonymous       bool       `json:"anonymous,omitempty"`
}

// ABIInput represents an input/output parameter in JSON ABI.
type ABIInput struct {
	Name       string     `json:"name"`
	Type       string     `json:"type"`
	Indexed    bool       `json:"indexed,omitempty"`
	Components []ABIInput `json:"components,omitempty"`
}

// ParseItems parses the ABI JSON into individual items for inspection.
func ParseItems(jsonABI []byte) ([]ABIItem, error) {
	var items []ABIItem
	if err := json.Unmarshal(jsonABI, &items); err != nil {
		return nil, fmt.Errorf("failed to parse ABI items: %w", err)
	}
	return items, nil
}
