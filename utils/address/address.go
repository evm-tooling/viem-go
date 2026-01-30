package address

import (
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strings"

	"golang.org/x/crypto/sha3"
)

var (
	// addressRegex matches a valid Ethereum address format
	addressRegex = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)

	// ErrInvalidAddress is returned when an address is not valid
	ErrInvalidAddress = errors.New("invalid address")
)

// Address represents an Ethereum address as a string.
type Address = string

// IsAddressOptions configures address validation behavior.
type IsAddressOptions struct {
	// Strict enables checksum validation. Default is true.
	Strict bool
}

// IsAddress checks if a string is a valid Ethereum address.
// By default, it validates the checksum if the address contains mixed case.
func IsAddress(address string, opts ...IsAddressOptions) bool {
	strict := true
	if len(opts) > 0 {
		strict = opts[0].Strict
	}

	// Check basic format
	if !addressRegex.MatchString(address) {
		return false
	}

	// If all lowercase, it's valid (no checksum to verify)
	if strings.ToLower(address) == address {
		return true
	}

	// If strict mode and contains uppercase, verify checksum
	if strict {
		checksummed := ChecksumAddress(address)
		return checksummed == address
	}

	return true
}

// ChecksumAddress converts an address to EIP-55 checksum format.
// Optionally supports EIP-1191 chain-specific checksums (not recommended for general use).
func ChecksumAddress(address string, chainId ...int64) string {
	// Get lowercase address without 0x
	addr := strings.ToLower(strings.TrimPrefix(address, "0x"))
	addr = strings.TrimPrefix(addr, "0X")

	// Prepare the string to hash
	var hashInput string
	if len(chainId) > 0 && chainId[0] > 0 {
		// EIP-1191: include chain ID
		hashInput = fmt.Sprintf("%d0x%s", chainId[0], addr)
	} else {
		hashInput = addr
	}

	// Keccak256 hash
	hash := keccak256([]byte(hashInput))

	// Apply checksum
	result := make([]byte, 40)
	for i := 0; i < 40; i++ {
		c := addr[i]
		hashByte := hash[i/2]

		var nibble byte
		if i%2 == 0 {
			nibble = hashByte >> 4
		} else {
			nibble = hashByte & 0x0f
		}

		if nibble >= 8 && c >= 'a' && c <= 'f' {
			result[i] = c - 32 // Convert to uppercase
		} else {
			result[i] = c
		}
	}

	return "0x" + string(result)
}

// GetAddress validates an address and returns it in checksummed format.
// Returns an error if the address is invalid.
func GetAddress(address string, chainId ...int64) (string, error) {
	if !IsAddress(address, IsAddressOptions{Strict: false}) {
		return "", fmt.Errorf("%w: %s", ErrInvalidAddress, address)
	}

	if len(chainId) > 0 {
		return ChecksumAddress(address, chainId[0]), nil
	}
	return ChecksumAddress(address), nil
}

// IsAddressEqual compares two addresses for equality (case-insensitive).
func IsAddressEqual(a, b string) (bool, error) {
	if !IsAddress(a, IsAddressOptions{Strict: false}) {
		return false, fmt.Errorf("%w: %s", ErrInvalidAddress, a)
	}
	if !IsAddress(b, IsAddressOptions{Strict: false}) {
		return false, fmt.Errorf("%w: %s", ErrInvalidAddress, b)
	}
	return strings.EqualFold(a, b), nil
}

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

func keccak256(data []byte) []byte {
	h := sha3.NewLegacyKeccak256()
	h.Write(data)
	return h.Sum(nil)
}

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
