package errors

import "fmt"

// AbiConstructorNotFoundError is returned when a constructor is not found on the ABI.
type AbiConstructorNotFoundError struct {
	DocsPath string
}

func (e *AbiConstructorNotFoundError) Error() string {
	msg := "A constructor was not found on the ABI.\nMake sure you are using the correct ABI and that the constructor exists on it."
	if e.DocsPath != "" {
		msg += fmt.Sprintf("\nDocs: https://viem.sh%s", e.DocsPath)
	}
	return msg
}

// AbiConstructorParamsNotFoundError is returned when constructor args are provided but inputs are not found.
type AbiConstructorParamsNotFoundError struct {
	DocsPath string
}

func (e *AbiConstructorParamsNotFoundError) Error() string {
	msg := "Constructor arguments were provided (`args`), but a constructor parameters (`inputs`) were not found on the ABI.\nMake sure you are using the correct ABI, and that the `inputs` attribute on the constructor exists."
	if e.DocsPath != "" {
		msg += fmt.Sprintf("\nDocs: https://viem.sh%s", e.DocsPath)
	}
	return msg
}

// AbiDecodingDataSizeInvalidError is returned when data size is not in 32-byte increments.
type AbiDecodingDataSizeInvalidError struct {
	Data string
	Size int
}

func (e *AbiDecodingDataSizeInvalidError) Error() string {
	msg := fmt.Sprintf("Data size of %d bytes is invalid.\nSize must be in increments of 32 bytes (size %% 32 === 0).\nData: %s (%d bytes)", e.Size, e.Data, e.Size)
	return msg
}

// AbiDecodingDataSizeTooSmallError is returned when data is too small for the given parameters.
type AbiDecodingDataSizeTooSmallError struct {
	Data   string
	Params string
	Size   int
}

func (e *AbiDecodingDataSizeTooSmallError) Error() string {
	msg := fmt.Sprintf("Data size of %d bytes is too small for given parameters.\nParams: (%s)\nData:   %s (%d bytes)", e.Size, e.Params, e.Data, e.Size)
	return msg
}

// AbiDecodingZeroDataError is returned when trying to decode zero data with ABI parameters.
type AbiDecodingZeroDataError struct{}

func (e *AbiDecodingZeroDataError) Error() string {
	return `Cannot decode zero data ("0x") with ABI parameters.`
}

// AbiEncodingArrayLengthMismatchError is returned when array length does not match expected.
type AbiEncodingArrayLengthMismatchError struct {
	ExpectedLength int
	GivenLength    int
	Type           string
}

func (e *AbiEncodingArrayLengthMismatchError) Error() string {
	return fmt.Sprintf("ABI encoding array length mismatch for type %s.\nExpected length: %d\nGiven length: %d", e.Type, e.ExpectedLength, e.GivenLength)
}

// AbiEncodingBytesSizeMismatchError is returned when bytes size does not match expected.
type AbiEncodingBytesSizeMismatchError struct {
	ExpectedSize int
	Value       string
	ActualSize  int
}

func (e *AbiEncodingBytesSizeMismatchError) Error() string {
	return fmt.Sprintf(
		`Size of bytes "%s" (bytes%d) does not match expected size (bytes%d).`,
		e.Value, e.ActualSize, e.ExpectedSize,
	)
}

// AbiEncodingLengthMismatchError is returned when params/values length mismatch.
type AbiEncodingLengthMismatchError struct {
	ExpectedLength int
	GivenLength    int
}

func (e *AbiEncodingLengthMismatchError) Error() string {
	return fmt.Sprintf("ABI encoding params/values length mismatch.\nExpected length (params): %d\nGiven length (values): %d", e.ExpectedLength, e.GivenLength)
}

// AbiErrorInputsNotFoundError is returned when error args are provided but inputs not found.
type AbiErrorInputsNotFoundError struct {
	ErrorName string
	DocsPath  string
}

func (e *AbiErrorInputsNotFoundError) Error() string {
	msg := fmt.Sprintf(`Arguments (%sargs%s) were provided to "%s", but "%s" on the ABI does not contain any parameters (%sinputs%s).`, "`", "`", e.ErrorName, e.ErrorName, "`", "`")
	msg += "\nCannot encode error result without knowing what the parameter types are."
	msg += "\nMake sure you are using the correct ABI and that the inputs exist on it."
	if e.DocsPath != "" {
		msg += fmt.Sprintf("\nDocs: https://viem.sh%s", e.DocsPath)
	}
	return msg
}

// AbiErrorNotFoundError is returned when an error is not found on the ABI.
type AbiErrorNotFoundError struct {
	ErrorName string
	DocsPath  string
}

