package errors

import (
	"strings"
	"testing"
)

func testErrorContains(t *testing.T, err error, contains []string) {
	t.Helper()
	got := err.Error()
	for _, sub := range contains {
		if !strings.Contains(got, sub) {
			t.Errorf("Error() = %q, want to contain %q", got, sub)
		}
	}
}

func TestAbiConstructorNotFoundError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *AbiConstructorNotFoundError
		contains []string
	}{
		{"basic", &AbiConstructorNotFoundError{}, []string{"constructor was not found"}},
		{"with docs", &AbiConstructorNotFoundError{DocsPath: "/docs/abi"}, []string{"constructor was not found", "viem.sh/docs/abi"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { testErrorContains(t, tt.err, tt.contains) })
	}
}

func TestAbiConstructorParamsNotFoundError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *AbiConstructorParamsNotFoundError
		contains []string
	}{
		{"basic", &AbiConstructorParamsNotFoundError{}, []string{"constructor parameters"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { testErrorContains(t, tt.err, tt.contains) })
	}
}

func TestAbiDecodingDataSizeInvalidError_Error(t *testing.T) {
	err := &AbiDecodingDataSizeInvalidError{Data: "0x1234", Size: 10}
	testErrorContains(t, err, []string{"invalid", "32 bytes", "0x1234", "10"})
}

func TestAbiDecodingDataSizeTooSmallError_Error(t *testing.T) {
	err := &AbiDecodingDataSizeTooSmallError{Data: "0x12", Params: "uint256", Size: 2}
	testErrorContains(t, err, []string{"too small for given parameters", "uint256", "0x12"})
}

func TestAbiDecodingZeroDataError_Error(t *testing.T) {
	err := &AbiDecodingZeroDataError{}
	testErrorContains(t, err, []string{"Cannot decode zero data"})
}

func TestAbiEncodingArrayLengthMismatchError_Error(t *testing.T) {
	err := &AbiEncodingArrayLengthMismatchError{ExpectedLength: 2, GivenLength: 3, Type: "uint256[]"}
	testErrorContains(t, err, []string{"length mismatch", "uint256[]", "2", "3"})
}

func TestAbiEncodingBytesSizeMismatchError_Error(t *testing.T) {
	err := &AbiEncodingBytesSizeMismatchError{ExpectedSize: 32, Value: "0x1234", ActualSize: 2}
	testErrorContains(t, err, []string{"does not match expected size", "0x1234", "bytes32", "bytes2"})
}

func TestAbiEncodingLengthMismatchError_Error(t *testing.T) {
	err := &AbiEncodingLengthMismatchError{ExpectedLength: 1, GivenLength: 2}
	testErrorContains(t, err, []string{"params/values length mismatch", "1", "2"})
}

func TestAbiErrorInputsNotFoundError_Error(t *testing.T) {
	err := &AbiErrorInputsNotFoundError{ErrorName: "MyError", DocsPath: "/docs"}
	testErrorContains(t, err, []string{"MyError", "inputs"})
}

func TestAbiErrorNotFoundError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *AbiErrorNotFoundError
		contains []string
	}{
		{"with name", &AbiErrorNotFoundError{ErrorName: "Foo"}, []string{"Foo", "not found on ABI"}},
		{"without name", &AbiErrorNotFoundError{}, []string{"not found on ABI"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { testErrorContains(t, tt.err, tt.contains) })
	}
}

func TestAbiErrorSignatureNotFoundError_Error(t *testing.T) {
	err := &AbiErrorSignatureNotFoundError{Signature: "0xabcd1234", DocsPath: "/docs"}
	testErrorContains(t, err, []string{"0xabcd1234", "not found on ABI", "4byte.sourcify.dev"})
}

func TestAbiEventSignatureEmptyTopicsError_Error(t *testing.T) {
	err := &AbiEventSignatureEmptyTopicsError{}
	testErrorContains(t, err, []string{"empty topics"})
}

