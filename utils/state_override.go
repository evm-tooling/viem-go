package utils

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ChefBingbong/viem-go/utils/address"
	"github.com/ChefBingbong/viem-go/utils/encoding"
)

// State override errors
var (
	// ErrInvalidBytesLength is returned when a hex value has an invalid length.
	ErrInvalidBytesLength = errors.New("invalid bytes length")

	// ErrInvalidAddress is returned when an address is invalid.
	ErrInvalidStateAddress = errors.New("invalid address")

	// ErrStateAssignmentConflict is returned when both state and stateDiff are set.
	ErrStateAssignmentConflict = errors.New("state and stateDiff are mutually exclusive")

	// ErrAccountStateConflict is returned when the same address appears multiple times.
	ErrAccountStateConflict = errors.New("state for account already set")
)

// StateMapping represents a storage slot to value mapping.
type StateMapping struct {
	Slot  string `json:"slot"`  // 32-byte hex storage slot
	Value string `json:"value"` // 32-byte hex value
}

// AccountStateOverride represents state overrides for a single account.
type AccountStateOverride struct {
	Address   string          `json:"address"`             // Account address
	Balance   *big.Int        `json:"balance,omitempty"`   // Override balance
	Nonce     *uint64         `json:"nonce,omitempty"`     // Override nonce
	Code      string          `json:"code,omitempty"`      // Override contract code
	State     []StateMapping  `json:"state,omitempty"`     // Replace entire storage
	StateDiff []StateMapping  `json:"stateDiff,omitempty"` // Modify specific storage slots
}

// StateOverride is a list of account state overrides for RPC calls.
type StateOverride []AccountStateOverride

// RpcStateMapping is the RPC format for state mappings (slot -> value map).
type RpcStateMapping map[string]string

// RpcAccountStateOverride is the RPC format for account state overrides.
type RpcAccountStateOverride struct {
	Balance   string          `json:"balance,omitempty"`
	Nonce     string          `json:"nonce,omitempty"`
	Code      string          `json:"code,omitempty"`
	State     RpcStateMapping `json:"state,omitempty"`
	StateDiff RpcStateMapping `json:"stateDiff,omitempty"`
}

// RpcStateOverride is the RPC format for state overrides (address -> override map).
type RpcStateOverride map[string]*RpcAccountStateOverride

// SerializeStateMapping converts a StateMapping slice to RPC format.
// Each slot and value must be exactly 66 characters (0x + 64 hex chars = 32 bytes).
func SerializeStateMapping(stateMapping []StateMapping) (RpcStateMapping, error) {
	if len(stateMapping) == 0 {
		return nil, nil
	}

	result := make(RpcStateMapping)
	for _, sm := range stateMapping {
		if len(sm.Slot) != 66 {
			return nil, fmt.Errorf("%w: slot has size %d, expected 66 (32 bytes)", ErrInvalidBytesLength, len(sm.Slot))
		}
		if len(sm.Value) != 66 {
			return nil, fmt.Errorf("%w: value has size %d, expected 66 (32 bytes)", ErrInvalidBytesLength, len(sm.Value))
		}
		result[sm.Slot] = sm.Value
	}

	return result, nil
}

// SerializeAccountStateOverride converts an AccountStateOverride to RPC format.
func SerializeAccountStateOverride(override AccountStateOverride) (*RpcAccountStateOverride, error) {
	result := &RpcAccountStateOverride{}

	if override.Code != "" {
		result.Code = override.Code
	}

	if override.Balance != nil {
		result.Balance = encoding.NumberToHex(override.Balance)
	}

	if override.Nonce != nil {
		result.Nonce = encoding.NumberToHex(big.NewInt(int64(*override.Nonce)))
	}

	if len(override.State) > 0 {
		stateMapping, err := SerializeStateMapping(override.State)
		if err != nil {
			return nil, err
		}
		result.State = stateMapping
	}

	if len(override.StateDiff) > 0 {
		if result.State != nil {
			return nil, ErrStateAssignmentConflict
		}
		stateDiff, err := SerializeStateMapping(override.StateDiff)
		if err != nil {
			return nil, err
		}
		result.StateDiff = stateDiff
	}

	return result, nil
}

// SerializeStateOverride converts a StateOverride to RPC format.
// Validates addresses and ensures no duplicate addresses.
func SerializeStateOverride(overrides StateOverride) (RpcStateOverride, error) {
	if len(overrides) == 0 {
		return nil, nil
	}

	result := make(RpcStateOverride)
	for _, override := range overrides {
		addr := override.Address

		// Validate address
		if !address.IsAddress(addr) {
			return nil, fmt.Errorf("%w: %s", ErrInvalidStateAddress, addr)
		}

		// Normalize address to lowercase for comparison
		addrLower := strings.ToLower(addr)

		// Check for duplicate
		if _, exists := result[addrLower]; exists {
			return nil, fmt.Errorf("%w: %s", ErrAccountStateConflict, addr)
		}

		// Serialize account state
		accountState, err := SerializeAccountStateOverride(override)
		if err != nil {
			return nil, err
		}

		result[addr] = accountState
	}

	return result, nil
}

// NewStateOverride creates a new StateOverride from individual account overrides.
func NewStateOverride(overrides ...AccountStateOverride) StateOverride {
	return StateOverride(overrides)
}

// NewStateMapping creates a new StateMapping.
func NewStateMapping(slot, value string) StateMapping {
	return StateMapping{
		Slot:  slot,
		Value: value,
	}
}
