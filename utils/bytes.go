package utils

import (
	"encoding/binary"
	"encoding/hex"
	"math/big"
)

// ByteConverter provides a fluent API for byte conversions.
type ByteConverter struct {
	data []byte
	err  error
}

// FromBytes creates a new ByteConverter from a byte slice.
func FromBytes(b []byte) *ByteConverter {
	return &ByteConverter{data: b}
}

// ToHex converts the bytes to a hex string with 0x prefix.
func (c *ByteConverter) ToHex() string {
	return BytesToHex(c.data)
}

// ToInt converts the bytes to an int64.
func (c *ByteConverter) ToInt() int64 {
	return BytesToInt(c.data)
}

// ToUint converts the bytes to a uint64.
func (c *ByteConverter) ToUint() uint64 {
	return BytesToUint(c.data)
}

// ToBigInt converts the bytes to a *big.Int.
func (c *ByteConverter) ToBigInt() *big.Int {
	return BytesToBigInt(c.data)
}

// ToBool converts the bytes to a boolean.
func (c *ByteConverter) ToBool() bool {
	return BytesToBool(c.data)
}

// ToBytes returns the underlying byte slice.
func (c *ByteConverter) ToBytes() []byte {
	return c.data
}

func BytesToInt(b []byte) int64 {
	// Pad to 8 bytes if needed
	if len(b) < 8 {
		padded := make([]byte, 8)
		copy(padded[8-len(b):], b)
		b = padded
	}
	return int64(binary.BigEndian.Uint64(b))
}

func BytesToUint(b []byte) uint64 {
	if len(b) < 8 {
		padded := make([]byte, 8)
		copy(padded[8-len(b):], b)
		b = padded
	}
	return binary.BigEndian.Uint64(b)
}

func BytesToBigInt(b []byte) *big.Int {
	return new(big.Int).SetBytes(b)
}

func BytesToBool(b []byte) bool {
	for _, v := range b {
		if v != 0 {
			return true
		}
	}
	return false
}

func IntToBytes(n int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(n))
	return b
}

func IntToBytesMinimal(n int64) []byte {
	if n == 0 {
		return []byte{0}
	}
	b := IntToBytes(n)
	// Strip leading zeros
	for i, v := range b {
		if v != 0 {
			return b[i:]
		}
	}
	return []byte{0}
}

func UintToBytes(n uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, n)
	return b
}

func BigIntToBytes(n *big.Int) []byte {
	if n == nil {
		return []byte{0}
	}
	return n.Bytes()
}

func BigIntToBytesPadded(n *big.Int, size int) []byte {
	b := n.Bytes()
	if len(b) >= size {
		return b
	}
	padded := make([]byte, size)
	copy(padded[size-len(b):], b)
	return padded
}

func BoolToBytes(v bool) []byte {
	if v {
		return []byte{1}
	}
	return []byte{0}
}

func BytesToHex(b []byte) string {
	return "0x" + hex.EncodeToString(b)
}

func BytesToHexUnprefixed(b []byte) string {
	return hex.EncodeToString(b)
}
