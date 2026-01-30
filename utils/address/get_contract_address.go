package address

import (
	"errors"
	"math/big"
	"strings"
)

// GetCreateAddressOptions configures CREATE address computation.
type GetCreateAddressOptions struct {
	From  string
	Nonce uint64
}

// GetCreate2AddressOptions configures CREATE2 address computation.
type GetCreate2AddressOptions struct {
	From         string
	Salt         []byte // 32 bytes
	Bytecode     []byte // Raw bytecode (will be hashed)
	BytecodeHash []byte // Pre-computed bytecode hash (32 bytes)
}

// GetContractAddress computes a contract address using either CREATE or CREATE2.
//
// Example:
//
//	// CREATE
//	addr, _ := GetContractAddress("CREATE", GetCreateAddressOptions{
//	  From:  "0x1a1e021a302c237453d3d45c7b82b19ceeb7e2e6",
//	  Nonce: 0,
//	})
//	// "0xFBA3912Ca04dd458c843e2EE08967fC04f3579c2"
//
//	// CREATE2
//	addr, _ := GetContractAddress("CREATE2", GetCreate2AddressOptions{
//	  From:     "0x1a1e021a302c237453d3d45c7b82b19ceeb7e2e6",
//	  Bytecode: bytecode,
//	  Salt:     salt,
//	})
func GetContractAddress(opcode string, opts any) (string, error) {
	switch opcode {
	case "CREATE2":
		createOpts, ok := opts.(GetCreate2AddressOptions)
		if !ok {
			return "", errors.New("invalid options for CREATE2")
		}
		return GetCreate2Address(createOpts)
	default:
		createOpts, ok := opts.(GetCreateAddressOptions)
		if !ok {
			return "", errors.New("invalid options for CREATE")
		}
		return GetCreateAddress(createOpts)
	}
}

// GetCreateAddress computes the address of a contract deployed using CREATE opcode.
// address = keccak256(rlp([sender, nonce]))[12:]
//
// Example:
//
//	addr, _ := GetCreateAddress(GetCreateAddressOptions{
//	  From:  "0x1a1e021a302c237453d3d45c7b82b19ceeb7e2e6",
//	  Nonce: 0,
//	})
//	// "0xFBA3912Ca04dd458c843e2EE08967fC04f3579c2"
func GetCreateAddress(opts GetCreateAddressOptions) (string, error) {
	from, err := GetAddress(opts.From)
	if err != nil {
		return "", err
	}

	fromBytes := hexToBytes(from)
	nonceBytes := encodeNonce(opts.Nonce)

	// RLP encode [from, nonce]
	rlpEncoded := rlpEncodeList([][]byte{fromBytes, nonceBytes})

	// Keccak256 and take last 20 bytes
	hash := keccak256(rlpEncoded)
	addressBytes := hash[12:]

	return GetAddress(bytesToHex(addressBytes))
}

// GetCreate2Address computes the address of a contract deployed using CREATE2 opcode.
// address = keccak256(0xff ++ sender ++ salt ++ keccak256(bytecode))[12:]
//
// Example:
//
//	addr, _ := GetCreate2Address(GetCreate2AddressOptions{
//	  From:     "0x1a1e021a302c237453d3d45c7b82b19ceeb7e2e6",
//	  Bytecode: bytecode,
//	  Salt:     salt,
//	})
func GetCreate2Address(opts GetCreate2AddressOptions) (string, error) {
	from, err := GetAddress(opts.From)
	if err != nil {
		return "", err
	}

	fromBytes := hexToBytes(from)

	// Pad salt to 32 bytes
	salt := padLeft(opts.Salt, 32)

	// Get bytecode hash
	var bytecodeHash []byte
	if len(opts.BytecodeHash) > 0 {
		bytecodeHash = opts.BytecodeHash
	} else if len(opts.Bytecode) > 0 {
		bytecodeHash = keccak256(opts.Bytecode)
	} else {
		return "", errors.New("bytecode or bytecodeHash is required")
	}

	// Concatenate: 0xff ++ from ++ salt ++ bytecodeHash
	data := make([]byte, 1+20+32+32)
	data[0] = 0xff
	copy(data[1:21], fromBytes)
	copy(data[21:53], salt)
	copy(data[53:85], bytecodeHash)

	// Keccak256 and take last 20 bytes
	hash := keccak256(data)
	addressBytes := hash[12:]

	return GetAddress(bytesToHex(addressBytes))
}

// Helper functions

func hexToBytes(s string) []byte {
	s = strings.TrimPrefix(s, "0x")
	s = strings.TrimPrefix(s, "0X")
	if len(s)%2 != 0 {
		s = "0" + s
	}
	b := make([]byte, len(s)/2)
	for i := 0; i < len(b); i++ {
		b[i] = hexCharToByte(s[i*2])<<4 | hexCharToByte(s[i*2+1])
	}
	return b
}

func hexCharToByte(c byte) byte {
	switch {
	case c >= '0' && c <= '9':
		return c - '0'
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10
	default:
		return 0
	}
}

func bytesToHex(b []byte) string {
	const hexChars = "0123456789abcdef"
	result := make([]byte, len(b)*2+2)
	result[0] = '0'
	result[1] = 'x'
	for i, v := range b {
		result[2+i*2] = hexChars[v>>4]
		result[2+i*2+1] = hexChars[v&0x0f]
	}
	return string(result)
}

func padLeft(b []byte, size int) []byte {
	if len(b) >= size {
		return b[:size]
	}
	padded := make([]byte, size)
	copy(padded[size-len(b):], b)
	return padded
}

func encodeNonce(nonce uint64) []byte {
	if nonce == 0 {
		return []byte{}
	}
	n := new(big.Int).SetUint64(nonce)
	b := n.Bytes()
	// Remove leading zeros (except keep at least one byte)
	for len(b) > 1 && b[0] == 0 {
		b = b[1:]
	}
	return b
}

// RLP encoding helpers

func rlpEncodeBytes(b []byte) []byte {
	if len(b) == 1 && b[0] < 0x80 {
		return b
	}
	if len(b) <= 55 {
		return append([]byte{byte(0x80 + len(b))}, b...)
	}
	lenBytes := encodeLength(len(b))
	return append(append([]byte{byte(0xb7 + len(lenBytes))}, lenBytes...), b...)
}

func rlpEncodeList(items [][]byte) []byte {
	var encoded []byte
	for _, item := range items {
		encoded = append(encoded, rlpEncodeBytes(item)...)
	}

	if len(encoded) <= 55 {
		return append([]byte{byte(0xc0 + len(encoded))}, encoded...)
	}
	lenBytes := encodeLength(len(encoded))
	return append(append([]byte{byte(0xf7 + len(lenBytes))}, lenBytes...), encoded...)
}

func encodeLength(length int) []byte {
	if length < 256 {
		return []byte{byte(length)}
	}
	if length < 65536 {
		return []byte{byte(length >> 8), byte(length)}
	}
	if length < 16777216 {
		return []byte{byte(length >> 16), byte(length >> 8), byte(length)}
	}
	return []byte{byte(length >> 24), byte(length >> 16), byte(length >> 8), byte(length)}
}