func (e *AbiErrorNotFoundError) Error() string {
	msg := "Error "
	if e.ErrorName != "" {
		msg += fmt.Sprintf(`"%s" `, e.ErrorName)
	}
	msg += "not found on ABI.\nMake sure you are using the correct ABI and that the error exists on it."
	if e.DocsPath != "" {
		msg += fmt.Sprintf("\nDocs: https://viem.sh%s", e.DocsPath)
	}
	return msg
}

// AbiErrorSignatureNotFoundError is returned when encoded error signature is not found.
type AbiErrorSignatureNotFoundError struct {
	Signature string
	DocsPath  string
}

func (e *AbiErrorSignatureNotFoundError) Error() string {
	msg := fmt.Sprintf(`Encoded error signature "%s" not found on ABI.`, e.Signature)
	msg += "\nMake sure you are using the correct ABI and that the error exists on it."
	msg += fmt.Sprintf("\nYou can look up the decoded signature here: https://4byte.sourcify.dev/?q=%s.", e.Signature)
	if e.DocsPath != "" {
		msg += fmt.Sprintf("\nDocs: https://viem.sh%s", e.DocsPath)
	}
	return msg
}

// AbiEventSignatureEmptyTopicsError is returned when trying to extract signature from empty topics.
type AbiEventSignatureEmptyTopicsError struct {
	DocsPath string
}

func (e *AbiEventSignatureEmptyTopicsError) Error() string {
	msg := "Cannot extract event signature from empty topics."
	if e.DocsPath != "" {
		msg += fmt.Sprintf("\nDocs: https://viem.sh%s", e.DocsPath)
	}
	return msg
}

// AbiEventSignatureNotFoundError is returned when encoded event signature is not found.
type AbiEventSignatureNotFoundError struct {
	Signature string
	DocsPath  string
}

func (e *AbiEventSignatureNotFoundError) Error() string {
	msg := fmt.Sprintf(`Encoded event signature "%s" not found on ABI.`, e.Signature)
	msg += "\nMake sure you are using the correct ABI and that the event exists on it."
	msg += fmt.Sprintf("\nYou can look up the signature here: https://4byte.sourcify.dev/?q=%s.", e.Signature)
	if e.DocsPath != "" {
		msg += fmt.Sprintf("\nDocs: https://viem.sh%s", e.DocsPath)
	}
	return msg
}

// AbiEventNotFoundError is returned when an event is not found on the ABI.
type AbiEventNotFoundError struct {
	EventName string
	DocsPath  string
}

func (e *AbiEventNotFoundError) Error() string {
	msg := "Event "
	if e.EventName != "" {
		msg += fmt.Sprintf(`"%s" `, e.EventName)
	}
	msg += "not found on ABI.\nMake sure you are using the correct ABI and that the event exists on it."
	if e.DocsPath != "" {
		msg += fmt.Sprintf("\nDocs: https://viem.sh%s", e.DocsPath)
	}
	return msg
}

// AbiFunctionNotFoundError is returned when a function is not found on the ABI.
type AbiFunctionNotFoundError struct {
	FunctionName string
	DocsPath     string
}

func (e *AbiFunctionNotFoundError) Error() string {
	msg := "Function "
	if e.FunctionName != "" {
		msg += fmt.Sprintf(`"%s" `, e.FunctionName)
	}
	msg += "not found on ABI.\nMake sure you are using the correct ABI and that the function exists on it."
	if e.DocsPath != "" {
		msg += fmt.Sprintf("\nDocs: https://viem.sh%s", e.DocsPath)
	}
	return msg
}

// AbiFunctionOutputsNotFoundError is returned when function has no outputs on ABI.
type AbiFunctionOutputsNotFoundError struct {
	FunctionName string
	DocsPath     string
}

func (e *AbiFunctionOutputsNotFoundError) Error() string {
	msg := fmt.Sprintf(`Function "%s" does not contain any %soutputs%s on ABI.`, e.FunctionName, "`", "`")
	msg += "\nCannot decode function result without knowing what the parameter types are."
	msg += "\nMake sure you are using the correct ABI and that the function exists on it."
	if e.DocsPath != "" {
		msg += fmt.Sprintf("\nDocs: https://viem.sh%s", e.DocsPath)
	}
	return msg
}

// AbiFunctionSignatureNotFoundError is returned when encoded function signature is not found.
type AbiFunctionSignatureNotFoundError struct {
	Signature string
	DocsPath  string
}

