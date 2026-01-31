package chain

func ExtractChain(chains []*Chain, chainID int64) (*Chain, error) {
	if len(chains) == 0 {
		return nil, ErrChainNotFound
	}

	if chainID < 0 {
		return nil, ErrInvalidChainID
	}

	for _, chain := range chains {
		if chain.ID == chainID {
			targetChain := *chain
			return &targetChain, nil
		}
	}
	return nil, ErrChainNotFound
}
