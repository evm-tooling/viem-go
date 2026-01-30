package signature

import (
	"errors"
	"math/big"
	"strings"
)

// SerializeCompactSignature converts a CompactSignature into hex format.
//
// Example:
//
//	hex, err := SerializeCompactSignature(&CompactSignature{
//		R: "0x68a020a209d3d56c46f38cc50a33f704f4a9a10a59377f8dd762ac66910e9b90",
//		YParityAndS: "0x7e865ad05c4035ab5792787d4a0297a43617ae897930a6fe4d822b8faea52064",
//	})
//	// "0x68a020a209d3d56c46f38cc50a33f704f4a9a10a59377f8dd762ac66910e9b907e865ad05c4035ab5792787d4a0297a43617ae897930a6fe4d822b8faea52064"
func SerializeCompactSignature(sig *CompactSignature) (string, error) {
	if sig == nil {
		return "", errors.New("compact signature is nil")
	}

	// Parse r and yParityAndS as big integers
	r, err := hexToBigInt(sig.R)
	if err != nil {
		return "", err
	}

	yParityAndS, err := hexToBigInt(sig.YParityAndS)
	if err != nil {
		return "", err
	}

	// Format as 64-byte compact signature
	rHex := padLeft(r.Text(16), 64)
	yParityAndSHex := padLeft(yParityAndS.Text(16), 64)

	return "0x" + rHex + yParityAndSHex, nil
}

// SerializeCompactSignatureBytes converts a CompactSignature into bytes.
func SerializeCompactSignatureBytes(sig *CompactSignature) ([]byte, error) {
	hexStr, err := SerializeCompactSignature(sig)
	if err != nil {
		return nil, err
	}
	return hexToBytes(hexStr), nil
}

// hexToBigInt converts a hex string to a big.Int.
func hexToBigInt(s string) (*big.Int, error) {
	s = strings.TrimPrefix(s, "0x")
	s = strings.TrimPrefix(s, "0X")

	if len(s) == 0 {
		return big.NewInt(0), nil
	}

	n := new(big.Int)
	_, ok := n.SetString(s, 16)
	if !ok {
		return nil, errors.New("invalid hex string")
	}

	return n, nil
}
