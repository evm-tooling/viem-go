package wallet

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ChefBingbong/viem-go/utils/encoding"
)

// SwitchChainParameters contains the parameters for the SwitchChain action.
// This mirrors viem's SwitchChainParameters type.
type SwitchChainParameters struct {
	// ID is the chain ID to switch to.
	ID int64
}

// switchEthereumChainParams is the internal request format for wallet_switchEthereumChain.
type switchEthereumChainParams struct {
	ChainID string `json:"chainId"`
}

// SwitchChain switches the target chain in a wallet.
//
// This is equivalent to viem's `switchChain` action.
//
// JSON-RPC Method: wallet_switchEthereumChain (EIP-3326)
//
// Example:
//
//	err := wallet.SwitchChain(ctx, client, wallet.SwitchChainParameters{
//	    ID: 10, // Optimism
//	})
func SwitchChain(ctx context.Context, client Client, params SwitchChainParameters) error {
	chainIDHex := encoding.NumberToHex(big.NewInt(params.ID))

	reqParams := switchEthereumChainParams{
		ChainID: chainIDHex,
	}

	_, err := client.Request(ctx, "wallet_switchEthereumChain", reqParams)
	if err != nil {
		return fmt.Errorf("wallet_switchEthereumChain failed: %w", err)
	}

	return nil
}
