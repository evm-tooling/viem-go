package signature

import (
	"fmt"
	"math/big"
	"strings"
)

// ParseCompactSignature parses a hex formatted compact signature into a structured CompactSignature.
//
// Example:
//
//	sig, err := ParseCompactSignature("0x68a020a209d3d56c46f38cc50a33f704f4a9a10a59377f8dd762ac66910e9b907e865ad05c4035ab5792787d4a0297a43617ae897930a6fe4d822b8faea52064")
//	// sig.R = "0x68a020a209d3d56c46f38cc50a33f704f4a9a10a59377f8dd762ac66910e9b90"
//	// sig.YParityAndS = "0x7e865ad05c4035ab5792787d4a0297a43617ae897930a6fe4d822b8faea52064"
func ParseCompactSignature(signatureHex string) (*CompactSignature, error) {
	// Remove 0x prefix and validate
	sigHex := strings.TrimPrefix(signatureHex, "0x")
	sigHex = strings.TrimPrefix(sigHex, "0X")

	if len(sigHex) != 128 {
		return nil, fmt.Errorf("%w: expected 64 bytes (128 hex chars), got %d", ErrInvalidSignatureLength, len(sigHex)/2)
	}

	// Extract r and yParityAndS (each 32 bytes = 64 hex chars)
	rHex := "0x" + sigHex[0:64]
	yParityAndSHex := "0x" + sigHex[64:128]

	return &CompactSignature{
		R:           rHex,
		YParityAndS: yParityAndSHex,
	}, nil
}

// ParseCompactSignatureBytes parses a 64-byte compact signature.
func ParseCompactSignatureBytes(sig []byte) (*CompactSignature, error) {
	if len(sig) != 64 {
		return nil, fmt.Errorf("%w: expected 64 bytes, got %d", ErrInvalidSignatureLength, len(sig))
	}
	return ParseCompactSignature(bytesToHex(sig))
}

// numberToHex converts a big.Int to a hex string with optional size padding.
func numberToHex(n *big.Int, size int) string {
	if n == nil {
		return "0x0"
	}
	hexStr := n.Text(16)
	if size > 0 {
		targetLen := size * 2
		if len(hexStr) < targetLen {
			hexStr = strings.Repeat("0", targetLen-len(hexStr)) + hexStr
		}
	}
	return "0x" + hexStr
}
