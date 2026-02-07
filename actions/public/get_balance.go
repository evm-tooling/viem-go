package public

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	json "github.com/goccy/go-json"

	"github.com/ethereum/go-ethereum/common"
)

// GetBalanceParameters contains the parameters for the GetBalance action.
// This mirrors viem's GetBalanceParameters type.
type GetBalanceParameters struct {
	// Address is the address to get the balance of. Required.
	Address common.Address

	// BlockNumber is the block number to get the balance at.
	// Mutually exclusive with BlockTag.
	BlockNumber *uint64

	// BlockTag is the block tag to get the balance at (e.g., "latest", "pending").
	// Mutually exclusive with BlockNumber.
	BlockTag BlockTag
}

// GetBalanceReturnType is the return type for the GetBalance action.
// It represents the balance in wei.
type GetBalanceReturnType = *big.Int

// GetBalance returns the balance of an address in wei.
//
// This is equivalent to viem's `getBalance` action.
//
// JSON-RPC Method: eth_getBalance
//
// Example:
//
//	balance, err := public.GetBalance(ctx, client, public.GetBalanceParameters{
//	    Address: common.HexToAddress("0x..."),
//	})
//	// balance is in wei, use formatEther to convert to ETH
func GetBalance(ctx context.Context, client Client, params GetBalanceParameters) (GetBalanceReturnType, error) {
	// Determine block tag
	blockTag := resolveBlockTag(client, params.BlockNumber, params.BlockTag)

	// Execute the request
	resp, err := client.Request(ctx, "eth_getBalance", params.Address.Hex(), blockTag)
	if err != nil {
		return nil, fmt.Errorf("eth_getBalance failed: %w", err)
	}

	var hexBalance string
	if unmarshalErr := json.Unmarshal(resp.Result, &hexBalance); unmarshalErr != nil {
		return nil, fmt.Errorf("failed to unmarshal balance: %w", unmarshalErr)
	}

	// Parse the balance
	balance, err := parseHexBigInt(hexBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to parse balance: %w", err)
	}

	return balance, nil
}

// parseHexBigInt parses a hex string to *big.Int.
func parseHexBigInt(hexStr string) (*big.Int, error) {
	hexStr = strings.TrimPrefix(hexStr, "0x")
	if hexStr == "" {
		return big.NewInt(0), nil
	}

	n := new(big.Int)
	_, ok := n.SetString(hexStr, 16)
	if !ok {
		return nil, fmt.Errorf("invalid hex string: %s", hexStr)
	}
	return n, nil
}
