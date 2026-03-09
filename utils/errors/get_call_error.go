// Package errors provides error handling utilities.
package errors

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/utils/unit"
)

// CallErrorParams contains parameters for wrapping call errors.
// Mirrors viem's CallParameters passed to CallExecutionError.
type CallErrorParams struct {
	From                 *common.Address
	To                   *common.Address
	Data                 []byte
	Value                *big.Int
	Gas                  *uint64
	GasPrice             *big.Int
	MaxFeePerGas         *big.Int
	MaxPriorityFeePerGas *big.Int
	ChainID              *int64
	NativeCurrencySymbol string // e.g. "POL", "ETH" for value formatting
	FormattedData        string // Optional: e.g. "name()" or "transfer(0x..., 1)"
	StateOverride        string // Optional: pretty-printed state override for error display
}

// CallExecutionError wraps an error with call context.
// Mirrors viem's CallExecutionError with consistent formatting.
type CallExecutionError struct {
	Cause   error
	To      *common.Address
	Data    []byte
	ChainID *int64
	params  CallErrorParams
}

func (e *CallExecutionError) Error() string {
	// Use cause's message as base (preserves HttpRequestError, RpcRequestError etc.)
	base := "call execution failed"
	if e.Cause != nil {
		base = e.Cause.Error()
	}

	prettyArgs := e.formatRawCallArgs()
	if prettyArgs == "" {
		return "CallExecutionError: " + base
	}

	// Mirror viem format: error type prefix + cause message + blank line + Raw Call Arguments
	return "CallExecutionError: " + base + "\n\nRaw Call Arguments:\n" + prettyArgs
}

func (e *CallExecutionError) formatRawCallArgs() string {
	p := e.params
	var parts []string

	if p.From != nil {
		parts = append(parts, fmt.Sprintf("  from:   %s", p.From.Hex()))
	}
	if p.To != nil {
		parts = append(parts, fmt.Sprintf("  to:     %s", p.To.Hex()))
	}
	if p.Value != nil && p.Value.Sign() > 0 {
		valStr := unit.FormatEther(p.Value)
		symbol := p.NativeCurrencySymbol
		if symbol == "" {
			symbol = "ETH"
		}
		parts = append(parts, fmt.Sprintf("  value:  %s %s", valStr, symbol))
	}
	if len(p.Data) > 0 {
		hex := common.Bytes2Hex(p.Data)
		if len(hex) > 66 {
			hex = hex[:66] + "..."
		}
		dataLine := fmt.Sprintf("  data:   0x%s", hex)
		if p.FormattedData != "" {
			dataLine += fmt.Sprintf(" (%s)", p.FormattedData)
		}
		parts = append(parts, dataLine)
	}
	if p.Gas != nil {
		parts = append(parts, fmt.Sprintf("  gas:    %d", *p.Gas))
	}
	if p.GasPrice != nil && p.GasPrice.Sign() > 0 {
		parts = append(parts, fmt.Sprintf("  gasPrice: %s gwei", unit.FormatGwei(p.GasPrice)))
	}
	if p.MaxFeePerGas != nil && p.MaxFeePerGas.Sign() > 0 {
		parts = append(parts, fmt.Sprintf("  maxFeePerGas: %s gwei", unit.FormatGwei(p.MaxFeePerGas)))
	}
	if p.MaxPriorityFeePerGas != nil && p.MaxPriorityFeePerGas.Sign() > 0 {
		parts = append(parts, fmt.Sprintf("  maxPriorityFeePerGas: %s gwei", unit.FormatGwei(p.MaxPriorityFeePerGas)))
	}
	if p.ChainID != nil {
		parts = append(parts, fmt.Sprintf("  chainId: %d", *p.ChainID))
	}
	if p.StateOverride != "" {
		parts = append(parts, "  State Override:")
		parts = append(parts, p.StateOverride)
	}

	if len(parts) == 0 {
		// Fallback to minimal info
		if p.To != nil {
			parts = append(parts, fmt.Sprintf("  to: %s", p.To.Hex()))
		}
		if len(p.Data) > 0 {
			hex := common.Bytes2Hex(p.Data)
			if len(hex) > 66 {
				hex = hex[:66] + "..."
			}
			parts = append(parts, fmt.Sprintf("  data: 0x%s", hex))
		}
		if p.ChainID != nil {
			parts = append(parts, fmt.Sprintf("  chainId: %d", *p.ChainID))
		}
	}

	return strings.Join(parts, "\n")
}

func (e *CallExecutionError) Unwrap() error {
	return e.Cause
}

// GetCallError wraps an error with call execution context.
// Provides better error messages with relevant call details (mirrors viem's getCallError).
func GetCallError(err error, params CallErrorParams) error {
	if err == nil {
		return nil
	}

	return &CallExecutionError{
		Cause:   err,
		To:      params.To,
		Data:    params.Data,
		ChainID: params.ChainID,
		params:  params,
	}
}

// DataError is the interface for errors that contain revert data (e.g. RPCRequestError).
type DataError interface {
	ErrorData() interface{}
}

// GetRevertErrorData extracts revert data from an error if available.
// Walks the error chain to find errors implementing DataError (e.g. transport.RPCRequestError).
// Falls back to parsing hex from the error string for RPCs that embed data in the message.
// Returns nil if the error doesn't contain revert data.
func GetRevertErrorData(err error) []byte {
	if err == nil {
		return nil
	}
	var dataErr DataError
	if errors.As(err, &dataErr) {
		d := dataErr.ErrorData()
		if d != nil {
			if data, ok := d.([]byte); ok {
				return data
			}
			if s, ok := d.(string); ok {
				return common.FromHex(s)
			}
		}
	}

	// Fallback: parse hex from error string (some RPCs embed data in message)
	return parseHexFromErrorString(err.Error())
}

// parseHexFromErrorString extracts a hex string (0x...) from an error message.
func parseHexFromErrorString(s string) []byte {
	idx := strings.Index(s, "0x")
	if idx < 0 {
		idx = strings.Index(s, "0X")
	}
	if idx < 0 {
		return nil
	}
	hexStr := s[idx:]
	end := 2
	for i := 2; i < len(hexStr); i++ {
		c := hexStr[i]
		if (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F') {
			end = i + 1
		} else {
			break
		}
	}
	if end <= 2 {
		return nil
	}
	return common.FromHex(hexStr[:end])
}
