package errors

import "testing"

func TestEnsAvatarInvalidMetadataError_Error(t *testing.T) {
	err := &EnsAvatarInvalidMetadataError{Data: "{}"}
	testErrorContains(t, err, []string{"metadata", "image", "image_url"})
}

func TestEnsAvatarInvalidNftUriError_Error(t *testing.T) {
	err := &EnsAvatarInvalidNftUriError{Reason: "invalid format"}
	testErrorContains(t, err, []string{"invalid", "invalid format"})
}

func TestEnsAvatarUriResolutionError_Error(t *testing.T) {
	err := &EnsAvatarUriResolutionError{URI: "https://example.com/avatar"}
	testErrorContains(t, err, []string{"Unable to resolve", "https://example.com/avatar"})
}

func TestEnsAvatarUnsupportedNamespaceError_Error(t *testing.T) {
	err := &EnsAvatarUnsupportedNamespaceError{Namespace: "erc20"}
	testErrorContains(t, err, []string{"not supported", "erc721", "erc1155"})
}

func TestEnsInvalidChainIdError_Error(t *testing.T) {
	err := &EnsInvalidChainIdError{ChainID: 999999}
	testErrorContains(t, err, []string{"Invalid", "chainId", "999999"})
}
