package ens

import (
	"strings"

	"golang.org/x/net/idna"
)

// Normalize normalizes an ENS name according to ENSIP-15 (UTS-46).
// This handles Unicode normalization, case folding, and validation.
//
// Example:
//
//	normalized, _ := Normalize("Vitalik.ETH")
//	// "vitalik.eth"
//
//	normalized, _ := Normalize("wevm.eth")
//	// "wevm.eth"
//
// @see https://docs.ens.domains/contract-api-reference/name-processing#normalising-names
// @see https://github.com/ensdomains/docs/blob/9edf9443de4333a0ea7ec658a870672d5d180d53/ens-improvement-proposals/ensip-15-normalization-standard.md
func Normalize(name string) (string, error) {
	if name == "" {
		return "", nil
	}

	// Use IDNA profile for ENS normalization
	// This handles UTS-46 processing including:
	// - Unicode normalization (NFC)
	// - Case folding (lowercase)
	// - Punycode encoding/decoding
	profile := idna.New(
		idna.MapForLookup(),
		idna.Transitional(false), // Use non-transitional processing for ENS
	)

	// Process each label
	labels := strings.Split(name, ".")
	normalizedLabels := make([]string, len(labels))

	for i, label := range labels {
		// Check if it's an encoded labelhash - don't normalize these
		if len(label) == 66 && label[0] == '[' && label[65] == ']' {
			normalizedLabels[i] = label
			continue
		}

		normalized, err := profile.ToUnicode(label)
		if err != nil {
			return "", err
		}
		normalizedLabels[i] = normalized
	}

	return strings.Join(normalizedLabels, "."), nil
}

// MustNormalize normalizes an ENS name, panicking on error.
func MustNormalize(name string) string {
	normalized, err := Normalize(name)
	if err != nil {
		panic(err)
	}
	return normalized
}
