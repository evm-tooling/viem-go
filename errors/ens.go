package errors

import "fmt"

// EnsAvatarInvalidMetadataError is returned when ENS avatar metadata is invalid.
type EnsAvatarInvalidMetadataError struct {
	Data string
}

func (e *EnsAvatarInvalidMetadataError) Error() string {
	msg := "Unable to extract image from metadata. The metadata may be malformed or invalid."
	msg += "\n- Metadata must be a JSON object with at least an `image`, `image_url` or `image_data` property."
	if e.Data != "" {
		msg += "\n\nProvided data: " + e.Data
	}
	return msg
}

// EnsAvatarInvalidNftUriError is returned when ENS NFT avatar URI is invalid.
type EnsAvatarInvalidNftUriError struct {
	Reason string
}

func (e *EnsAvatarInvalidNftUriError) Error() string {
	return "ENS NFT avatar URI is invalid. " + e.Reason
}

// EnsAvatarUriResolutionError is returned when ENS avatar URI cannot be resolved.
type EnsAvatarUriResolutionError struct {
	URI string
}

func (e *EnsAvatarUriResolutionError) Error() string {
	return fmt.Sprintf("Unable to resolve ENS avatar URI \"%s\". The URI may be malformed, invalid, or does not respond with a valid image.", e.URI)
}

// EnsAvatarUnsupportedNamespaceError is returned when ENS NFT namespace is not supported.
type EnsAvatarUnsupportedNamespaceError struct {
	Namespace string
}

func (e *EnsAvatarUnsupportedNamespaceError) Error() string {
	return fmt.Sprintf("ENS NFT avatar namespace \"%s\" is not supported. Must be \"erc721\" or \"erc1155\".", e.Namespace)
}

// EnsInvalidChainIdError is returned when ENS chain ID is invalid.
type EnsInvalidChainIdError struct {
	ChainID int64
}

func (e *EnsInvalidChainIdError) Error() string {
	return fmt.Sprintf("Invalid ENSIP-11 chainId: %d. Must be between 0 and 0x7fffffff, or 1.", e.ChainID)
}
