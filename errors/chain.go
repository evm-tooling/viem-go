package errors

import "fmt"

// ChainDoesNotSupportContract is returned when a chain does not support a required contract.
type ChainDoesNotSupportContract struct {
	ChainID      int64
	ChainName    string
	ContractName string
	BlockNumber  *uint64
	BlockCreated *uint64
}

func (e *ChainDoesNotSupportContract) Error() string {
	msg := fmt.Sprintf("Chain \"%s\" does not support contract \"%s\".", e.ChainName, e.ContractName)
	msg += "\nThis could be due to any of the following:"
	if e.BlockNumber != nil && e.BlockCreated != nil && *e.BlockCreated > *e.BlockNumber {
		msg += fmt.Sprintf("\n- The contract \"%s\" was not deployed until block %d (current block %d).", e.ContractName, *e.BlockCreated, *e.BlockNumber)
	} else {
		msg += fmt.Sprintf("\n- The chain does not have the contract \"%s\" configured.", e.ContractName)
	}
	return msg
}

// ChainMismatchError is returned when wallet chain does not match target chain.
type ChainMismatchError struct {
	ChainID      int64
	ChainName    string
	CurrentChainID int64
}

func (e *ChainMismatchError) Error() string {
	msg := fmt.Sprintf("The current chain of the wallet (id: %d) does not match the target chain for the transaction (id: %d – %s).", e.CurrentChainID, e.ChainID, e.ChainName)
	msg += fmt.Sprintf("\nCurrent Chain ID:  %d", e.CurrentChainID)
	msg += fmt.Sprintf("\nExpected Chain ID: %d – %s", e.ChainID, e.ChainName)
	return msg
}

// ChainNotFoundError is returned when no chain was provided to the request.
type ChainNotFoundError struct{}

func (e *ChainNotFoundError) Error() string {
	return "No chain was provided to the request.\nPlease provide a chain with the `chain` argument on the Action, or by supplying a `chain` to WalletClient."
}

// ClientChainNotConfiguredError is returned when no chain was provided to the Client.
type ClientChainNotConfiguredError struct{}

func (e *ClientChainNotConfiguredError) Error() string {
	return "No chain was provided to the Client."
}

// InvalidChainIdError is returned when chain ID is invalid.
type InvalidChainIdError struct {
	ChainID *int64
}

func (e *InvalidChainIdError) Error() string {
	if e.ChainID != nil {
		return fmt.Sprintf("Chain ID \"%d\" is invalid.", *e.ChainID)
	}
	return "Chain ID is invalid."
}
