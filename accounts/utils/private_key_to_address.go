package utils

import (
	"encoding/hex"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

// PrivateKeyToAddress converts an ECDSA private key to an Ethereum address.
//
// The function derives the public key from the private key, then computes the address
// from the public key using keccak256 and taking the last 20 bytes.
//
// Example:
//
//	address, err := PrivateKeyToAddress("0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
//	// "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
func PrivateKeyToAddress(privateKey string) (string, error) {
	// Parse the private key
	key := strings.TrimPrefix(privateKey, "0x")
	key = strings.TrimPrefix(key, "0X")

	ecdsaKey, err := crypto.HexToECDSA(key)
	if err != nil {
		return "", ErrInvalidPrivateKey
	}

	// Get the public key (uncompressed, without 04 prefix for the address calculation)
	publicKey := crypto.FromECDSAPub(&ecdsaKey.PublicKey)

	// Convert to hex and pass to PublicKeyToAddress
	publicKeyHex := "0x" + hex.EncodeToString(publicKey)

	return PublicKeyToAddress(publicKeyHex)
}

// MustPrivateKeyToAddress converts a private key to address or panics on error.
func MustPrivateKeyToAddress(privateKey string) string {
	addr, err := PrivateKeyToAddress(privateKey)
	if err != nil {
		panic(err)
	}
	return addr
}

// PrivateKeyToPublicKey derives the public key from a private key.
//
// Returns the public key in uncompressed format (65 bytes with 04 prefix).
//
// Example:
//
//	publicKey, err := PrivateKeyToPublicKey("0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
//	// "0x04bfcab88f42f8d0a..."
func PrivateKeyToPublicKey(privateKey string) (string, error) {
	// Parse the private key
	key := strings.TrimPrefix(privateKey, "0x")
	key = strings.TrimPrefix(key, "0X")

	ecdsaKey, err := crypto.HexToECDSA(key)
	if err != nil {
		return "", ErrInvalidPrivateKey
	}

	// Get the public key (uncompressed format with 04 prefix)
	publicKey := crypto.FromECDSAPub(&ecdsaKey.PublicKey)

	return "0x" + hex.EncodeToString(publicKey), nil
}

// MustPrivateKeyToPublicKey converts a private key to public key or panics on error.
func MustPrivateKeyToPublicKey(privateKey string) string {
	pubKey, err := PrivateKeyToPublicKey(privateKey)
	if err != nil {
		panic(err)
	}
	return pubKey
}
