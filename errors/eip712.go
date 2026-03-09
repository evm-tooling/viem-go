package errors

import "fmt"

// Eip712DomainNotFoundError is returned when EIP-712 domain is not found on contract.
type Eip712DomainNotFoundError struct {
	Address string
}

func (e *Eip712DomainNotFoundError) Error() string {
	msg := fmt.Sprintf("No EIP-712 domain found on contract \"%s\".", e.Address)
	msg += "\nEnsure that:"
	msg += fmt.Sprintf("\n- The contract is deployed at the address \"%s\".", e.Address)
	msg += "\n- `eip712Domain()` function exists on the contract."
	msg += "\n- `eip712Domain()` function matches signature to ERC-5267 specification."
	return msg
}
