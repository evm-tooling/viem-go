package errors

// RpcError represents a JSON-RPC or Ethereum RPC error (EIP-1474).
type RpcError struct {
	Code         int
	ShortMessage string
	Cause        error
}

func (e *RpcError) Error() string {
	msg := e.ShortMessage
	if msg == "" {
		msg = "An unknown RPC error occurred."
	}
	if e.Cause != nil {
		msg += "\nDetails: " + e.Cause.Error()
	}
	return msg
}

func (e *RpcError) Unwrap() error {
	return e.Cause
}

// ProviderRpcError represents an EIP-1193 provider error.
type ProviderRpcError struct {
	RpcError
	Data interface{}
}

// RPC error codes (EIP-1474)
const (
	RpcCodeParse          = -32700
	RpcCodeInvalidRequest = -32600
	RpcCodeMethodNotFound = -32601
	RpcCodeInvalidParams  = -32602
	RpcCodeInternal       = -32603
	RpcCodeInvalidInput   = -32000
	RpcCodeResourceNotFound    = -32001
	RpcCodeResourceUnavailable = -32002
	RpcCodeTransactionRejected = -32003
	RpcCodeMethodNotSupported  = -32004
	RpcCodeLimitExceeded       = -32005
	RpcCodeVersionUnsupported  = -32006
)

// Provider error codes (EIP-1193, EIP-5792)
const (
	RpcCodeUserRejected      = 4001
	RpcCodeUnauthorized      = 4100
	RpcCodeUnsupportedMethod = 4200
	RpcCodeDisconnected      = 4900
	RpcCodeChainDisconnected = 4901
	RpcCodeSwitchChain       = 4902
	RpcCodeUnsupportedCapability = 5700
	RpcCodeUnsupportedChainId   = 5710
	RpcCodeDuplicateId         = 5720
	RpcCodeUnknownBundleId      = 5730
	RpcCodeBundleTooLarge       = 5740
	RpcCodeAtomicRejected       = 5750
	RpcCodeAtomicityNotSupported = 5760
	RpcCodeWalletConnectSession = 7000
)

// NewParseRpcError creates a Parse error (-32700).
func NewParseRpcError(cause error) *RpcError {
	return &RpcError{
		Code:         RpcCodeParse,
		ShortMessage: "Invalid JSON was received by the server. An error occurred on the server while parsing the JSON text.",
		Cause:        cause,
	}
}

// NewInvalidRequestRpcError creates an Invalid request error (-32600).
func NewInvalidRequestRpcError(cause error) *RpcError {
	return &RpcError{
		Code:         RpcCodeInvalidRequest,
		ShortMessage: "JSON is not a valid request object.",
		Cause:        cause,
	}
}

// NewMethodNotFoundRpcError creates a Method not found error (-32601).
func NewMethodNotFoundRpcError(cause error, method string) *RpcError {
	msg := "The method does not exist / is not available."
	if method != "" {
		msg = "The method \"" + method + "\" does not exist / is not available."
	}
	return &RpcError{Code: RpcCodeMethodNotFound, ShortMessage: msg, Cause: cause}
}

// NewInvalidParamsRpcError creates an Invalid params error (-32602).
func NewInvalidParamsRpcError(cause error) *RpcError {
	return &RpcError{
		Code:         RpcCodeInvalidParams,
		ShortMessage: "Invalid parameters were provided to the RPC method.\nDouble check you have provided the correct parameters.",
		Cause:        cause,
	}
}

// NewInternalRpcError creates an Internal error (-32603).
func NewInternalRpcError(cause error) *RpcError {
	return &RpcError{
		Code:         RpcCodeInternal,
		ShortMessage: "An internal error was received.",
		Cause:        cause,
	}
}

// NewInvalidInputRpcError creates an Invalid input error (-32000).
func NewInvalidInputRpcError(cause error) *RpcError {
	return &RpcError{
		Code:         RpcCodeInvalidInput,
		ShortMessage: "Missing or invalid parameters.\nDouble check you have provided the correct parameters.",
		Cause:        cause,
	}
}

// NewResourceNotFoundRpcError creates a Resource not found error (-32001).
func NewResourceNotFoundRpcError(cause error) *RpcError {
	return &RpcError{
		Code:         RpcCodeResourceNotFound,
		ShortMessage: "Requested resource not found.",
		Cause:        cause,
	}
}

// NewResourceUnavailableRpcError creates a Resource unavailable error (-32002).
func NewResourceUnavailableRpcError(cause error) *RpcError {
	return &RpcError{
		Code:         RpcCodeResourceUnavailable,
		ShortMessage: "Requested resource not available.",
		Cause:        cause,
	}
}

// NewTransactionRejectedRpcError creates a Transaction rejected error (-32003).
func NewTransactionRejectedRpcError(cause error) *RpcError {
	return &RpcError{
		Code:         RpcCodeTransactionRejected,
		ShortMessage: "Transaction creation failed.",
		Cause:        cause,
	}
}

// NewMethodNotSupportedRpcError creates a Method not supported error (-32004).
func NewMethodNotSupportedRpcError(cause error, method string) *RpcError {
	msg := "Method is not supported."
	if method != "" {
		msg = "Method \"" + method + "\" is not supported."
	}
	return &RpcError{Code: RpcCodeMethodNotSupported, ShortMessage: msg, Cause: cause}
}

// NewLimitExceededRpcError creates a Limit exceeded error (-32005).
func NewLimitExceededRpcError(cause error) *RpcError {
	return &RpcError{
		Code:         RpcCodeLimitExceeded,
		ShortMessage: "Request exceeds defined limit.",
		Cause:        cause,
	}
}