func TestAbiEventSignatureNotFoundError_Error(t *testing.T) {
	err := &AbiEventSignatureNotFoundError{Signature: "0xdeadbeef"}
	testErrorContains(t, err, []string{"0xdeadbeef", "not found on ABI"})
}

func TestAbiEventNotFoundError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *AbiEventNotFoundError
		contains []string
	}{
		{"with name", &AbiEventNotFoundError{EventName: "Transfer"}, []string{"Transfer", "not found on ABI"}},
		{"without name", &AbiEventNotFoundError{}, []string{"not found on ABI"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { testErrorContains(t, tt.err, tt.contains) })
	}
}

func TestAbiFunctionNotFoundError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *AbiFunctionNotFoundError
		contains []string
	}{
		{"with name", &AbiFunctionNotFoundError{FunctionName: "balanceOf"}, []string{"balanceOf", "not found on ABI"}},
		{"without name", &AbiFunctionNotFoundError{}, []string{"not found on ABI"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { testErrorContains(t, tt.err, tt.contains) })
	}
}

func TestAbiFunctionOutputsNotFoundError_Error(t *testing.T) {
	err := &AbiFunctionOutputsNotFoundError{FunctionName: "getData"}
	testErrorContains(t, err, []string{"getData", "outputs"})
}

func TestAbiFunctionSignatureNotFoundError_Error(t *testing.T) {
	err := &AbiFunctionSignatureNotFoundError{Signature: "0x12345678"}
	testErrorContains(t, err, []string{"0x12345678", "not found on ABI"})
}

func TestAbiItemAmbiguityError_Error(t *testing.T) {
	err := &AbiItemAmbiguityError{TypeX: "uint256", ItemX: "foo()", TypeY: "uint8", ItemY: "bar()"}
	testErrorContains(t, err, []string{"ambiguous types", "uint256", "foo()", "uint8", "bar()"})
}

func TestBytesSizeMismatchError_Error(t *testing.T) {
	err := &BytesSizeMismatchError{ExpectedSize: 32, GivenSize: 20}
	testErrorContains(t, err, []string{"Expected bytes32", "got bytes20"})
}

func TestDecodeLogDataMismatch_Error(t *testing.T) {
	err := &DecodeLogDataMismatch{Data: "0x12", Params: "address", Size: 2}
	testErrorContains(t, err, []string{"too small for non-indexed", "address", "0x12"})
}

func TestDecodeLogTopicsMismatch_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *DecodeLogTopicsMismatch
		contains []string
	}{
		{"with param", &DecodeLogTopicsMismatch{ParamName: "owner", EventSig: "Transfer(address)"}, []string{"Expected a topic", "owner", "Transfer(address)"}},
		{"without param", &DecodeLogTopicsMismatch{EventSig: "Transfer()"}, []string{"Expected a topic", "Transfer()"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { testErrorContains(t, tt.err, tt.contains) })
	}
}

func TestInvalidAbiEncodingTypeError_Error(t *testing.T) {
	err := &InvalidAbiEncodingTypeError{Type: "invalid"}
	testErrorContains(t, err, []string{"not a valid encoding type", "invalid"})
}

func TestInvalidAbiDecodingTypeError_Error(t *testing.T) {
	err := &InvalidAbiDecodingTypeError{Type: "bad"}
	testErrorContains(t, err, []string{"not a valid decoding type", "bad"})
}

func TestInvalidArrayError_Error(t *testing.T) {
	err := &InvalidArrayError{Value: "not-an-array"}
	testErrorContains(t, err, []string{"not a valid array", "not-an-array"})
}

func TestInvalidDefinitionTypeError_Error(t *testing.T) {
	err := &InvalidDefinitionTypeError{Type: "foo"}
	testErrorContains(t, err, []string{"not a valid definition type", "function", "event", "error"})
}

func TestUnsupportedPackedAbiType_Error(t *testing.T) {
	err := &UnsupportedPackedAbiType{Type: "tuple"}
	testErrorContains(t, err, []string{"not supported for packed encoding", "tuple"})
}
