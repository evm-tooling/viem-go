package chain

import (
	"github.com/ChefBingbong/viem-go/utils/address"
)

// ChainBlockExplorer represents a block explorer for a chain.
type ChainBlockExplorer struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	ApiURL string `json:"apiUrl,omitempty"`
}

// ChainContract represents a contract address on a chain.
type ChainContract struct {
	Address      address.Address `json:"address"`
	BlockCreated *uint64         `json:"blockCreated,omitempty"`
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

// Chain is the basic chain definition, mirroring viem's Chain type.
// It omits formatters, fees, serializers, and other chain config for simplicity.
type Chain struct {
	ID                              int64                         `json:"id"`
	Name                            string                        `json:"name"`
	NativeCurrency                  ChainNativeCurrency           `json:"nativeCurrency"`
	RpcUrls                         map[string]ChainRpcUrls       `json:"rpcUrls"`
	BlockExplorers                  map[string]ChainBlockExplorer `json:"blockExplorers,omitempty"`
	BlockTime                       *int64                        `json:"blockTime,omitempty"`
	Contracts                       map[string]ChainContract      `json:"contracts,omitempty"`
	EnsTlds                         []string                      `json:"ensTlds,omitempty"`
	SourceID                        *int64                        `json:"sourceId,omitempty"`
	Testnet                         bool                          `json:"testnet,omitempty"`
	ExperimentalPreconfirmationTime *int64                        `json:"experimental_preconfirmationTime,omitempty"`
}
