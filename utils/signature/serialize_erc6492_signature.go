package signature

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// SerializeErc6492SignatureParams contains parameters for serializing an ERC-6492 signature.
type SerializeErc6492SignatureParams struct {
	// Address is the ERC-4337 Account Factory address to use for counterfactual verification.
	Address string
	// Data is the calldata to pass to deploy account (if not deployed).
	Data string
	// Signature is the original signature.
	Signature string
}

// SerializeErc6492Signature serializes an ERC-6492 flavoured signature into hex format.
//
// Example:
//
//	hex, err := SerializeErc6492Signature(SerializeErc6492SignatureParams{
//		Address:   "0xCAFEBABECAFEBABECAFEBABECAFEBABECAFEBABE",
//		Data:      "0xdeadbeef",
//		Signature: "0xa461f509...",
//	})
func SerializeErc6492Signature(params SerializeErc6492SignatureParams) (string, error) {
	// Define the ABI types
	addressType, _ := abi.NewType("address", "", nil)
	bytesType, _ := abi.NewType("bytes", "", nil)

	args := abi.Arguments{
		{Type: addressType},
		{Type: bytesType},
		{Type: bytesType},
	}

	// Convert values
	address := common.HexToAddress(params.Address)
	data := hexToBytes(params.Data)
	signature := hexToBytes(params.Signature)

	// Pack the data
	encoded, err := args.Pack(address, data, signature)
	if err != nil {
		return "", err
	}

	// Append magic bytes
	magicBytes := hexToBytes(Erc6492MagicBytes)
	result := append(encoded, magicBytes...)

	return bytesToHex(result), nil
}

// SerializeErc6492SignatureBytes returns the serialized signature as bytes.
func SerializeErc6492SignatureBytes(params SerializeErc6492SignatureParams) ([]byte, error) {
	hexStr, err := SerializeErc6492Signature(params)
	if err != nil {
		return nil, err
	}
	return hexToBytes(hexStr), nil
}
