package errors

import "fmt"

// BaseFeeScalarError is returned when baseFeeMultiplier is invalid.
type BaseFeeScalarError struct{}

func (e *BaseFeeScalarError) Error() string {
	return "`baseFeeMultiplier` must be greater than 1."
}

// Eip1559FeesNotSupportedError is returned when chain does not support EIP-1559.
type Eip1559FeesNotSupportedError struct{}

func (e *Eip1559FeesNotSupportedError) Error() string {
	return "Chain does not support EIP-1559 fees."
}

// MaxFeePerGasTooLowError is returned when maxFeePerGas is less than maxPriorityFeePerGas.
type MaxFeePerGasTooLowError struct {
	MaxPriorityFeePerGas string
}

func (e *MaxFeePerGasTooLowError) Error() string {
	return fmt.Sprintf("`maxFeePerGas` cannot be less than the `maxPriorityFeePerGas` (%s gwei).", e.MaxPriorityFeePerGas)
}
