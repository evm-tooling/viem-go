package wallet

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	viemabi "github.com/ChefBingbong/viem-go/abi"
	"github.com/ChefBingbong/viem-go/utils/data"
	"github.com/ChefBingbong/viem-go/utils/encoding"
)

// Call represents a single call in a batch.
// This mirrors viem's Call type.
type Call struct {
	// To is the target address.
	To string `json:"to,omitempty"`
	// Data is the calldata (hex string).
	Data string `json:"data,omitempty"`
	// Value is the amount of ETH to send.
	Value *big.Int `json:"value,omitempty"`

	// ABI is the contract ABI for encoding function calls (optional).
	// If provided, FunctionName and Args are used to encode the data.
	ABI any `json:"-"` // []byte, string, or *abi.ABI
	// FunctionName is the function to call (used with ABI).
	FunctionName string `json:"-"`
	// Args are the function arguments (used with ABI).
	Args []any `json:"-"`

	// DataSuffix is data to append to the end of the calldata.
	DataSuffix string `json:"-"`
}

// SendCallsParameters contains the parameters for the SendCalls action.
// This mirrors viem's SendCallsParameters type.
type SendCallsParameters struct {
	// Account is the account to send from. If nil, uses the client's account.
	Account Account

	// Calls is the list of calls to send.
	Calls []Call

	// Capabilities is the optional capabilities to request.
	Capabilities map[string]any `json:"capabilities,omitempty"`

	// ForceAtomic when true, requires the batch to execute atomically. Default: false.
	ForceAtomic bool

	// ID is an optional identifier for the call batch.
	ID string

	// Version is the EIP-5792 version. Default: "2.0.0".
	Version string

	// ExperimentalFallback when true, falls back to eth_sendTransaction if
	// wallet_sendCalls is not supported.
	ExperimentalFallback bool

	// ExperimentalFallbackDelay is the delay (in ms) between fallback transactions.
	// Default: 32ms.
	ExperimentalFallbackDelay *int
}

// SendCallsReturnType is the return type for the SendCalls action.
// This mirrors viem's SendCallsReturnType type.
type SendCallsReturnType struct {
	ID           string         `json:"id"`
	Capabilities map[string]any `json:"capabilities,omitempty"`
}

// sendCallsRpcCall is the formatted call for the RPC request.
type sendCallsRpcCall struct {
	Data  string `json:"data,omitempty"`
	To    string `json:"to,omitempty"`
	Value string `json:"value,omitempty"`
}

// sendCallsRpcParams is the RPC params for wallet_sendCalls.
type sendCallsRpcParams struct {
	AtomicRequired bool               `json:"atomicRequired"`
	Calls          []sendCallsRpcCall `json:"calls"`
	Capabilities   map[string]any     `json:"capabilities,omitempty"`
	ChainID        string             `json:"chainId"`
	From           string             `json:"from,omitempty"`
	ID             string             `json:"id,omitempty"`
	Version        string             `json:"version"`
}

// SendCalls requests the connected wallet to send a batch of calls.
//
// This is equivalent to viem's `sendCalls` action.
//
// JSON-RPC Method: wallet_sendCalls (EIP-5792)
//
// Example:
//
//	result, err := wallet.SendCalls(ctx, client, wallet.SendCallsParameters{
//	    Calls: []wallet.Call{
//	        {Data: "0xdeadbeef", To: "0x70997970c51812dc3a010c7d01b50e0d17dc79c8"},
//	        {To: "0x70997970c51812dc3a010c7d01b50e0d17dc79c8", Value: big.NewInt(69420)},
//	    },
//	})
func SendCalls(ctx context.Context, client Client, params SendCallsParameters) (*SendCallsReturnType, error) {
	// Resolve account
	account := params.Account
	if account == nil {
		account = client.Account()
	}

	// Resolve version
	version := params.Version
	if version == "" {
		version = "2.0.0"
	}

	// Resolve chain
	ch := client.Chain()
	if ch == nil {
		return nil, fmt.Errorf("chain is required for wallet_sendCalls")
	}

	// Propagate client.DataSuffix() to capabilities if not already set.
	// Mirrors viem's: if (client.dataSuffix && !parameters.capabilities?.dataSuffix) { ... }
	capabilities := params.Capabilities
	if clientSuffix := client.DataSuffix(); len(clientSuffix) > 0 {
		if capabilities == nil {
			capabilities = make(map[string]any)
		}
		if _, hasDataSuffix := capabilities["dataSuffix"]; !hasDataSuffix {
			capabilities["dataSuffix"] = map[string]any{
				"value":    encoding.BytesToHex(clientSuffix),
				"optional": true,
			}
		}
	}

	// Encode calls (mirrors viem's calls.map encoding)
	rpcCalls := make([]sendCallsRpcCall, len(params.Calls))
	for i, call := range params.Calls {
		callData := call.Data

		// If ABI is provided, encode the function data
		if call.ABI != nil {
			parsedABI, parseErr := parseABIForCalls(call.ABI)
			if parseErr != nil {
				return nil, fmt.Errorf("failed to parse ABI for call %d: %w", i, parseErr)
			}
			encoded, encErr := parsedABI.EncodeFunctionData(call.FunctionName, call.Args...)
			if encErr != nil {
				return nil, fmt.Errorf("failed to encode function data for call %d: %w", i, encErr)
			}
			callData = "0x" + fmt.Sprintf("%x", encoded)
		}

		// Apply data suffix if present
		if call.DataSuffix != "" && callData != "" {
			callData = data.ConcatHex(callData, call.DataSuffix)
		}

		rpcCall := sendCallsRpcCall{
			Data: callData,
			To:   call.To,
		}
		if call.Value != nil {
			rpcCall.Value = encoding.NumberToHex(call.Value)
		}
		rpcCalls[i] = rpcCall
	}

	// Build RPC params
	rpcParams := sendCallsRpcParams{
		AtomicRequired: params.ForceAtomic,
		Calls:          rpcCalls,
		Capabilities:   capabilities,
		ChainID:        encoding.NumberToHex(big.NewInt(ch.ID)),
		ID:             params.ID,
		Version:        version,
	}
	if account != nil {
		rpcParams.From = account.Address().Hex()
	}

	// Send wallet_sendCalls request
	resp, err := client.Request(ctx, "wallet_sendCalls", rpcParams)
	if err != nil {
		// Handle fallback to eth_sendTransaction
		if params.ExperimentalFallback && isMethodNotSupportedError(err) {
			return sendCallsFallback(ctx, client, account, params, rpcCalls)
		}
		return nil, fmt.Errorf("wallet_sendCalls failed: %w", err)
	}

	// Response can be a string (just id) or an object
	var result SendCallsReturnType
	// Try to unmarshal as object first
	if unmarshalErr := json.Unmarshal(resp.Result, &result); unmarshalErr != nil {
		// Try as string (just the id)
		var idStr string
		if strErr := json.Unmarshal(resp.Result, &idStr); strErr != nil {
			return nil, fmt.Errorf("failed to unmarshal sendCalls response: %w", unmarshalErr)
		}
		result.ID = idStr
	}

	return &result, nil
}

