package signature

import (
	"golang.org/x/crypto/sha3"
)

// HashMessage computes the Ethereum Signed Message hash.
// This prepends the standard prefix and then hashes with keccak256.
//
// Example:
//
//	hash := HashMessage(NewSignableMessage("hello world"))
//	// "0xd9eba16ed0ecae432b71fe008c98cc872bb4cc214d3220a36f365326cf807d68"
func HashMessage(message SignableMessage) string {
	prefixed := ToPrefixedMessage(message)
	return keccak256Hex(prefixed)
}

// HashMessageBytes returns the hash as bytes.
func HashMessageBytes(message SignableMessage) []byte {
	prefixed := ToPrefixedMessageBytes(message)
	return keccak256Bytes(prefixed)
}

// keccak256Hex computes keccak256 hash and returns hex string.
func keccak256Hex(data any) string {
	return bytesToHex(keccak256Bytes(data))
}

// keccak256Bytes computes keccak256 hash and returns bytes.
func keccak256Bytes(data any) []byte {
	var b []byte
	switch v := data.(type) {
	case []byte:
		b = v
	case string:
		if isHex(v) {
			b = hexToBytes(v)
		} else {
			b = []byte(v)
		}
	default:
		return nil
	}

	h := sha3.NewLegacyKeccak256()
	h.Write(b)
	return h.Sum(nil)
}
