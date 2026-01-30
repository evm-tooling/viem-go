package signature

import (
	"errors"
	"math/big"
	"strings"
)

// SerializeSignature converts a Signature into hex format.
//
// Example:
//
//	hex, err := SerializeSignature(&Signature{
//		R: "0x6e100a352ec6ad1b70802290e18aeed190704973570f3b8ed42cb9808e2ea6bf",
//		S: "0x4a90a229a244495b41890987806fcbd2d5d23fc0dbe5f5256c2613c039d76db8",
//		YParity: 1,
//	})
//	// "0x6e100a352ec6ad1b70802290e18aeed190704973570f3b8ed42cb9808e2ea6bf4a90a229a244495b41890987806fcbd2d5d23fc0dbe5f5256c2613c039d76db81c"
func SerializeSignature(sig *Signature) (string, error) {
	if sig == nil {
		return "", errors.New("signature is nil")
	}

	// Determine yParity from either yParity or v
	yParity, err := getYParity(sig)
	if err != nil {
		return "", err
	}

	// Extract r and s without 0x prefix
	r := strings.TrimPrefix(sig.R, "0x")
	r = strings.TrimPrefix(r, "0X")
	s := strings.TrimPrefix(sig.S, "0x")
	s = strings.TrimPrefix(s, "0X")

	// Pad r and s to 64 characters (32 bytes)
	r = padLeft(r, 64)
	s = padLeft(s, 64)

	// Add recovery byte (27 = 0x1b for yParity=0, 28 = 0x1c for yParity=1)
	var vByte string
	if yParity == 0 {
		vByte = "1b"
	} else {
		vByte = "1c"
	}

	return "0x" + r + s + vByte, nil
}

// SerializeSignatureBytes converts a Signature into bytes.
func SerializeSignatureBytes(sig *Signature) ([]byte, error) {
	hexStr, err := SerializeSignature(sig)
	if err != nil {
		return nil, err
	}
	return hexToBytes(hexStr), nil
}

// getYParity extracts the yParity value from a signature.
func getYParity(sig *Signature) (int, error) {
	// If yParity is explicitly set (0 or 1), use it
	if sig.YParity == 0 || sig.YParity == 1 {
		return sig.YParity, nil
	}

	// Otherwise, derive from v
	if sig.V != nil {
		v := sig.V.Int64()
		if v == 27 {
			return 0, nil
		}
		if v == 28 {
			return 1, nil
		}
		// EIP-155: v = chainId * 2 + 35 + yParity
		// yParity = (v - 35) % 2
		if v >= 35 {
			mod := new(big.Int).Mod(sig.V, big.NewInt(2))
			if mod.Int64() == 0 {
				return 1, nil
			}
			return 0, nil
		}
	}

	return 0, errors.New("invalid v or yParity value")
}

// padLeft pads a string with zeros on the left to reach the target length.
func padLeft(s string, length int) string {
	if len(s) >= length {
		return s
	}
	return strings.Repeat("0", length-len(s)) + s
}
