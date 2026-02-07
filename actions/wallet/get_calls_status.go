package wallet

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	json "github.com/goccy/go-json"

	"github.com/ChefBingbong/viem-go/utils/data"
	"github.com/ChefBingbong/viem-go/utils/encoding"
)

// FallbackMagicIdentifier is the magic suffix used to identify fallback call batches
// (sent via individual eth_sendTransaction calls instead of wallet_sendCalls).
// This mirrors viem's fallbackMagicIdentifier.
const FallbackMagicIdentifier = "0x5792579257925792579257925792579257925792579257925792579257925792"

// FallbackTransactionErrorMagicIdentifier is the identifier for failed transactions
// in a fallback call batch. This mirrors viem's fallbackTransactionErrorMagicIdentifier.
var FallbackTransactionErrorMagicIdentifier = "0x0000000000000000000000000000000000000000000000000000000000000000"

// GetCallsStatusParameters contains the parameters for the GetCallsStatus action.
// This mirrors viem's GetCallsStatusParameters type.
type GetCallsStatusParameters struct {
	// ID is the identifier of the call batch.
	ID string
}

// CallsStatusReceipt represents a receipt in the calls status response.
type CallsStatusReceipt struct {
	BlockHash        string            `json:"blockHash,omitempty"`
	BlockNumber      *big.Int          `json:"blockNumber,omitempty"`
	GasUsed          *big.Int          `json:"gasUsed,omitempty"`
	Logs             []json.RawMessage `json:"logs,omitempty"`
	Status           string            `json:"status,omitempty"`
	TransactionHash  string            `json:"transactionHash,omitempty"`
	TransactionIndex string            `json:"transactionIndex,omitempty"`
}

// GetCallsStatusReturnType is the return type for the GetCallsStatus action.
// This mirrors viem's GetCallsStatusReturnType type.
type GetCallsStatusReturnType struct {
	Atomic       bool                 `json:"atomic"`
	ChainID      *int64               `json:"chainId,omitempty"`
	Receipts     []CallsStatusReceipt `json:"receipts"`
	StatusCode   int                  `json:"statusCode"`
	Status       string               `json:"status,omitempty"`
	Version      string               `json:"version"`
	Capabilities map[string]any       `json:"capabilities,omitempty"`
}

// rpcCallsStatusResponse is the raw response from wallet_getCallsStatus.
type rpcCallsStatusResponse struct {
	Atomic       bool              `json:"atomic"`
	ChainID      string            `json:"chainId,omitempty"`
	Receipts     []rpcCallsReceipt `json:"receipts,omitempty"`
	Status       interface{}       `json:"status"` // can be int or string
	Version      string            `json:"version"`
	Capabilities map[string]any    `json:"capabilities,omitempty"`
}

// rpcCallsReceipt is the raw receipt from the RPC response.
type rpcCallsReceipt struct {
	BlockHash        string            `json:"blockHash,omitempty"`
	BlockNumber      string            `json:"blockNumber,omitempty"`
	GasUsed          string            `json:"gasUsed,omitempty"`
	Logs             []json.RawMessage `json:"logs,omitempty"`
	Status           string            `json:"status,omitempty"`
	TransactionHash  string            `json:"transactionHash,omitempty"`
	TransactionIndex string            `json:"transactionIndex,omitempty"`
}

// receiptStatuses maps hex status to human-readable status strings.
// This mirrors viem's receiptStatuses.
var receiptStatuses = map[string]string{
	"0x0": "reverted",
	"0x1": "success",
}

// GetCallsStatus returns the status of a call batch that was sent via SendCalls.
//
// This is equivalent to viem's `getCallsStatus` action.
//
// JSON-RPC Method: wallet_getCallsStatus (EIP-5792)
//
// Example:
//
//	status, err := wallet.GetCallsStatus(ctx, client, wallet.GetCallsStatusParameters{
//	    ID: "0xdeadbeef",
//	})
func GetCallsStatus(ctx context.Context, client Client, params GetCallsStatusParameters) (*GetCallsStatusReturnType, error) {
	id := params.ID

	// Check if this is a fallback call batch (sent via individual transactions)
	isFallback := strings.HasSuffix(id, strings.TrimPrefix(FallbackMagicIdentifier, "0x"))
	if isFallback {
		return getCallsStatusFallback(ctx, client, id)
	}

	// Standard wallet_getCallsStatus RPC call
	return getCallsStatusRPC(ctx, client, id)
}

