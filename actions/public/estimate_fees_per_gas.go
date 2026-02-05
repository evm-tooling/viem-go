package public

import (
	"context"
	"fmt"
	"math"
	"math/big"
)

// FeeValuesType represents the type of fee values to return.
type FeeValuesType string

const (
	// FeeValuesTypeLegacy returns a legacy gas price.
	FeeValuesTypeLegacy FeeValuesType = "legacy"
	// FeeValuesTypeEIP1559 returns EIP-1559 fee values.
	FeeValuesTypeEIP1559 FeeValuesType = "eip1559"
)

// EstimateFeesPerGasParameters contains the parameters for the
// EstimateFeesPerGas action.
//
// This mirrors the behavior of viem's `estimateFeesPerGas`:
//   - For EIP-1559 chains, returns maxFeePerGas & maxPriorityFeePerGas.
//   - For legacy chains, returns gasPrice.
//   - Applies a base fee multiplier (default 1.2) to provide a safety buffer.
type EstimateFeesPerGasParameters struct {
	// Type is the type of fee values to return.
	// Defaults to FeeValuesTypeEIP1559.
	Type FeeValuesType

	// BaseFeeMultiplier is the multiplier applied to the base fee per gas
	// (or gas price for legacy chains) when computing fees.
	// Defaults to 1.2 (20% buffer).
	BaseFeeMultiplier *float64
}

// EstimateFeesPerGasReturnType represents the estimated fees per gas.
//
// For Type == FeeValuesTypeEIP1559:
//   - MaxFeePerGas and MaxPriorityFeePerGas are populated.
//
// For Type == FeeValuesTypeLegacy:
//   - GasPrice is populated.
type EstimateFeesPerGasReturnType struct {
	Type                 FeeValuesType
	MaxFeePerGas         *big.Int
	MaxPriorityFeePerGas *big.Int
	GasPrice             *big.Int
}

// EstimateFeesPerGas returns an estimate for the fees per gas (in wei) for a
// transaction to be likely included in the next block.
//
// This is equivalent to viem's `estimateFeesPerGas` action.
func EstimateFeesPerGas(
	ctx context.Context,
	client Client,
	params EstimateFeesPerGasParameters,
) (*EstimateFeesPerGasReturnType, error) {
	feeType := params.Type
	if feeType == "" {
		feeType = FeeValuesTypeEIP1559
	}

	// Resolve multiplier (default 1.2).
	baseFeeMultiplier := 1.2
	if params.BaseFeeMultiplier != nil {
		baseFeeMultiplier = *params.BaseFeeMultiplier
	}
	if baseFeeMultiplier < 1 {
		return nil, &BaseFeeScalarError{Multiplier: baseFeeMultiplier}
	}

	// Fetch latest block once (used by both paths).
	block, err := GetBlock(ctx, client, GetBlockParameters{
		BlockTag: BlockTagLatest,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest block: %w", err)
	}

	switch feeType {
	case FeeValuesTypeEIP1559:
		if block.BaseFeePerGas == nil {
			return nil, fmt.Errorf("EIP-1559 fees not supported: missing baseFeePerGas on block")
		}

		maxPriorityFeePerGas, err := EstimateMaxPriorityFeePerGas(ctx, client, EstimateMaxPriorityFeePerGasParameters{
			Block: block,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to estimate maxPriorityFeePerGas: %w", err)
		}

		baseFeePerGas := applyBaseFeeMultiplier(block.BaseFeePerGas, baseFeeMultiplier)
		maxFeePerGas := new(big.Int).Add(baseFeePerGas, maxPriorityFeePerGas)

		return &EstimateFeesPerGasReturnType{
			Type:                 FeeValuesTypeEIP1559,
			MaxFeePerGas:         maxFeePerGas,
			MaxPriorityFeePerGas: maxPriorityFeePerGas,
		}, nil

	case FeeValuesTypeLegacy:
		gasPrice, err := GetGasPrice(ctx, client)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch gas price: %w", err)
		}

		adjustedGasPrice := applyBaseFeeMultiplier(gasPrice, baseFeeMultiplier)
		return &EstimateFeesPerGasReturnType{
			Type:     FeeValuesTypeLegacy,
			GasPrice: adjustedGasPrice,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported fee values type: %s", feeType)
	}
}

// applyBaseFeeMultiplier applies the base fee multiplier using integer math
// to avoid floating point precision issues.
func applyBaseFeeMultiplier(base *big.Int, multiplier float64) *big.Int {
	if base == nil {
		return nil
	}

	// Determine decimal precision of multiplier (up to 18 decimals).
	decimals := 0
	for f := multiplier; f != math.Trunc(f) && decimals < 18; {
		f *= 10
		decimals++
	}

	denominator := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	numerator := new(big.Int).Mul(
		big.NewInt(int64(math.Round(multiplier*math.Pow(10, float64(decimals))))),
		big.NewInt(1),
	)

	result := new(big.Int).Mul(base, numerator)
	result.Div(result, denominator)
	return result
}
