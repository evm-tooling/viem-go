package wallet

import (
	"context"
	"encoding/json"
	"fmt"
)

// WatchAssetOptions contains the token options for the WatchAsset action.
type WatchAssetOptions struct {
	// Address is the address of the token contract.
	Address string `json:"address"`
	// Symbol is a ticker symbol or shorthand, up to 11 characters.
	Symbol string `json:"symbol"`
	// Decimals is the number of token decimals.
	Decimals uint8 `json:"decimals"`
	// Image is a string URL of the token logo.
	Image string `json:"image,omitempty"`
}

// WatchAssetParameters contains the parameters for the WatchAsset action.
// This mirrors viem's WatchAssetParams type from EIP-747.
type WatchAssetParameters struct {
	// Type is the token type (e.g., "ERC20").
	Type string `json:"type"`
	// Options contains the token details.
	Options WatchAssetOptions `json:"options"`
}

// WatchAssetReturnType is the return type for the WatchAsset action.
type WatchAssetReturnType = bool

// WatchAsset requests that the wallet tracks a specified token.
//
// This is equivalent to viem's `watchAsset` action.
//
// JSON-RPC Method: wallet_watchAsset (EIP-747)
//
// Example:
//
//	success, err := wallet.WatchAsset(ctx, client, wallet.WatchAssetParameters{
//	    Type: "ERC20",
//	    Options: wallet.WatchAssetOptions{
//	        Address:  "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
//	        Decimals: 18,
//	        Symbol:   "WETH",
//	    },
//	})
func WatchAsset(ctx context.Context, client Client, params WatchAssetParameters) (WatchAssetReturnType, error) {
	resp, err := client.Request(ctx, "wallet_watchAsset", params)
	if err != nil {
		return false, fmt.Errorf("wallet_watchAsset failed: %w", err)
	}

	var added bool
	if unmarshalErr := json.Unmarshal(resp.Result, &added); unmarshalErr != nil {
		return false, fmt.Errorf("failed to unmarshal watch asset result: %w", unmarshalErr)
	}

	return added, nil
}