// getCallsStatusFallback handles fallback call batches (sent via eth_sendTransaction).
// This mirrors viem's isTransactions branch in getCallsStatus.
func getCallsStatusFallback(ctx context.Context, client Client, id string) (*GetCallsStatusReturnType, error) {
	// Strip 0x prefix for processing
	hex := strings.TrimPrefix(id, "0x")

	// The last 64 chars are the magic identifier, the 64 before that are the chainId
	if len(hex) < 128 {
		return nil, fmt.Errorf("invalid fallback call batch id: too short")
	}

	// Extract chain ID (64 hex chars = 32 bytes before the magic identifier)
	chainIDHex := "0x" + hex[len(hex)-128:len(hex)-64]
	chainIDHex = data.TrimHex(chainIDHex, data.TrimLeft)

	chainID, err := encoding.HexToNumber(chainIDHex, false)
	if err != nil {
		return nil, fmt.Errorf("failed to parse chain ID: %w", err)
	}

	// Extract transaction hashes (everything before chainId + magic)
	hashesHex := hex[:len(hex)-128]

	// Split into 64-char chunks (32-byte hashes)
	var hashes []string
	for i := 0; i < len(hashesHex); i += 64 {
		end := i + 64
		if end > len(hashesHex) {
			end = len(hashesHex)
		}
		hashes = append(hashes, hashesHex[i:end])
	}

	// Fetch receipts for each hash
	var receipts []CallsStatusReceipt
	allConfirmed := true
	anySuccess := false
	anyReverted := false

	for _, hash := range hashes {
		// Skip error placeholder hashes
		if hash == strings.TrimPrefix(FallbackTransactionErrorMagicIdentifier, "0x") {
			continue
		}

		resp, reqErr := client.Request(ctx, "eth_getTransactionReceipt", "0x"+hash)
		if reqErr != nil {
			return nil, fmt.Errorf("eth_getTransactionReceipt failed: %w", reqErr)
		}

		if resp.Result == nil || string(resp.Result) == "null" {
			allConfirmed = false
			continue
		}

		var rpcReceipt rpcCallsReceipt
		if unmarshalErr := json.Unmarshal(resp.Result, &rpcReceipt); unmarshalErr != nil {
			return nil, fmt.Errorf("failed to unmarshal receipt: %w", unmarshalErr)
		}

		receipt := formatCallsReceipt(rpcReceipt)
		receipts = append(receipts, receipt)

		switch rpcReceipt.Status {
		case "0x1":
			anySuccess = true
		case "0x0":
			anyReverted = true
		}
	}

	// Determine status code (mirrors viem's status logic)
	var statusCode int
	if !allConfirmed {
		statusCode = 100 // pending
	} else if anySuccess && !anyReverted {
		statusCode = 200 // success
	} else if anyReverted && !anySuccess {
		statusCode = 500 // complete failure
	} else {
		statusCode = 600 // partial failure
	}

	status, _ := resolveStatusFromCode(statusCode)

	return &GetCallsStatusReturnType{
		Atomic:     false,
		ChainID:    &chainID,
		Receipts:   receipts,
		StatusCode: statusCode,
		Status:     status,
		Version:    "2.0.0",
	}, nil
}

// getCallsStatusRPC handles standard wallet_getCallsStatus RPC calls.
func getCallsStatusRPC(ctx context.Context, client Client, id string) (*GetCallsStatusReturnType, error) {
	resp, err := client.Request(ctx, "wallet_getCallsStatus", id)
	if err != nil {
		return nil, fmt.Errorf("wallet_getCallsStatus failed: %w", err)
	}

	var raw rpcCallsStatusResponse
	if unmarshalErr := json.Unmarshal(resp.Result, &raw); unmarshalErr != nil {
		return nil, fmt.Errorf("failed to unmarshal calls status: %w", unmarshalErr)
	}

	// Parse status code from the response (can be int or string for backwards compat)
	statusCode := parseStatusCode(raw.Status)
	status, resolvedCode := resolveStatusFromCode(statusCode)

	// Parse chain ID
	var chainID *int64
	if raw.ChainID != "" {
		cid, cidErr := encoding.HexToNumber(raw.ChainID, false)
		if cidErr == nil {
			chainID = &cid
		}
	}

	// Format receipts
	receipts := make([]CallsStatusReceipt, len(raw.Receipts))
	for i, r := range raw.Receipts {
		receipts[i] = formatCallsReceipt(r)
	}

	return &GetCallsStatusReturnType{
		Atomic:       raw.Atomic,
		ChainID:      chainID,
		Receipts:     receipts,
		StatusCode:   resolvedCode,
		Status:       status,
		Version:      ifEmpty(raw.Version, "2.0.0"),
		Capabilities: raw.Capabilities,
	}, nil
}

// formatCallsReceipt converts an RPC receipt to a CallsStatusReceipt with parsed fields.
func formatCallsReceipt(r rpcCallsReceipt) CallsStatusReceipt {
	receipt := CallsStatusReceipt{
		BlockHash:        r.BlockHash,
		Logs:             r.Logs,
		TransactionHash:  r.TransactionHash,
		TransactionIndex: r.TransactionIndex,
	}

	// Parse blockNumber from hex
	if r.BlockNumber != "" {
		bn, err := encoding.HexToBigInt(r.BlockNumber, false)
		if err == nil {
			receipt.BlockNumber = bn
		}
	}

	// Parse gasUsed from hex
	if r.GasUsed != "" {
		gu, err := encoding.HexToBigInt(r.GasUsed, false)
		if err == nil {
			receipt.GasUsed = gu
		}
	}

	// Map status from hex to human-readable (mirrors viem's receiptStatuses)
	if mapped, ok := receiptStatuses[r.Status]; ok {
		receipt.Status = mapped
	} else {
		receipt.Status = r.Status
	}

	return receipt
}

// resolveStatusFromCode maps a status code to a status string.
// This mirrors viem's status resolution logic.
func resolveStatusFromCode(code int) (string, int) {
	switch {
	case code >= 100 && code < 200:
		return "pending", code
	case code >= 200 && code < 300:
		return "success", code
	case code >= 300 && code < 700:
		return "failure", code
	default:
		return "", code
	}
}

// parseStatusCode extracts an integer status code from a response value
// that may be an int or a string (for backwards compatibility).
func parseStatusCode(v interface{}) int {
	switch s := v.(type) {
	case float64:
		return int(s)
	case int:
		return s
	case string:
		// Backwards compatibility with older string statuses
		switch s {
		case "CONFIRMED":
			return 200
		case "PENDING":
			return 100
		default:
			return 0
		}
	default:
		return 0
	}
}

// ifEmpty returns fallback if s is empty.
func ifEmpty(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}
