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

// GetEstimateGasErrorParams contains parameters for wrapping estimate gas errors.
// Mirrors viem's getEstimateGasError parameters (partial EstimateGasParameters).
type GetEstimateGasErrorParams struct {
	Account              *common.Address
	Chain                *chain.Chain
	Data                 []byte
	Gas                  *uint64
	GasPrice             *big.Int
	MaxFeePerGas         *big.Int
	MaxPriorityFeePerGas *big.Int
	Nonce                *uint64
	To                   *common.Address
	Value                *big.Int
}

// GetEstimateGasError wraps an error with estimate gas execution context.
// Classifies the cause via GetNodeError, then wraps in EstimateGasExecutionError.
// Mirrors viem's getEstimateGasError.
func GetEstimateGasError(err error, params GetEstimateGasErrorParams) error {
	if err == nil {
		return nil
	}

	nodeParams := GetNodeErrorParams{
		Gas:                  params.Gas,
		GasPrice:             params.GasPrice,
		MaxFeePerGas:         params.MaxFeePerGas,
		MaxPriorityFeePerGas: params.MaxPriorityFeePerGas,
		Nonce:                params.Nonce,
		Value:                params.Value,
	}
	if params.To != nil {
		s := params.To.Hex()
		nodeParams.To = &s
	}

	cause := GetNodeError(err, nodeParams)
	// If we got UnknownNodeError, use original err for cause (mirrors viem)
	if IsUnknownNodeError(cause) {
		cause = err
	}

	metaMessages := formatEstimateGasArguments(params)
	return &pkgerrors.EstimateGasExecutionError{
		Cause:        cause,
		MetaMessages: metaMessages,
	}
}

func formatEstimateGasArguments(p GetEstimateGasErrorParams) []string {
	var parts []string
	if p.Account != nil {
		parts = append(parts, fmt.Sprintf("  from:    %s", p.Account.Hex()))
	}
	if p.To != nil {
		parts = append(parts, fmt.Sprintf("  to:      %s", p.To.Hex()))
	}
	if p.Value != nil && p.Value.Sign() > 0 {
		symbol := "ETH"
		if p.Chain != nil && p.Chain.NativeCurrency.Symbol != "" {
			symbol = p.Chain.NativeCurrency.Symbol
		}
		parts = append(parts, fmt.Sprintf("  value:   %s %s", unit.FormatEther(p.Value), symbol))
	}
	if len(p.Data) > 0 {
		hex := fmt.Sprintf("0x%x", p.Data)
		if len(hex) > 66 {
			hex = hex[:66] + "..."
		}
		parts = append(parts, fmt.Sprintf("  data:    %s", hex))
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
	return []string{"Estimate Gas Arguments:", strings.Join(parts, "\n")}
}
