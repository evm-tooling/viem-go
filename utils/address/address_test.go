package address_test

import (
	"testing"

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
		{"lowercase valid", "0xa0cf798816d4b9b9866b5330eea46a18382f251e", true, true},
		{"checksummed valid", "0xA0Cf798816D4b9b9866b5330EEa46a18382f251e", true, true},
		{"all caps valid", "0xA0CF798816D4B9B9866B5330EEA46A18382F251E", false, true},

		// Invalid addresses
		{"wrong length", "0xa0cf798816d4b9b9866b5330eea46a18382f251", true, false},
		{"no prefix", "a0cf798816d4b9b9866b5330eea46a18382f251e", true, false},
		{"invalid chars", "0xa0cf798816d4b9b9866b5330eea46a18382f251z", true, false},
		{"empty", "", true, false},
		{"just prefix", "0x", true, false},

		// Checksum validation
		{"invalid checksum strict", "0xa0Cf798816D4b9b9866b5330EEa46a18382f251E", true, false},
		{"invalid checksum non-strict", "0xa0Cf798816D4b9b9866b5330EEa46a18382f251E", false, true},
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
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"address 1",
			"0xa0cf798816d4b9b9866b5330eea46a18382f251e",
			"0xA0Cf798816D4b9b9866b5330EEa46a18382f251e",
		},
		{
			"address 2",
			"0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266",
			"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		},
		{
			"address 3",
			"0x70997970c51812dc3a010c7d01b50e0d17dc79c8",
			"0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
		},
		{
			"address 4",
			"0x3c44cdddb6a900fa2b585dd299e03d12fa4293bc",
			"0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC",
		},
		{
			"all zeros",
			"0x0000000000000000000000000000000000000000",
			"0x0000000000000000000000000000000000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := address.ChecksumAddress(tt.input)
			if result != tt.expected {
				t.Errorf("ChecksumAddress(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetAddress(t *testing.T) {
	t.Run("valid lowercase", func(t *testing.T) {
		result, err := address.GetAddress("0xa0cf798816d4b9b9866b5330eea46a18382f251e")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		expected := "0xA0Cf798816D4b9b9866b5330EEa46a18382f251e"
		if result != expected {
			t.Errorf("got %q, want %q", result, expected)
		}
	})

	t.Run("valid checksummed", func(t *testing.T) {
		result, err := address.GetAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		expected := "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
		if result != expected {
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
			"0xa5cC3c03994DB5b0d9A5EEdD10CabaB0813678AC",
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
			"high nonce",
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
			if result != tt.expected {
				t.Errorf("GetCreateAddress(from=%q, nonce=%d) = %q, want %q", tt.from, tt.nonce, result, tt.expected)
			}
		})
	}
}

func TestGetCreate2Address(t *testing.T) {
	// Test using bytecode directly
	t.Run("with bytecode", func(t *testing.T) {
		// Bytecode from viem test: 0x6394198df16000526103ff60206004601c335afa6040516060f3
		bytecode, _ := hexToBytes("0x6394198df16000526103ff60206004601c335afa6040516060f3")
		// Salt from "hello world" as bytes
		salt := []byte("hello world")

		result, err := address.GetCreate2Address(address.GetCreate2AddressOptions{
			From:     "0x1a1e021a302c237453d3d45c7b82b19ceeb7e2e6",
			Salt:     salt,
			Bytecode: bytecode,
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		expected := "0x59fbB593ABe27Cb193b6ee5C5DC7bbde312290aB"
		if result != expected {
			t.Errorf("got %q, want %q", result, expected)
		}
	})
}

func hexToBytes(s string) ([]byte, error) {
	s = s[2:] // strip 0x
	b := make([]byte, len(s)/2)
	for i := 0; i < len(b); i++ {
		b[i] = hexCharToByte(s[i*2])<<4 | hexCharToByte(s[i*2+1])
	}
	return b, nil
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

func TestGetCreate2AddressOld(t *testing.T) {
	tests := []struct {
		name         string
		from         string
		salt         []byte
		bytecodeHash []byte
		expected     string
	}{
		{
			"basic with zero hash",
			"0x0000000000000000000000000000000000000000",
			make([]byte, 32), // all zeros
			make([]byte, 32), // all zeros bytecodeHash
			"",               // We'll verify it doesn't error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := address.GetCreate2Address(address.GetCreate2AddressOptions{
				From:         tt.from,
				Salt:         tt.salt,
				BytecodeHash: tt.bytecodeHash,
			})
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("GetCreate2Address() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestGetCreate2AddressWithBytecode(t *testing.T) {
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
			From:  "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
			Nonce: 0,
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		expected := "0x5FbDB2315678afecb367f032d93F642f64180aa3"
		if result != expected {
			t.Errorf("got %q, want %q", result, expected)
		}
	})

	t.Run("CREATE2 opcode", func(t *testing.T) {
		result, err := address.GetContractAddress("CREATE2", address.GetCreate2AddressOptions{
			From:         "0x0000000000000000000000000000000000000000",
			Salt:         make([]byte, 32),
			BytecodeHash: make([]byte, 32),
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		expected := "0x4D1A2e2bB4F88F0250f26Ffff098B0b30B26BF38"
		if result != expected {
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

// Benchmark tests

func BenchmarkIsAddress(b *testing.B) {
	addr := "0xa5cC3c03994DB5b0d9A5EEdD10CabaB0813678AC"
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
			From:  "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
			Nonce: 0,
		})
	}
}