func (e *AbiFunctionSignatureNotFoundError) Error() string {
	msg := fmt.Sprintf(`Encoded function signature "%s" not found on ABI.`, e.Signature)
	msg += "\nMake sure you are using the correct ABI and that the function exists on it."
	msg += fmt.Sprintf("\nYou can look up the signature here: https://4byte.sourcify.dev/?q=%s.", e.Signature)
	if e.DocsPath != "" {
		msg += fmt.Sprintf("\nDocs: https://viem.sh%s", e.DocsPath)
	}
	return msg
}

// AbiItemAmbiguityError is returned when ambiguous types are found in overloaded ABI items.
type AbiItemAmbiguityError struct {
	TypeX string
	ItemX string
	TypeY string
	ItemY string
}

func (e *AbiItemAmbiguityError) Error() string {
	msg := "Found ambiguous types in overloaded ABI items."
	msg += fmt.Sprintf("\n`%s` in `%s`, and", e.TypeX, e.ItemX)
	msg += fmt.Sprintf("\n`%s` in `%s`", e.TypeY, e.ItemY)
	msg += "\n\nThese types encode differently and cannot be distinguished at runtime."
	msg += "\nRemove one of the ambiguous items in the ABI."
	return msg
}

// BytesSizeMismatchError is returned when bytes size does not match expected.
type BytesSizeMismatchError struct {
	ExpectedSize int
	GivenSize    int
}

func (e *BytesSizeMismatchError) Error() string {
	return fmt.Sprintf("Expected bytes%d, got bytes%d.", e.ExpectedSize, e.GivenSize)
}

// DecodeLogDataMismatch is returned when data size is too small for non-indexed event parameters.
type DecodeLogDataMismatch struct {
	Data   string
	Params string
	Size   int
}

func (e *DecodeLogDataMismatch) Error() string {
	msg := fmt.Sprintf("Data size of %d bytes is too small for non-indexed event parameters.", e.Size)
	msg += fmt.Sprintf("\nParams: (%s)", e.Params)
	msg += fmt.Sprintf("\nData:   %s (%d bytes)", e.Data, e.Size)
	return msg
}

// DecodeLogTopicsMismatch is returned when a topic is expected for indexed parameter but not found.
type DecodeLogTopicsMismatch struct {
	ParamName string
	EventSig  string
}

func (e *DecodeLogTopicsMismatch) Error() string {
	msg := "Expected a topic for indexed event parameter"
	if e.ParamName != "" {
		msg += fmt.Sprintf(` "%s"`, e.ParamName)
	}
	msg += fmt.Sprintf(` on event "%s".`, e.EventSig)
	return msg
}

// InvalidAbiEncodingTypeError is returned when type is not a valid encoding type.
type InvalidAbiEncodingTypeError struct {
	Type    string
	DocsPath string
}

func (e *InvalidAbiEncodingTypeError) Error() string {
	msg := fmt.Sprintf(`Type "%s" is not a valid encoding type.`, e.Type)
	msg += "\nPlease provide a valid ABI type."
	if e.DocsPath != "" {
		msg += fmt.Sprintf("\nDocs: https://viem.sh%s", e.DocsPath)
	}
	return msg
}

// InvalidAbiDecodingTypeError is returned when type is not a valid decoding type.
type InvalidAbiDecodingTypeError struct {
	Type    string
	DocsPath string
}

func (e *InvalidAbiDecodingTypeError) Error() string {
	msg := fmt.Sprintf(`Type "%s" is not a valid decoding type.`, e.Type)
	msg += "\nPlease provide a valid ABI type."
	if e.DocsPath != "" {
		msg += fmt.Sprintf("\nDocs: https://viem.sh%s", e.DocsPath)
	}
	return msg
}

// InvalidArrayError is returned when value is not a valid array.
type InvalidArrayError struct {
	Value interface{}
}

func (e *InvalidArrayError) Error() string {
	return fmt.Sprintf(`Value "%v" is not a valid array.`, e.Value)
}

// InvalidDefinitionTypeError is returned when type is not a valid definition type.
type InvalidDefinitionTypeError struct {
	Type string
}

func (e *InvalidDefinitionTypeError) Error() string {
	msg := fmt.Sprintf(`"%s" is not a valid definition type.`, e.Type)
	msg += "\nValid types: \"function\", \"event\", \"error\""
	return msg
}

// UnsupportedPackedAbiType is returned when type is not supported for packed encoding.
type UnsupportedPackedAbiType struct {
	Type interface{}
}

func (e *UnsupportedPackedAbiType) Error() string {
	return fmt.Sprintf(`Type "%v" is not supported for packed encoding.`, e.Type)
}

