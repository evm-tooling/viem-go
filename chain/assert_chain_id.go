package chain

// AssertCurrentChain validates that the current chain ID matches the expected chain.
// Returns ErrChainNotFound if chain is nil.
// Returns *ChainMismatchError if currentChainID does not match chain.ID.
//
// Equivalent to viem's assertCurrentChain.
func AssertCurrentChain(chain *Chain, currentChainID int64) error {
	if chain == nil {
		return ErrChainNotFound
	}
	if currentChainID != chain.ID {
		return &ChainMismatchError{
			Chain:          *chain,
			CurrentChainID: currentChainID,
		}
	}
	return nil
}
