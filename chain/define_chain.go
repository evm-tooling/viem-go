package chain

// DefineChain returns a chain definition, mirroring viem's defineChain.
// It copies the chain so callers can reuse the same struct without mutation.
// Use "default" as the key in RpcUrls and BlockExplorers for the primary entry.
func DefineChain(c Chain) Chain {
	// Copy so the returned chain is independent
	out := Chain{
		ID:                              c.ID,
		Name:                            c.Name,
		NativeCurrency:                  c.NativeCurrency,
		RpcUrls:                         copyRpcUrls(c.RpcUrls),
		BlockExplorers:                  copyBlockExplorers(c.BlockExplorers),
		BlockTime:                       copyInt64Ptr(c.BlockTime),
		Contracts:                       copyContracts(c.Contracts),
		EnsTlds:                         copyStrings(c.EnsTlds),
		SourceID:                        copyInt64Ptr(c.SourceID),
		Testnet:                         c.Testnet,
		ExperimentalPreconfirmationTime: copyInt64Ptr(c.ExperimentalPreconfirmationTime),
	}
	return out
}

// DefaultRpcUrl returns the first HTTP URL from the "default" RPC entry, or empty string if not set.
func (c *Chain) DefaultRpcUrl() string {
	if c.RpcUrls == nil {
		return ""
	}
	urls, ok := c.RpcUrls["default"]
	if !ok || len(urls.HTTP) == 0 {
		return ""
	}
	return urls.HTTP[0]
}

// DefaultBlockExplorer returns the block explorer for the "default" key, or empty struct if not set.
func (c *Chain) DefaultBlockExplorer() ChainBlockExplorer {
	if c.BlockExplorers == nil {
		return ChainBlockExplorer{}
	}
	return c.BlockExplorers["default"]
}

func copyRpcUrls(m map[string]ChainRpcUrls) map[string]ChainRpcUrls {
	if m == nil {
		return nil
	}
	out := make(map[string]ChainRpcUrls, len(m))
	for k, v := range m {
		out[k] = ChainRpcUrls{
			HTTP:      copyStrings(v.HTTP),
			WebSocket: copyStrings(v.WebSocket),
		}
	}
	return out
}

func copyBlockExplorers(m map[string]ChainBlockExplorer) map[string]ChainBlockExplorer {
	if m == nil {
		return nil
	}
	out := make(map[string]ChainBlockExplorer, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func copyContracts(m map[string]ChainContract) map[string]ChainContract {
	if m == nil {
		return nil
	}
	out := make(map[string]ChainContract, len(m))
	for k, v := range m {
		cp := ChainContract{Address: v.Address}
		if v.BlockCreated != nil {
			bc := *v.BlockCreated
			cp.BlockCreated = &bc
		}
		out[k] = cp
	}
	return out
}

func copyStrings(s []string) []string {
	if s == nil {
		return nil
	}
	out := make([]string, len(s))
	copy(out, s)
	return out
}

func copyInt64Ptr(p *int64) *int64 {
	if p == nil {
		return nil
	}
	v := *p
	return &v
}
