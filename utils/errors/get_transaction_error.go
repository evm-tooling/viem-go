// Package errors provides error handling utilities.
package errors

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/chain"
	pkgerrors "github.com/ChefBingbong/viem-go/errors"
	"github.com/ChefBingbong/viem-go/utils/unit"
)

// GetTransactionErrorParams contains parameters for wrapping transaction errors.
// Mirrors viem's GetTransactionErrorParameters (partial SendTransactionParameters).
type GetTransactionErrorParams struct {
	Account              *common.Address
	Chain                *chain.Chain
	Data                 string
	Gas                  *uint64
	GasPrice             *big.Int
	MaxFeePerGas         *big.Int
	MaxPriorityFeePerGas *big.Int
	Nonce                *uint64
	To                   *string
	Value                *big.Int
}

// GetTransactionError wraps an error with transaction execution context.
// Classifies the cause via GetNodeError, then wraps in TransactionExecutionError.
// Mirrors viem's getTransactionError.
func GetTransactionError(err error, params GetTransactionErrorParams) error {
	if err == nil {
		return nil
	}

	nodeParams := GetNodeErrorParams{
		Gas:                  params.Gas,
		GasPrice:             params.GasPrice,
		MaxFeePerGas:         params.MaxFeePerGas,
		MaxPriorityFeePerGas: params.MaxPriorityFeePerGas,
		Nonce:                params.Nonce,
		To:                   params.To,
		Value:                params.Value,
	}

	cause := GetNodeError(err, nodeParams)
	// If we got UnknownNodeError, use original err for cause (mirrors viem)
	if IsUnknownNodeError(cause) {
		cause = err
	}

	metaMessages := formatRequestArguments(params)
	return &pkgerrors.TransactionExecutionError{
		Cause:        cause,
		MetaMessages: metaMessages,
	}
}

func formatRequestArguments(p GetTransactionErrorParams) []string {
	var parts []string
	if p.Chain != nil {
		parts = append(parts, fmt.Sprintf("  chain:   %s (id: %d)", p.Chain.Name, p.Chain.ID))
	}
	if p.Account != nil {
		parts = append(parts, fmt.Sprintf("  from:    %s", p.Account.Hex()))
	}
	if p.To != nil {
		parts = append(parts, fmt.Sprintf("  to:      %s", *p.To))
	}
	if p.Value != nil && p.Value.Sign() > 0 {
		symbol := "ETH"
		if p.Chain != nil && p.Chain.NativeCurrency.Symbol != "" {
			symbol = p.Chain.NativeCurrency.Symbol
		}
		parts = append(parts, fmt.Sprintf("  value:   %s %s", unit.FormatEther(p.Value), symbol))
	}
	if p.Data != "" {
		data := p.Data
		if len(data) > 66 {
			data = data[:66] + "..."
		}
		parts = append(parts, fmt.Sprintf("  data:    %s", data))
	}
	if p.Gas != nil {
		parts = append(parts, fmt.Sprintf("  gas:     %d", *p.Gas))
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
	if p.Nonce != nil {
		parts = append(parts, fmt.Sprintf("  nonce:   %d", *p.Nonce))
	}

	if len(parts) == 0 {
		return nil
	}
	return []string{"Request Arguments:", strings.Join(parts, "\n")}
}
