package wallet

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ChefBingbong/viem-go/chain"
	"github.com/ChefBingbong/viem-go/utils/encoding"
)

// AddChainParameters contains the parameters for the AddChain action.
// This mirrors viem's AddChainParameters type.
type AddChainParameters struct {
	// Chain is the chain to add to the wallet.
	Chain *chain.Chain
}

// addEthereumChainParams is the internal request format for wallet_addEthereumChain.
type addEthereumChainParams struct {
	ChainID           string                    `json:"chainId"`
	ChainName         string                    `json:"chainName"`
	NativeCurrency    chain.ChainNativeCurrency `json:"nativeCurrency"`
	RpcUrls           []string                  `json:"rpcUrls"`
	BlockExplorerUrls []string                  `json:"blockExplorerUrls,omitempty"`
}

// AddChain adds an EVM chain to the wallet.
//
// This is equivalent to viem's `addChain` action.
//
// JSON-RPC Method: wallet_addEthereumChain (EIP-3085)
//
// Example:
//
//	err := wallet.AddChain(ctx, client, wallet.AddChainParameters{
//	    Chain: optimism,
//	})
func AddChain(ctx context.Context, client Client, params AddChainParameters) error {
	ch := params.Chain
	if ch == nil {
		return fmt.Errorf("chain is required")
	}

	// Convert chain ID to hex using existing encoding utility
	chainIDHex := encoding.NumberToHex(big.NewInt(ch.ID))

	// Build RPC URLs from default entry
	var rpcUrls []string
	if defaultUrls, ok := ch.RpcUrls["default"]; ok {
		rpcUrls = defaultUrls.HTTP
	}

	// Build block explorer URLs
	var blockExplorerUrls []string
	if ch.BlockExplorers != nil {
		for _, explorer := range ch.BlockExplorers {
			blockExplorerUrls = append(blockExplorerUrls, explorer.URL)
		}
	}

	reqParams := addEthereumChainParams{
		ChainID:        chainIDHex,
		ChainName:      ch.Name,
		NativeCurrency: ch.NativeCurrency,
		RpcUrls:        rpcUrls,
	}
	if len(blockExplorerUrls) > 0 {
		reqParams.BlockExplorerUrls = blockExplorerUrls
	}

	_, err := client.Request(ctx, "wallet_addEthereumChain", reqParams)
	if err != nil {
		return fmt.Errorf("wallet_addEthereumChain failed: %w", err)
	}

	return nil
}
