package ens

import (
	"regexp"
)

var hexRegex = regexp.MustCompile(`^0x[a-fA-F0-9]+$`)

// EncodedLabelToLabelhash converts an encoded label (e.g., "[abc123...]") to a labelhash.
// Returns empty string if the label is not a valid encoded labelhash.
//
// Encoded labels are in the format: [<64 hex chars>]
//
// Example:
//
//	hash := EncodedLabelToLabelhash("[4f5b812789fc606be1b3b16908db13fc7a9adf7ca72641f84d75b47069d3d7f0]")
//	// "0x4f5b812789fc606be1b3b16908db13fc7a9adf7ca72641f84d75b47069d3d7f0"
//
//	hash := EncodedLabelToLabelhash("eth")
//	// "" (not an encoded label)
func EncodedLabelToLabelhash(label string) string {
	if len(label) != 66 {
		return ""
	}
	if label[0] != '[' {
		return ""
	}
	if label[65] != ']' {
		return ""
	}

	hash := "0x" + label[1:65]
	if !hexRegex.MatchString(hash) {
		return ""
	}

	return hash
}
