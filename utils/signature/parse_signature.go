package signature

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
)

var (
	// ErrInvalidSignatureLength is returned when the signature has an invalid length.
	ErrInvalidSignatureLength = errors.New("invalid signature length")
	// ErrInvalidYParityOrV is returned when the yParity or v value is invalid.
	ErrInvalidYParityOrV = errors.New("invalid yParityOrV value")
)

// ParseSignature parses a hex formatted signature into a structured Signature.
//
// Example:
//
//	sig, err := ParseSignature("0x6e100a352ec6ad1b70802290e18aeed190704973570f3b8ed42cb9808e2ea6bf4a90a229a244495b41890987806fcbd2d5d23fc0dbe5f5256c2613c039d76db81c")
//	// sig.R = "0x6e100a352ec6ad1b70802290e18aeed190704973570f3b8ed42cb9808e2ea6bf"
//	// sig.S = "0x4a90a229a244495b41890987806fcbd2d5d23fc0dbe5f5256c2613c039d76db8"
//	// sig.YParity = 1
//	// sig.V = 28
func ParseSignature(signatureHex string) (*Signature, error) {
	// Remove 0x prefix and validate
	sigHex := strings.TrimPrefix(signatureHex, "0x")
	sigHex = strings.TrimPrefix(sigHex, "0X")

	if len(sigHex) != 130 {
		return nil, fmt.Errorf("%w: expected 65 bytes (130 hex chars), got %d", ErrInvalidSignatureLength, len(sigHex)/2)
	}

	// Extract r and s (first 64 bytes = 128 hex chars)
	rHex := "0x" + sigHex[0:64]
	sHex := "0x" + sigHex[64:128]

	// Extract yParityOrV (last byte)
	yParityOrVHex := sigHex[128:130]
	yParityOrV, err := hexToInt(yParityOrVHex)
	if err != nil {
		return nil, fmt.Errorf("failed to parse yParityOrV: %w", err)
	}

	// Determine v and yParity
	var v *big.Int
	var yParity int

	switch yParityOrV {
	case 0:
		yParity = 0
	case 1:
		yParity = 1
	case 27:
		v = big.NewInt(27)
		yParity = 0
	case 28:
		v = big.NewInt(28)
		yParity = 1
	default:
		return nil, fmt.Errorf("%w: %d", ErrInvalidYParityOrV, yParityOrV)
	}

	return &Signature{
		R:       rHex,
		S:       sHex,
		V:       v,
		YParity: yParity,
	}, nil
}

// ParseSignatureBytes parses a 65-byte signature into a structured Signature.
func ParseSignatureBytes(sig []byte) (*Signature, error) {
	if len(sig) != 65 {
		return nil, fmt.Errorf("%w: expected 65 bytes, got %d", ErrInvalidSignatureLength, len(sig))
	}
	return ParseSignature(bytesToHex(sig))
}

// hexToInt converts a hex string to int.
func hexToInt(s string) (int, error) {
	s = strings.TrimPrefix(s, "0x")
	s = strings.TrimPrefix(s, "0X")

	n := new(big.Int)
	_, ok := n.SetString(s, 16)
	if !ok {
		return 0, errors.New("invalid hex string")
	}

	return int(n.Int64()), nil
}
