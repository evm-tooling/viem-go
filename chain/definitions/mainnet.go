package definitions

import (
	"github.com/ChefBingbong/viem-go/chain"
	"github.com/ChefBingbong/viem-go/utils/address"
)

// Mainnet is the Ethereum mainnet chain definition.
var Mainnet = chain.DefineChain(chain.Chain{
	ID:   1,
	Name: "Ethereum",
	NativeCurrency: chain.ChainNativeCurrency{
		Name:     "Ether",
		Symbol:   "ETH",
		Decimals: 18,
	},
	BlockTime: int64Ptr(12_000),
	RpcUrls: map[string]chain.ChainRpcUrls{
		"default": {
			HTTP: []string{"https://eth.merkle.io"},
		},
	},
	BlockExplorers: map[string]chain.ChainBlockExplorer{
		"default": {
			Name:   "Etherscan",
			URL:    "https://etherscan.io",
			ApiURL: "https://api.etherscan.io/api",
		},
	},
	Contracts: map[string]chain.ChainContract{
		"ensUniversalResolver": {
			Address:      address.Address("0xeeeeeeee14d718c2b47d9923deab1335e144eeee"),
			BlockCreated: uint64Ptr(23_085_558),
		},
		"multicall3": {
			Address:      address.Address("0xca11bde05977b3631167028862be2a173976ca11"),
			BlockCreated: uint64Ptr(14_353_601),
		},
	},
})
