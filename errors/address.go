package errors

import "fmt"

// InvalidAddressError is returned when an address is invalid.
type InvalidAddressError struct {
	Address string
}

func (e *InvalidAddressError) Error() string {
	msg := fmt.Sprintf("Address \"%s\" is invalid.", e.Address)
	msg += "\n- Address must be a hex value of 20 bytes (40 hex characters)."
	msg += "\n- Address must match its checksum counterpart."
	return msg
}
