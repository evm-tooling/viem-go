package definitions

import (
	"github.com/ChefBingbong/viem-go/chain"
	"github.com/ChefBingbong/viem-go/utils/address"
)

// Avalanche is the Avalanche C-Chain definition.
var Avalanche = chain.DefineChain(chain.Chain{
	ID:   43_114,
	Name: "Avalanche",
	NativeCurrency: chain.ChainNativeCurrency{
		Name:     "Avalanche",
		Symbol:   "AVAX",
		Decimals: 18,
	},
	BlockTime: int64Ptr(1_700),
	RpcUrls: map[string]chain.ChainRpcUrls{
		"default": {
			HTTP: []string{"https://api.avax.network/ext/bc/C/rpc"},
		},
	},
	BlockExplorers: map[string]chain.ChainBlockExplorer{
		"default": {
			Name:   "SnowTrace",
			URL:    "https://snowtrace.io",
			ApiURL: "https://api.snowtrace.io",
		},
	},
	Contracts: map[string]chain.ChainContract{
		"multicall3": {
			Address:      address.Address("0xca11bde05977b3631167028862be2a173976ca11"),
			BlockCreated: uint64Ptr(11_907_934),
		},
	},
})
