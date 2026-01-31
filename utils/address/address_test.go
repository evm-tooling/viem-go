package address_test

import (
	"testing"

	"golang.org/x/crypto/sha3"

	"github.com/ChefBingbong/viem-go/utils/address"
)

func TestIsAddress(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		strict   bool
		expected bool
	}{
		// Valid addresses
		{"lowercase valid", "0xa5cc3c03994db5b0d9a5eedd10cabab0813678ac", true, true},
		{"checksummed valid", "0xa5cc3c03994DB5b0d9A5eEdD10CabaB0813678AC", true, true},
		{"all lowercase valid strict", "0xa0cf798816d4b9b9866b5330eea46a18382f251e", true, true},

		// Invalid addresses
		{"wrong length", "0xa5cc3c03994db5b0d9a5eedd10cabab0813678a", true, false},
		{"no prefix", "a5cc3c03994db5b0d9a5eedd10cabab0813678ac", true, false},
		{"invalid chars", "0xa5cc3c03994db5b0d9a5eedd10cabab0813678zz", true, false},
		{"empty", "", true, false},
		{"just prefix", "0x", true, false},

		// Checksum validation
		{"invalid checksum strict", "0xa5cc3c03994DB5b0d9A5eEdD10CabaB0813678Ac", true, false},
		{"invalid checksum non-strict", "0xa5cc3c03994DB5b0d9A5eEdD10CabaB0813678Ac", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := address.IsAddress(tt.input, address.IsAddressOptions{Strict: tt.strict})
			if result != tt.expected {
				t.Errorf("IsAddress(%q, strict=%v) = %v, want %v", tt.input, tt.strict, result, tt.expected)
			}
		})
	}
}

