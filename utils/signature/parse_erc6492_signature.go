package signature

import (
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// ParseErc6492Signature parses a hex-formatted ERC-6492 flavored signature.
// If the signature is not in ERC-6492 format, then the underlying (original) signature is returned.
//
// Example:
//
//	parsed, err := ParseErc6492Signature("0x000000000000000000000000cafebabecafebabecafebabecafebabecafebabe000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000004deadbeef000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000041a461f509887bd19e312c0c58467ce8ff8e300d3c1a90b608a760c5b80318eaf15fe57c96f9175d6cd4daad4663763baa7e78836e067d0163e9a2ccf2ff753f5b1b000000000000000000000000000000000000000000000000000000000000006492649264926492649264926492649264926492649264926492649264926492")
//	// parsed.Address = "0xCAFEBABECAFEBABECAFEBABECAFEBABECAFEBABE"
//	// parsed.Data = "0xdeadbeef"
//	// parsed.Signature = "0xa461f509887bd19e312c0c58467ce8ff8e300d3c1a90b608a760c5b80318eaf15fe57c96f9175d6cd4daad4663763baa7e78836e067d0163e9a2ccf2ff753f5b1b"
func ParseErc6492Signature(signature string) (*Erc6492Signature, error) {
	// If not an ERC-6492 signature, return the original
	if !IsErc6492Signature(signature) {
		return &Erc6492Signature{
			Signature: signature,
		}, nil
	}

	// Decode the ABI-encoded data (excluding the magic bytes)
	sigBytes := hexToBytes(signature)

	// The signature format is: abi.encode(address, bytes, bytes) + magicBytes
	// Remove the last 32 bytes (magic bytes)
	encodedData := sigBytes[:len(sigBytes)-32]

	// Define the ABI types
	addressType, _ := abi.NewType("address", "", nil)
	bytesType, _ := abi.NewType("bytes", "", nil)

	args := abi.Arguments{
		{Type: addressType},
		{Type: bytesType},
		{Type: bytesType},
	}

	// Unpack the data
	unpacked, err := args.Unpack(encodedData)
	if err != nil {
		return nil, err
	}

	if len(unpacked) != 3 {
		return nil, errors.New("invalid ERC-6492 signature: expected 3 decoded values")
	}

	address, ok := unpacked[0].(common.Address)
	if !ok {
		return nil, errors.New("invalid ERC-6492 signature: failed to decode address")
	}

	data, ok := unpacked[1].([]byte)
	if !ok {
		return nil, errors.New("invalid ERC-6492 signature: failed to decode data")
	}

	innerSignature, ok := unpacked[2].([]byte)
	if !ok {
		return nil, errors.New("invalid ERC-6492 signature: failed to decode inner signature")
	}

	return &Erc6492Signature{
		Address:   address.Hex(),
		Data:      bytesToHex(data),
		Signature: bytesToHex(innerSignature),
	}, nil
}
