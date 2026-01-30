package ens

import "strings"

// EncodeLabelhash encodes a hash as an ENS encoded label.
// This is the inverse of EncodedLabelToLabelhash.
//
// Example:
//
//	encoded := EncodeLabelhash("0x4f5b812789fc606be1b3b16908db13fc7a9adf7ca72641f84d75b47069d3d7f0")
//	// "[4f5b812789fc606be1b3b16908db13fc7a9adf7ca72641f84d75b47069d3d7f0]"
func EncodeLabelhash(hash string) string {
	// Remove 0x prefix if present
	hash = strings.TrimPrefix(hash, "0x")
	hash = strings.TrimPrefix(hash, "0X")
	return "[" + hash + "]"
}
