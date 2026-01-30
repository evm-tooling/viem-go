package ens

import (
	"strings"

	"golang.org/x/crypto/sha3"
)

// Namehash computes the namehash of an ENS name according to EIP-137.
//
// The namehash algorithm recursively hashes each label from right to left,
// starting with 32 zero bytes as the initial hash.
//
// Note: Since ENS names prohibit certain forbidden characters (e.g. underscore)
// and have other validation rules, you likely want to normalize ENS names
// with UTS-46 normalization before passing them to namehash.
// You can use the Normalize function for this.
//
// Example:
//
//	hash := Namehash("vitalik.eth")
//	// "0xee6c4522aab0003e8d14cd40a6af439055fd2577951148c14b6cea9a53475835"
//
//	hash := Namehash("eth")
//	// "0x93cdeb708b7545dc668eb9280176169d1c33cfd8ed6f04690a0bcc88a93fc4ae"
//
// @see https://eips.ethereum.org/EIPS/eip-137
func Namehash(name string) string {
	return bytesToHex(NamehashBytes(name))
}

// NamehashBytes computes the namehash and returns raw bytes.
func NamehashBytes(name string) []byte {
	// Start with 32 zero bytes
	result := make([]byte, 32)

	if name == "" {
		return result
	}

	// Split into labels
	labels := strings.Split(name, ".")

	// Iterate in reverse order building up hash
	for i := len(labels) - 1; i >= 0; i-- {
		label := labels[i]

		// Get the labelhash
		var labelHash []byte
		if encoded := EncodedLabelToLabelhash(label); encoded != "" {
			labelHash = hexToBytes(encoded)
		} else {
			h := sha3.NewLegacyKeccak256()
			h.Write([]byte(label))
			labelHash = h.Sum(nil)
		}

		// Concatenate result with labelHash and hash
		h := sha3.NewLegacyKeccak256()
		h.Write(result)
		h.Write(labelHash)
		result = h.Sum(nil)
	}

	return result
}
