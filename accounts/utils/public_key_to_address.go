package utils

import (
	"encoding/hex"
	"errors"
	"strings"

	"github.com/ChefBingbong/viem-go/utils/address"
	"github.com/ChefBingbong/viem-go/utils/hash"
)

var (
	// ErrInvalidPublicKey is returned when the public key is invalid.
	ErrInvalidPublicKey = errors.New("invalid public key")
)

// PublicKeyToAddress converts an ECDSA public key to an Ethereum address.
//
// The public key should be in uncompressed format (65 bytes) or without the 04 prefix (64 bytes).
// The address is derived by taking the last 20 bytes of the keccak256 hash of the public key
// (excluding the 04 prefix).
//
// Example:
//
//	address, err := PublicKeyToAddress("0x04bfcab88f42f8d0a...")
//	// "0x1234567890123456789012345678901234567890"
func PublicKeyToAddress(publicKey string) (string, error) {
	// Remove 0x prefix
	pubKey := strings.TrimPrefix(publicKey, "0x")
	pubKey = strings.TrimPrefix(pubKey, "0X")

	// Validate length (should be 128 or 130 hex chars = 64 or 65 bytes)
	if len(pubKey) != 128 && len(pubKey) != 130 {
		return "", ErrInvalidPublicKey
	}

	// If it includes the 04 prefix, remove it
	// The 04 prefix indicates uncompressed public key format
	if len(pubKey) == 130 {
		if pubKey[:2] != "04" {
			return "", ErrInvalidPublicKey
		}
		pubKey = pubKey[2:]
	}

	// Decode the public key bytes
	pubKeyBytes, err := hex.DecodeString(pubKey)
	if err != nil {
		return "", ErrInvalidPublicKey
	}

	// Keccak256 hash of the public key (without 04 prefix)
	hashResult := hash.Keccak256Bytes(pubKeyBytes)

	// Take the last 20 bytes as the address
	addressBytes := hashResult[12:]
	addressHex := "0x" + hex.EncodeToString(addressBytes)

	// Return checksummed address
	return string(address.ChecksumAddress(addressHex)), nil
}

// MustPublicKeyToAddress converts a public key to address or panics on error.
func MustPublicKeyToAddress(publicKey string) string {
	addr, err := PublicKeyToAddress(publicKey)
	if err != nil {
		panic(err)
	}
	return addr
}
