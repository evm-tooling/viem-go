package wallet

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ChefBingbong/viem-go/utils/encoding"
)

// GetCapabilitiesParameters contains the parameters for the GetCapabilities action.
// This mirrors viem's GetCapabilitiesParameters type.
type GetCapabilitiesParameters struct {
	// Account is the account address to get capabilities for.
	// If nil, uses the client's account.
	Account *string

	// ChainID optionally restricts capabilities to a specific chain.
	// If set, returns capabilities only for this chain ID.
	// If nil, returns capabilities for all chains.
	ChainID *int64
}

// Capabilities represents a set of wallet capabilities as a flexible map.
// This mirrors viem's Capabilities type from EIP-5792.
type Capabilities = map[string]any

// ChainIDToCapabilities maps chain IDs to their capabilities.
// This mirrors viem's ChainIdToCapabilities type.
type ChainIDToCapabilities = map[int64]Capabilities

// GetCapabilitiesReturnType is the return type for the GetCapabilities action.
// It can be either a ChainIDToCapabilities (all chains) or Capabilities (single chain).
// The caller should use the appropriate type based on whether ChainID was provided.
type GetCapabilitiesReturnType = ChainIDToCapabilities

// GetCapabilities extracts capabilities that a connected wallet supports
// (e.g. paymasters, session keys, etc).
//
// This is equivalent to viem's `getCapabilities` action.
//
// JSON-RPC Method: wallet_getCapabilities (EIP-5792)
//
// Example (all chains):
//
//	capabilities, err := wallet.GetCapabilities(ctx, client, wallet.GetCapabilitiesParameters{})
//	// capabilities[1]["paymasterService"] => map[string]any{...}
//
// Example (specific chain):
//
//	chainID := int64(1)
//	capabilities, err := wallet.GetCapabilities(ctx, client, wallet.GetCapabilitiesParameters{
//	    ChainID: &chainID,
//	})
//	// capabilities[1] contains only mainnet capabilities
func GetCapabilities(ctx context.Context, client Client, params GetCapabilitiesParameters) (GetCapabilitiesReturnType, error) {
	// Determine account address
	var accountAddr *string
	if params.Account != nil {
		accountAddr = params.Account
	} else if acct := client.Account(); acct != nil {
		addr := acct.Address().Hex()
		accountAddr = &addr
	}

	// Build RPC params - mirrors viem's conditional params construction
	var rpcParams []any
	if params.ChainID != nil {
		chainIDHex := encoding.NumberToHex(big.NewInt(*params.ChainID))
		if accountAddr != nil {
			rpcParams = []any{*accountAddr, []string{chainIDHex}}
		} else {
			rpcParams = []any{nil, []string{chainIDHex}}
		}
	} else {
		if accountAddr != nil {
			rpcParams = []any{*accountAddr}
		}
	}

	resp, err := client.Request(ctx, "wallet_getCapabilities", rpcParams...)
	if err != nil {
		return nil, fmt.Errorf("wallet_getCapabilities failed: %w", err)
	}

	// Parse raw capabilities response - keys are hex chain IDs, values are capability objects
	var rawCapabilities map[string]map[string]any
	if unmarshalErr := json.Unmarshal(resp.Result, &rawCapabilities); unmarshalErr != nil {
		return nil, fmt.Errorf("failed to unmarshal capabilities: %w", unmarshalErr)
	}

	// Convert hex chain ID keys to int64 and normalize capability keys
	// Mirrors viem's post-processing: renaming "addSubAccount" to "unstable_addSubAccount"
	capabilities := make(ChainIDToCapabilities)
	for chainIDStr, caps := range rawCapabilities {
		chainID, parseErr := parseChainID(chainIDStr)
		if parseErr != nil {
			return nil, fmt.Errorf("failed to parse chain id %q: %w", chainIDStr, parseErr)
		}

		normalizedCaps := make(Capabilities)
		for key, value := range caps {
			// Normalize capability keys (mirrors viem's key renaming)
			if key == "addSubAccount" {
				key = "unstable_addSubAccount"
			}
			normalizedCaps[key] = value
		}
		capabilities[chainID] = normalizedCaps
	}

	return capabilities, nil
}

// parseChainID parses a chain ID from a hex string (e.g., "0x1") or decimal string.
func parseChainID(s string) (int64, error) {
	if len(s) >= 2 && (s[:2] == "0x" || s[:2] == "0X") {
		// Hex format
		n := new(big.Int)
		_, ok := n.SetString(s[2:], 16)
		if !ok {
			return 0, fmt.Errorf("invalid hex chain id: %s", s)
		}
		return n.Int64(), nil
	}
	// Decimal format
	return strconv.ParseInt(s, 10, 64)
}
