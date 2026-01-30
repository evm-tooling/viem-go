package signature

import (
	"errors"
)

// SignatureToCompactSignature converts a signature to an EIP-2098 compact signature.
// https://eips.ethereum.org/EIPS/eip-2098
//
// Example:
//
//	compact, err := SignatureToCompactSignature(&Signature{
//		R: "0x68a020a209d3d56c46f38cc50a33f704f4a9a10a59377f8dd762ac66910e9b90",
//		S: "0x7e865ad05c4035ab5792787d4a0297a43617ae897930a6fe4d822b8faea52064",
//		YParity: 0,
//	})
//	// compact.R = "0x68a020a209d3d56c46f38cc50a33f704f4a9a10a59377f8dd762ac66910e9b90"
//	// compact.YParityAndS = "0x7e865ad05c4035ab5792787d4a0297a43617ae897930a6fe4d822b8faea52064"
func SignatureToCompactSignature(sig *Signature) (*CompactSignature, error) {
	if sig == nil {
		return nil, errors.New("signature is nil")
	}

	// Get yParity from either yParity field or v
	yParity := sig.YParity
	if sig.V != nil && sig.V.Int64() >= 27 {
		yParity = int((sig.V.Int64() - 27) % 2)
	}

	// Convert s to bytes
	sBytes := hexToBytes(sig.S)
	if len(sBytes) == 0 {
		sBytes = make([]byte, 32)
	}

	// Pad to 32 bytes if needed
	if len(sBytes) < 32 {
		padded := make([]byte, 32)
		copy(padded[32-len(sBytes):], sBytes)
		sBytes = padded
	}

	// Set the top bit of the first byte if yParity is 1
	yParityAndS := make([]byte, 32)
	copy(yParityAndS, sBytes)
	if yParity == 1 {
		yParityAndS[0] |= 0x80
	}

	return &CompactSignature{
		R:           sig.R,
		YParityAndS: bytesToHex(yParityAndS),
	}, nil
}