// sendCallsFallback falls back to individual eth_sendTransaction calls.
// This mirrors viem's experimental_fallback branch.
func sendCallsFallback(
	ctx context.Context,
	client Client,
	account Account,
	params SendCallsParameters,
	rpcCalls []sendCallsRpcCall,
) (*SendCallsReturnType, error) {
	// Check for non-optional capabilities
	if params.Capabilities != nil {
		for _, cap := range params.Capabilities {
			if capMap, ok := cap.(map[string]any); ok {
				if optional, exists := capMap["optional"]; exists {
					if optBool, ok := optional.(bool); ok && !optBool {
						return nil, fmt.Errorf("non-optional capabilities are not supported on fallback to eth_sendTransaction")
					}
				} else {
					// No "optional" key means it's required
					return nil, fmt.Errorf("non-optional capabilities are not supported on fallback to eth_sendTransaction")
				}
			}
		}
	}

	// Check atomicity constraint
	if params.ForceAtomic && len(rpcCalls) > 1 {
		return nil, fmt.Errorf("forceAtomic is not supported on fallback to eth_sendTransaction")
	}

	// Resolve fallback delay
	fallbackDelay := 32
	if params.ExperimentalFallbackDelay != nil {
		fallbackDelay = *params.ExperimentalFallbackDelay
	}

	// Get chain ID
	var chainID int64
	if clientChain := client.Chain(); clientChain != nil {
		chainID = clientChain.ID
	}

	// Send each call individually
	hashes := make([]string, len(rpcCalls))
	allFailed := true

	for i, call := range rpcCalls {
		var value *big.Int
		if call.Value != "" {
			v, vErr := encoding.HexToBigInt(call.Value, false)
			if vErr == nil {
				value = v
			}
		}

		hash, txErr := SendTransaction(ctx, client, SendTransactionParameters{
			Account: account,
			Data:    call.Data,
			To:      call.To,
			Value:   value,
		})

		if txErr != nil {
			hashes[i] = FallbackTransactionErrorMagicIdentifier
		} else {
			hashes[i] = hash
			allFailed = false
		}

		// Delay between transactions (mirrors viem's experimental_fallbackDelay)
		if fallbackDelay > 0 && i < len(rpcCalls)-1 {
			time.Sleep(time.Duration(fallbackDelay) * time.Millisecond)
		}
	}

	if allFailed {
		return nil, fmt.Errorf("all fallback transactions failed")
	}

	// Build composite ID: concat(hashes..., chainId(32 bytes), magicIdentifier)
	chainIDHex, _ := encoding.NumberToHexWithSize(big.NewInt(chainID), 32, false)
	compositeID := data.ConcatHex(append(hashes, chainIDHex, FallbackMagicIdentifier)...)

	return &SendCallsReturnType{ID: compositeID}, nil
}

// isMethodNotSupportedError checks if an error indicates the RPC method is not supported.
func isMethodNotSupportedError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	lower := toLower(errStr)
	return contains(lower, "method not found") ||
		contains(lower, "method not supported") ||
		contains(lower, "does not exist") ||
		contains(lower, "is not available") ||
		contains(lower, "missing or invalid") ||
		contains(lower, "did not match any variant") ||
		contains(lower, "account upgraded to unsupported contract") ||
		contains(lower, "eip-7702 not supported") ||
		contains(lower, "unsupported wc_ method") ||
		contains(lower, "feature toggled misconfigured") ||
		contains(lower, "jsonrpcengine: response has no error or result")
}

// parseABIForCalls parses ABI from various formats.
func parseABIForCalls(abiParam any) (*viemabi.ABI, error) {
	switch v := abiParam.(type) {
	case *viemabi.ABI:
		return v, nil
	case []byte:
		return viemabi.Parse(v)
	case string:
		return viemabi.Parse([]byte(v))
	default:
		return nil, fmt.Errorf("ABI must be []byte, string, or *abi.ABI, got %T", abiParam)
	}
}

func toLower(s string) string {
	b := make([]byte, len(s))
	for i := range s {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		} else {
			b[i] = c
		}
	}
	return string(b)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
