package errors

// UrlRequiredError is returned when no URL was provided to the Transport.
type UrlRequiredError struct{}

func (e *UrlRequiredError) Error() string {
	return "No URL was provided to the Transport. Please provide a valid RPC URL to the Transport.\nDocs: https://viem.sh/docs/clients/intro"
}
