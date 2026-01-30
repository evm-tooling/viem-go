package ens

import (
	"encoding/hex"

	"golang.org/x/crypto/sha3"
)

// Labelhash computes the keccak256 hash of an ENS label.
//
// Note: Since ENS labels prohibit certain forbidden characters (e.g. underscore)
// and have other validation rules, you likely want to normalize ENS labels
// with UTS-46 normalization before passing them to labelhash.
// You can use the Normalize function for this.
//
// Example:
//
//	hash := Labelhash("eth")
//	// "0x4f5b812789fc606be1b3b16908db13fc7a9adf7ca72641f84d75b47069d3d7f0"
//
//	hash := Labelhash("vitalik")
//	// "0xaf2caa1c2ca1d027f1ac823b529d0a67cd144264b2789fa2ea4d63a67c7103cc"
func Labelhash(label string) string {
	return bytesToHex(LabelhashBytes(label))
}

// LabelhashBytes computes the labelhash and returns raw bytes.
func LabelhashBytes(label string) []byte {
	// Return 32 zero bytes for empty label
	if label == "" {
		return make([]byte, 32)
	}

	// Check if it's an encoded label
	if encoded := EncodedLabelToLabelhash(label); encoded != "" {
		return hexToBytes(encoded)
	}

	// Compute keccak256 of the label
	h := sha3.NewLegacyKeccak256()
	h.Write([]byte(label))
	return h.Sum(nil)
}

// Helper functions

func bytesToHex(b []byte) string {
	return "0x" + hex.EncodeToString(b)
}

func hexToBytes(s string) []byte {
	if len(s) >= 2 && (s[0:2] == "0x" || s[0:2] == "0X") {
		s = s[2:]
	}
	if len(s)%2 != 0 {
		s = "0" + s
	}
	b, _ := hex.DecodeString(s)
	return b
}
