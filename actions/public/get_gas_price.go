package public

import (
	"context"
	"fmt"
	"math/big"

	json "github.com/goccy/go-json"
)

// GetGasPriceReturnType is the return type for the GetGasPrice action.
// It represents the gas price in wei.
type GetGasPriceReturnType = *big.Int

// GetGasPrice returns the current price of gas (in wei).
//
// This is equivalent to viem's `getGasPrice` action.
//
// JSON-RPC Method: eth_gasPrice
//
// Example:
//
//	gasPrice, err := public.GetGasPrice(ctx, client)
//	// gasPrice is in wei, use formatGwei/formatEther to convert
func GetGasPrice(ctx context.Context, client Client) (GetGasPriceReturnType, error) {
	// Execute the request
	resp, err := client.Request(ctx, "eth_gasPrice")
	if err != nil {
		return nil, fmt.Errorf("eth_gasPrice failed: %w", err)
	}

	var hexGasPrice string
	if unmarshalErr := json.Unmarshal(resp.Result, &hexGasPrice); unmarshalErr != nil {
		return nil, fmt.Errorf("failed to unmarshal gas price: %w", unmarshalErr)
	}

	// Parse the gas price
	gasPrice, parseErr := parseHexBigInt(hexGasPrice)
	if parseErr != nil {
		return nil, fmt.Errorf("failed to parse gas price: %w", parseErr)
	}

	return gasPrice, nil
}
