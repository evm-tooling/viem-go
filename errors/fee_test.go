package errors

import "testing"

func TestBaseFeeScalarError_Error(t *testing.T) {
	err := &BaseFeeScalarError{}
	testErrorContains(t, err, []string{"baseFeeMultiplier", "greater than 1"})
}

func TestEip1559FeesNotSupportedError_Error(t *testing.T) {
	err := &Eip1559FeesNotSupportedError{}
	testErrorContains(t, err, []string{"EIP-1559"})
}

func TestMaxFeePerGasTooLowError_Error(t *testing.T) {
	err := &MaxFeePerGasTooLowError{MaxPriorityFeePerGas: "30"}
	testErrorContains(t, err, []string{"maxFeePerGas", "maxPriorityFeePerGas", "30"})
}