func TestChecksumAddress(t *testing.T) {
	// Test cases from viem's getAddress.test.ts
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"viem test 1",
			"0xa0cf798816d4b9b9866b5330eea46a18382f251e",
			"0xA0Cf798816D4b9b9866b5330EEa46a18382f251e",
		},
		{
			"viem test 2",
			"0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266",
			"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		},
		{
			"viem test 3",
			"0x70997970c51812dc3a010c7d01b50e0d17dc79c8",
			"0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
		},
		{
			"viem test 4",
			"0x3c44cdddb6a900fa2b585dd299e03d12fa4293bc",
			"0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC",
		},
		{
			"viem test 5",
			"0x90f79bf6eb2c4f870365e785982e1f101e93b906",
			"0x90F79bf6EB2c4f870365E785982E1f101E93b906",
		},
		{
			"viem test 6",
			"0x15d34aaf54267db7d7c367839aaf71a00a2c6a65",
			"0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65",
		},
		{
			"viem test 7",
			"0xa5cc3c03994db5b0d9a5eEdD10Cabab0813678ac",
			"0xa5cc3c03994DB5b0d9A5eEdD10CabaB0813678AC",
		},
		{
			"all zeros",
			"0x0000000000000000000000000000000000000000",
			"0x0000000000000000000000000000000000000000",
		},
		{
			"vitalik",
			"0xd8da6bf26964af9d7eed9e03e53415d37aa96045",
			"0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := address.ChecksumAddress(tt.input)
			if result != address.Address(tt.expected) {
				t.Errorf("ChecksumAddress(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetAddress(t *testing.T) {
	t.Run("valid lowercase", func(t *testing.T) {
		result, err := address.GetAddress("0xa5cc3c03994db5b0d9a5eEdD10Cabab0813678ac")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		expected := "0xa5cc3c03994DB5b0d9A5eEdD10CabaB0813678AC"
		if result != address.Address(expected) {
			t.Errorf("got %q, want %q", result, expected)
		}
	})

	t.Run("valid checksummed", func(t *testing.T) {
		result, err := address.GetAddress("0xa5cc3c03994DB5b0d9A5eEdD10CabaB0813678AC")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		expected := "0xa5cc3c03994DB5b0d9A5eEdD10CabaB0813678AC"
		if result != address.Address(expected) {
			t.Errorf("got %q, want %q", result, expected)
		}
	})

	t.Run("invalid address", func(t *testing.T) {
		_, err := address.GetAddress("0xinvalid")
		if err == nil {
			t.Error("expected error for invalid address")
		}
	})

	t.Run("invalid length", func(t *testing.T) {
		_, err := address.GetAddress("0xa5cc3c03994db5b0d9a5eedd10cabab0813678a")
		if err == nil {
			t.Error("expected error for invalid length")
		}
	})
}

func TestIsAddressEqual(t *testing.T) {
	t.Run("same address different case", func(t *testing.T) {
		result, err := address.IsAddressEqual(
			"0xa5cc3c03994db5b0d9a5eedd10cabab0813678ac",
			"0xa5cc3c03994DB5b0d9A5eEdD10CabaB0813678AC",
		)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !result {
			t.Error("expected addresses to be equal")
		}
	})

	t.Run("different addresses", func(t *testing.T) {
		result, err := address.IsAddressEqual(
			"0xa5cc3c03994db5b0d9a5eedd10cabab0813678ac",
			"0xd8da6bf26964af9d7eed9e03e53415d37aa96045",
		)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result {
			t.Error("expected addresses to be different")
		}
	})

	t.Run("invalid first address", func(t *testing.T) {
		_, err := address.IsAddressEqual("0xinvalid", "0xa5cc3c03994db5b0d9a5eedd10cabab0813678ac")
		if err == nil {
			t.Error("expected error for invalid address")
		}
	})

	t.Run("invalid second address", func(t *testing.T) {
		_, err := address.IsAddressEqual("0xa5cc3c03994db5b0d9a5eedd10cabab0813678ac", "0xinvalid")
		if err == nil {
			t.Error("expected error for invalid address")
		}
	})
}

func TestGetCreateAddress(t *testing.T) {
	// Test cases from viem's getContractAddress.test.ts
	tests := []struct {
		name     string
		from     string
		nonce    uint64
		expected string
	}{
		{
			"nonce 0",
			"0x1a1e021a302c237453d3d45c7b82b19ceeb7e2e6",
			0,
			"0xFBA3912Ca04dd458c843e2EE08967fC04f3579c2",
		},
		{
			"nonce 5",
			"0xc961145a54C96E3aE9bAA048c4F4D6b04C13916b",
			5,
			"0x30b3F7E5B61d6343Af9B4f98Ed92c003d8fc600F",
		},
		{
			"nonce 69420",
			"0xc961145a54C96E3aE9bAA048c4F4D6b04C13916b",
			69420,
			"0xDf2e056f7062790dF95A472f691670717Ae7b1B6",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := address.GetCreateAddress(address.GetCreateAddressOptions{
				From:  tt.from,
				Nonce: tt.nonce,
			})
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != address.Address(tt.expected) {
				t.Errorf("GetCreateAddress(from=%q, nonce=%d) = %q, want %q", tt.from, tt.nonce, result, tt.expected)
			}
		})
	}
}

func TestGetCreate2Address(t *testing.T) {
	t.Run("with bytecode", func(t *testing.T) {
		// Test case from viem's getContractAddress.test.ts
		bytecode := hexToBytes("0x6394198df16000526103ff60206004601c335afa6040516060f3")
		salt := []byte("hello world")

		result, err := address.GetCreate2Address(address.GetCreate2AddressOptions{
			From:     "0x1a1e021a302c237453d3d45c7b82b19ceeb7e2e6",
			Bytecode: bytecode,
			Salt:     salt,
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		expected := "0x59fbB593ABe27Cb193b6ee5C5DC7bbde312290aB"
		if result != address.Address(expected) {
			t.Errorf("GetCreate2Address() = %q, want %q", result, expected)
		}
	})

	t.Run("with bytecodeHash", func(t *testing.T) {
		// Same test as "with bytecode" but using pre-computed hash
		bytecode := hexToBytes("0x6394198df16000526103ff60206004601c335afa6040516060f3")
		// Pre-compute keccak256 of bytecode
		bytecodeHash := keccak256(bytecode)
		salt := []byte("hello world")

		result, err := address.GetCreate2Address(address.GetCreate2AddressOptions{
			From:         "0x1a1e021a302c237453d3d45c7b82b19ceeb7e2e6",
			BytecodeHash: bytecodeHash,
			Salt:         salt,
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		// Should match the "with bytecode" test result
		expected := "0x59fbB593ABe27Cb193b6ee5C5DC7bbde312290aB"
		if result != address.Address(expected) {
			t.Errorf("GetCreate2Address() = %q, want %q", result, expected)
		}
	})

	t.Run("missing bytecode", func(t *testing.T) {
		_, err := address.GetCreate2Address(address.GetCreate2AddressOptions{
			From: "0x0000000000000000000000000000000000000000",
			Salt: make([]byte, 32),
		})
		if err == nil {
			t.Error("expected error when bytecode and bytecodeHash are missing")
		}
	})
}

func TestGetContractAddress(t *testing.T) {
	t.Run("CREATE opcode", func(t *testing.T) {
		result, err := address.GetContractAddress("CREATE", address.GetCreateAddressOptions{
			From:  "0x1a1e021a302c237453d3d45c7b82b19ceeb7e2e6",
			Nonce: 0,
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		expected := "0xFBA3912Ca04dd458c843e2EE08967fC04f3579c2"
		if result != address.Address(expected) {
			t.Errorf("got %q, want %q", result, expected)
		}
	})

	t.Run("CREATE2 opcode", func(t *testing.T) {
		bytecode := hexToBytes("0x6394198df16000526103ff60206004601c335afa6040516060f3")
		salt := []byte("hello world")

		result, err := address.GetContractAddress("CREATE2", address.GetCreate2AddressOptions{
			From:     "0x1a1e021a302c237453d3d45c7b82b19ceeb7e2e6",
			Bytecode: bytecode,
			Salt:     salt,
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		expected := "0x59fbB593ABe27Cb193b6ee5C5DC7bbde312290aB"
		if result != address.Address(expected) {
			t.Errorf("got %q, want %q", result, expected)
		}
	})

	t.Run("invalid opcode options", func(t *testing.T) {
		_, err := address.GetContractAddress("CREATE", "invalid")
		if err == nil {
			t.Error("expected error for invalid options")
		}
	})
}

// Helper for tests
func keccak256(data []byte) []byte {
	h := sha3.NewLegacyKeccak256()
	h.Write(data)
	return h.Sum(nil)
}

func hexToBytes(s string) []byte {
	if len(s) >= 2 && s[0:2] == "0x" {
		s = s[2:]
	}
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

// Benchmark tests

func BenchmarkIsAddress(b *testing.B) {
	addr := "0xa5cc3c03994DB5b0d9A5eEdD10CabaB0813678AC"
	for i := 0; i < b.N; i++ {
		address.IsAddress(addr, address.IsAddressOptions{Strict: true})
	}
}

func BenchmarkChecksumAddress(b *testing.B) {
	addr := "0xa5cc3c03994db5b0d9a5eedd10cabab0813678ac"
	for i := 0; i < b.N; i++ {
		address.ChecksumAddress(addr)
	}
}

func BenchmarkGetCreateAddress(b *testing.B) {
	for i := 0; i < b.N; i++ {
		address.GetCreateAddress(address.GetCreateAddressOptions{
			From:  "0x1a1e021a302c237453d3d45c7b82b19ceeb7e2e6",
			Nonce: 0,
		})
	}
}
