package types

import (
	"github.com/ethereum/go-ethereum/common"
)

// ChainBlockExplorer represents a block explorer for a chain.
type ChainBlockExplorer struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	ApiURL string `json:"apiUrl,omitempty"`
}

// ChainContract represents a contract address on a chain.
type ChainContract struct {
	Address      common.Address `json:"address"`
	BlockCreated *uint64        `json:"blockCreated,omitempty"`
}

// ChainNativeCurrency represents the native currency of a chain.
type ChainNativeCurrency struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals uint8  `json:"decimals"`
}

// ChainRpcUrls represents RPC endpoints for a chain.
type ChainRpcUrls struct {
	HTTP      []string `json:"http"`
	WebSocket []string `json:"webSocket,omitempty"`
}

// ChainContracts contains well-known contract addresses.
type ChainContracts struct {
	Multicall3           *ChainContract `json:"multicall3,omitempty"`
	EnsRegistry          *ChainContract `json:"ensRegistry,omitempty"`
	EnsUniversalResolver *ChainContract `json:"ensUniversalResolver,omitempty"`
}

// Chain is the basic chain definition, mirroring viem's Chain type.
type Chain struct {
	ID                              int64                         `json:"id"`
	Name                            string                        `json:"name"`
	NativeCurrency                  ChainNativeCurrency           `json:"nativeCurrency"`
	RpcUrls                         map[string]ChainRpcUrls       `json:"rpcUrls"`
	BlockExplorers                  map[string]ChainBlockExplorer `json:"blockExplorers,omitempty"`
	BlockTime                       *int64                        `json:"blockTime,omitempty"`
	Contracts                       *ChainContracts               `json:"contracts,omitempty"`
	EnsTlds                         []string                      `json:"ensTlds,omitempty"`
	SourceID                        *int64                        `json:"sourceId,omitempty"`
	Testnet                         bool                          `json:"testnet,omitempty"`
	ExperimentalPreconfirmationTime *int64                        `json:"experimental_preconfirmationTime,omitempty"`
}

// GetDefaultHTTPUrl returns the first default HTTP RPC URL, or empty string if none.
func (c *Chain) GetDefaultHTTPUrl() string {
	if urls, ok := c.RpcUrls["default"]; ok && len(urls.HTTP) > 0 {
		return urls.HTTP[0]
	}
	return ""
}

// GetDefaultWSUrl returns the first default WebSocket RPC URL, or empty string if none.
func (c *Chain) GetDefaultWSUrl() string {
	if urls, ok := c.RpcUrls["default"]; ok && len(urls.WebSocket) > 0 {
		return urls.WebSocket[0]
	}
	return ""
}

// GetBlockExplorer returns the default block explorer, or nil if none.
func (c *Chain) GetBlockExplorer() *ChainBlockExplorer {
	if explorer, ok := c.BlockExplorers["default"]; ok {
		return &explorer
	}
	return nil
}
