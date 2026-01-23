package utils

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

// HexConverter provides a fluent API for hex conversions.
type HexConverter struct {
	data string
	err  error
}

// FromHex creates a new HexConverter from a hex string.
func FromHex(s string) *HexConverter {
	return &HexConverter{data: s}
}

// ToBytes converts the hex string to a byte slice.
func (c *HexConverter) ToBytes() ([]byte, error) {
	return HexToBytes(c.data)
}

// ToInt converts the hex string to an int64.
func (c *HexConverter) ToInt() (int64, error) {
	return HexToInt(c.data)
}

// ToUint converts the hex string to a uint64.
func (c *HexConverter) ToUint() (uint64, error) {
	return HexToUint(c.data)
}

// ToBigInt converts the hex string to a *big.Int.
func (c *HexConverter) ToBigInt() (*big.Int, error) {
	return HexToBigInt(c.data)
}

// ToBool converts the hex string to a boolean.
func (c *HexConverter) ToBool() (bool, error) {
	return HexToBool(c.data)
}

// String returns the original hex string.
func (c *HexConverter) String() string {
	return c.data
}

func HexToBytes(s string) ([]byte, error) {
	s = strip0x(s)
	// Pad odd-length hex strings
	if len(s)%2 != 0 {
		s = "0" + s
	}
	return hex.DecodeString(s)
}

func HexToInt(s string) (int64, error) {
	s = strip0x(s)
	return strconv.ParseInt(s, 16, 64)
}

func HexToUint(s string) (uint64, error) {
	s = strip0x(s)
	return strconv.ParseUint(s, 16, 64)
}

func HexToBigInt(s string) (*big.Int, error) {
	s = strip0x(s)
	n := new(big.Int)
	_, ok := n.SetString(s, 16)
	if !ok {
		return nil, fmt.Errorf("invalid hex string: %s", s)
	}
	return n, nil
}

func HexToBool(s string) (bool, error) {
	b, err := HexToBytes(s)
	if err != nil {
		return false, err
	}
	return BytesToBool(b), nil
}

func strip0x(s string) string {
	if len(s) >= 2 && (s[0:2] == "0x" || s[0:2] == "0X") {
		return s[2:]
	}
	return s
}

func IsValidHex(s string) bool {
	s = strip0x(s)
	if len(s) == 0 {
		return false
	}
	for _, c := range strings.ToLower(s) {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}

func PadHex(s string, byteLen int) string {
	s = strip0x(s)
	targetLen := byteLen * 2
	if len(s) >= targetLen {
		return "0x" + s
	}
	return "0x" + strings.Repeat("0", targetLen-len(s)) + s
}