// NewJsonRpcVersionUnsupportedError creates a JSON-RPC version not supported error (-32006).
func NewJsonRpcVersionUnsupportedError(cause error) *RpcError {
	return &RpcError{
		Code:         RpcCodeVersionUnsupported,
		ShortMessage: "Version of JSON-RPC protocol is not supported.",
		Cause:        cause,
	}
}

// NewUserRejectedRequestError creates a User rejected request error (4001).
func NewUserRejectedRequestError(cause error) *ProviderRpcError {
	return &ProviderRpcError{
		RpcError: RpcError{
			Code:         RpcCodeUserRejected,
			ShortMessage: "User rejected the request.",
			Cause:        cause,
		},
	}
}

// NewUnauthorizedProviderError creates an Unauthorized error (4100).
func NewUnauthorizedProviderError(cause error) *ProviderRpcError {
	return &ProviderRpcError{
		RpcError: RpcError{
			Code:         RpcCodeUnauthorized,
			ShortMessage: "The requested method and/or account has not been authorized by the user.",
			Cause:        cause,
		},
	}
}

// NewUnsupportedProviderMethodError creates an Unsupported method error (4200).
func NewUnsupportedProviderMethodError(cause error, method string) *ProviderRpcError {
	msg := "The Provider does not support the requested method."
	if method != "" {
		msg = "The Provider does not support the requested method \"" + method + "\"."
	}
	return &ProviderRpcError{
		RpcError: RpcError{Code: RpcCodeUnsupportedMethod, ShortMessage: msg, Cause: cause},
	}
}

// NewProviderDisconnectedError creates a Provider disconnected error (4900).
func NewProviderDisconnectedError(cause error) *ProviderRpcError {
	return &ProviderRpcError{
		RpcError: RpcError{
			Code:         RpcCodeDisconnected,
			ShortMessage: "The Provider is disconnected from all chains.",
			Cause:        cause,
		},
	}
}

// NewChainDisconnectedError creates a Chain disconnected error (4901).
func NewChainDisconnectedError(cause error) *ProviderRpcError {
	return &ProviderRpcError{
		RpcError: RpcError{
			Code:         RpcCodeChainDisconnected,
			ShortMessage: "The Provider is not connected to the requested chain.",
			Cause:        cause,
		},
	}
}

// NewSwitchChainError creates a Switch chain error (4902).
func NewSwitchChainError(cause error) *ProviderRpcError {
	return &ProviderRpcError{
		RpcError: RpcError{
			Code:         RpcCodeSwitchChain,
			ShortMessage: "An error occurred when attempting to switch chain.",
			Cause:        cause,
		},
	}
}

// NewUnsupportedNonOptionalCapabilityError creates error 5700.
func NewUnsupportedNonOptionalCapabilityError(cause error) *ProviderRpcError {
	return &ProviderRpcError{
		RpcError: RpcError{
			Code:         RpcCodeUnsupportedCapability,
			ShortMessage: "This Wallet does not support a capability that was not marked as optional.",
			Cause:        cause,
		},
	}
}

// NewUnsupportedChainIdError creates error 5710.
func NewUnsupportedChainIdError(cause error) *ProviderRpcError {
	return &ProviderRpcError{
		RpcError: RpcError{
			Code:         RpcCodeUnsupportedChainId,
			ShortMessage: "This Wallet does not support the requested chain ID.",
			Cause:        cause,
		},
	}
}

// NewDuplicateIdError creates error 5720.
func NewDuplicateIdError(cause error) *ProviderRpcError {
	return &ProviderRpcError{
		RpcError: RpcError{
			Code:         RpcCodeDuplicateId,
			ShortMessage: "There is already a bundle submitted with this ID.",
			Cause:        cause,
		},
	}
}

// NewUnknownBundleIdError creates error 5730.
func NewUnknownBundleIdError(cause error) *ProviderRpcError {
	return &ProviderRpcError{
		RpcError: RpcError{
			Code:         RpcCodeUnknownBundleId,
			ShortMessage: "This bundle id is unknown / has not been submitted",
			Cause:        cause,
		},
	}
}

// NewBundleTooLargeError creates error 5740.
func NewBundleTooLargeError(cause error) *ProviderRpcError {
	return &ProviderRpcError{
		RpcError: RpcError{
			Code:         RpcCodeBundleTooLarge,
			ShortMessage: "The call bundle is too large for the Wallet to process.",
			Cause:        cause,
		},
	}
}

// NewAtomicReadyWalletRejectedUpgradeError creates error 5750.
func NewAtomicReadyWalletRejectedUpgradeError(cause error) *ProviderRpcError {
	return &ProviderRpcError{
		RpcError: RpcError{
			Code:         RpcCodeAtomicRejected,
			ShortMessage: "The Wallet can support atomicity after an upgrade, but the user rejected the upgrade.",
			Cause:        cause,
		},
	}
}

// NewAtomicityNotSupportedError creates error 5760.
func NewAtomicityNotSupportedError(cause error) *ProviderRpcError {
	return &ProviderRpcError{
		RpcError: RpcError{
			Code:         RpcCodeAtomicityNotSupported,
			ShortMessage: "The wallet does not support atomic execution but the request requires it.",
			Cause:        cause,
		},
	}
}

// NewWalletConnectSessionSettlementError creates error 7000.
func NewWalletConnectSessionSettlementError(cause error) *ProviderRpcError {
	return &ProviderRpcError{
		RpcError: RpcError{
			Code:         RpcCodeWalletConnectSession,
			ShortMessage: "WalletConnect session settlement failed.",
			Cause:        cause,
		},
	}
}

// NewUnknownRpcError creates an unknown RPC error.
func NewUnknownRpcError(cause error) *RpcError {
	return &RpcError{
		Code:         -1,
		ShortMessage: "An unknown RPC error occurred.",
		Cause:        cause,
	}
}
