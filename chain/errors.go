package chain

import "fmt"

// ErrChainNotFound is returned when no chain was provided to the request.
var ErrChainNotFound = fmt.Errorf("chain: no chain was provided to the request")
var ErrInvalidChainsLen = fmt.Errorf("chain: no chains defined in chain array")
var ErrInvalidChainID = fmt.Errorf("chain: Invalid Chainid")

// ChainMismatchError is returned when the current chain ID does not match the expected chain.
type ChainMismatchError struct {
	Chain          Chain
	CurrentChainID int64
}

func (e *ChainMismatchError) Error() string {
	return fmt.Sprintf(
		"chain: the current chain of the wallet (id: %d) does not match the target chain for the transaction (id: %d – %s)",
		e.CurrentChainID, e.Chain.ID, e.Chain.Name,
	)
}
