package definitions

import (
	"github.com/ChefBingbong/viem-go/chain"
	"github.com/ChefBingbong/viem-go/utils/address"
)

// Optimism is the OP Mainnet (Optimism) chain definition.
var Optimism = chain.DefineChain(chain.Chain{
	ID:   10,
	Name: "OP Mainnet",
	NativeCurrency: chain.ChainNativeCurrency{
		Name:     "Ether",
		Symbol:   "ETH",
		Decimals: 18,
	},
	BlockTime: int64Ptr(2_000),
	SourceID:  int64Ptr(1), // mainnet L1
	RpcUrls: map[string]chain.ChainRpcUrls{
		"default": {
			HTTP: []string{"https://mainnet.optimism.io"},
		},
	},
	BlockExplorers: map[string]chain.ChainBlockExplorer{
		"default": {
			Name:   "Optimism Explorer",
			URL:    "https://optimistic.etherscan.io",
			ApiURL: "https://api-optimistic.etherscan.io/api",
		},
	},
	Contracts: map[string]chain.ChainContract{
		"multicall3": {
			Address:      address.Address("0xca11bde05977b3631167028862be2a173976ca11"),
			BlockCreated: uint64Ptr(4_286_263),
		},
		"portal": {
			Address:      address.Address("0xbEb5Fc579115071764c7423A4f12eDde41f106Ed"),
			BlockCreated: nil, // L1 sourceId indexed
		},
		"l1StandardBridge": {
			Address:      address.Address("0x99C9fc46f92E8a1c0deC1b1747d010903E884bE1"),
			BlockCreated: nil,
		},
	},
})
