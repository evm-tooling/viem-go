package address

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/sha3"
)

var (
	// ErrInvalidAddress is returned when an address is not valid
	ErrInvalidAddress = errors.New("invalid address")
)

// ChecksumAddress converts an address to EIP-55 checksum format.
// Optionally supports EIP-1191 chain-specific checksums (not recommended for general use).
//
// Warning: EIP-1191 checksum addresses are generally not backwards compatible with the
// wider Ethereum ecosystem, meaning it will break when validated against an application/tool
// that relies on EIP-55 checksum encoding (checksum without chainId).
//
// Example:
//
//	checksumAddress("0xa5cc3c03994db5b0d9a5eedd10cabab0813678ac")
//	// "0xa5cc3c03994DB5b0d9A5eEdD10CabaB0813678AC"
func ChecksumAddress(address string, chainId ...int64) string {
	// Get lowercase address without 0x
	addr := strings.ToLower(strings.TrimPrefix(address, "0x"))
	addr = strings.TrimPrefix(addr, "0X")

	// Prepare the string to hash
	var hashInput string
	if len(chainId) > 0 && chainId[0] > 0 {
		// EIP-1191: include chain ID
		hashInput = fmt.Sprintf("%d0x%s", chainId[0], addr)
	} else {
		hashInput = addr
	}

	// Keccak256 hash
	hash := keccak256([]byte(hashInput))

	// Apply checksum
	result := make([]byte, 40)
	for i := 0; i < 40; i++ {
		c := addr[i]
		hashByte := hash[i/2]

		var nibble byte
		if i%2 == 0 {
			nibble = hashByte >> 4
		} else {
			nibble = hashByte & 0x0f
		}

		if nibble >= 8 && c >= 'a' && c <= 'f' {
			result[i] = c - 32 // Convert to uppercase
		} else {
			result[i] = c
		}
	}

	return "0x" + string(result)
}

// GetAddress validates an address and returns it in checksummed format.
// Returns an error if the address is invalid.
//
// Example:
//
//	getAddress("0xa5cc3c03994db5b0d9a5eedd10cabab0813678ac")
//	// "0xa5cc3c03994DB5b0d9A5eEdD10CabaB0813678AC"
func GetAddress(address string, chainId ...int64) (string, error) {
	if !IsAddress(address, IsAddressOptions{Strict: false}) {
		return "", fmt.Errorf("%w: %s", ErrInvalidAddress, address)
	}

	if len(chainId) > 0 {
		return ChecksumAddress(address, chainId[0]), nil
	}
	return ChecksumAddress(address), nil
}

// keccak256 computes the Keccak-256 hash of input data.
func keccak256(data []byte) []byte {
	h := sha3.NewLegacyKeccak256()
	h.Write(data)
	return h.Sum(nil)
}
