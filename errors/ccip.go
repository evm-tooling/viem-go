package errors

// OffchainLookupError is returned when offchain lookup fails.
type OffchainLookupError struct {
	Cause            error
	CallbackSelector string
	Data             string
	ExtraData        string
	Sender           string
	URLs             []string
}

func (e *OffchainLookupError) Error() string {
	if e.Cause != nil {
		return e.Cause.Error()
	}
	msg := "An error occurred while fetching for an offchain result."
	if len(e.URLs) > 0 {
		msg += "\nOffchain Gateway Call:"
		msg += "\n  Gateway URL(s):"
		for _, u := range e.URLs {
			msg += "\n    " + u
		}
	}
	if e.Sender != "" {
		msg += "\n  Sender: " + e.Sender
	}
	if e.Data != "" {
		msg += "\n  Data: " + e.Data
	}
	if e.CallbackSelector != "" {
		msg += "\n  Callback selector: " + e.CallbackSelector
	}
	if e.ExtraData != "" {
		msg += "\n  Extra data: " + e.ExtraData
	}
	return msg
}

func (e *OffchainLookupError) Unwrap() error {
	return e.Cause
}

// OffchainLookupResponseMalformedError is returned when offchain response is malformed.
type OffchainLookupResponseMalformedError struct {
	URL    string
	Result string
}

func (e *OffchainLookupResponseMalformedError) Error() string {
	msg := "Offchain gateway response is malformed. Response data must be a hex value."
	if e.URL != "" {
		msg += "\nGateway URL: " + e.URL
	}
	if e.Result != "" {
		msg += "\nResponse: " + e.Result
	}
	return msg
}

// OffchainLookupSenderMismatchError is returned when sender does not match target.
type OffchainLookupSenderMismatchError struct {
	Sender string
	To     string
}

func (e *OffchainLookupSenderMismatchError) Error() string {
	msg := "Reverted sender address does not match target contract address (`to`)."
	if e.To != "" || e.Sender != "" {
		msg += "\nContract address: " + e.To
		msg += "\nOffchainLookup sender address: " + e.Sender
	}
	return msg
}
