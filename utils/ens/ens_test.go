package ens_test

import (
	"testing"

	"github.com/ChefBingbong/viem-go/utils/ens"
)

func TestLabelhash(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"eth",
			"eth",
			"0x4f5b812789fc606be1b3b16908db13fc7a9adf7ca72641f84d75b47069d3d7f0",
		},
		{
			"awkweb",
			"awkweb",
			"0x7aaad03ddcacc63166440f59c14a1a2c97ee381014b59c58f55b49ab05f31a38",
		},
		{
			"empty",
			"",
			"0x0000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			"encoded label",
			"[9c22ff5f21f0b81b113e63f7db6da94fedef11b2119b4088b89664fb9a3cb658]",
			"0x9c22ff5f21f0b81b113e63f7db6da94fedef11b2119b4088b89664fb9a3cb658",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ens.Labelhash(tt.input)
			if result != tt.expected {
				t.Errorf("Labelhash(%q) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNamehash(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"eth",
			"eth",
			"0x93cdeb708b7545dc668eb9280176169d1c33cfd8ed6f04690a0bcc88a93fc4ae",
		},
		{
			"alice.eth",
			"alice.eth",
			"0x787192fc5378cc32aa956ddfdedbf26b24e8d78e40109add0eea2c1a012c3dec",
		},
		{
			"iam.alice.eth",
			"iam.alice.eth",
			"0x5bec9e288ed3df984a80a1ac48538a7f19370794d676506adfbddefad210775b",
		},
		{
			"awkweb.eth",
			"awkweb.eth",
			"0x52d0f5fbf348925621be297a61b88ec492ebbbdfa9477d82892e2786020ad61c",
		},
		{
			"empty",
			"",
			"0x0000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			"encoded label",
			"[9c22ff5f21f0b81b113e63f7db6da94fedef11b2119b4088b89664fb9a3cb658].eth",
			"0xeb4f647bea6caa36333c816d7b46fdcb05f9466ecacc140ea8c66faf15b3d9f1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ens.Namehash(tt.input)
			if result != tt.expected {
				t.Errorf("Namehash(%q) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestEncodedLabelToLabelhash(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"valid encoded label",
			"[4f5b812789fc606be1b3b16908db13fc7a9adf7ca72641f84d75b47069d3d7f0]",
			"0x4f5b812789fc606be1b3b16908db13fc7a9adf7ca72641f84d75b47069d3d7f0",
		},
		{
			"not encoded - regular label",
			"eth",
			"",
		},
		{
			"wrong length",
			"[abc]",
			"",
		},
		{
			"missing brackets",
			"4f5b812789fc606be1b3b16908db13fc7a9adf7ca72641f84d75b47069d3d7f0",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ens.EncodedLabelToLabelhash(tt.input)
			if result != tt.expected {
				t.Errorf("EncodedLabelToLabelhash(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestEncodeLabelhash(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"with prefix",
			"0x4f5b812789fc606be1b3b16908db13fc7a9adf7ca72641f84d75b47069d3d7f0",
			"[4f5b812789fc606be1b3b16908db13fc7a9adf7ca72641f84d75b47069d3d7f0]",
		},
		{
			"without prefix",
			"4f5b812789fc606be1b3b16908db13fc7a9adf7ca72641f84d75b47069d3d7f0",
			"[4f5b812789fc606be1b3b16908db13fc7a9adf7ca72641f84d75b47069d3d7f0]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ens.EncodeLabelhash(tt.input)
			if result != tt.expected {
				t.Errorf("EncodeLabelhash(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{
			"lowercase",
			"vitalik.eth",
			"vitalik.eth",
			false,
		},
		{
			"uppercase",
			"VITALIK.ETH",
			"vitalik.eth",
			false,
		},
		{
			"mixed case",
			"Vitalik.ETH",
			"vitalik.eth",
			false,
		},
		{
			"empty",
			"",
			"",
			false,
		},
		{
			"with encoded label",
			"[4f5b812789fc606be1b3b16908db13fc7a9adf7ca72641f84d75b47069d3d7f0].eth",
			"[4f5b812789fc606be1b3b16908db13fc7a9adf7ca72641f84d75b47069d3d7f0].eth",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ens.Normalize(tt.input)
			if tt.hasError && err == nil {
				t.Errorf("Normalize(%q) expected error", tt.input)
			}
			if !tt.hasError && err != nil {
				t.Errorf("Normalize(%q) unexpected error: %v", tt.input, err)
			}
			if !tt.hasError && result != tt.expected {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestPacketToBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []byte
	}{
		{
			"simple",
			"eth",
			[]byte{3, 'e', 't', 'h', 0},
		},
		{
			"two labels",
			"vitalik.eth",
			[]byte{7, 'v', 'i', 't', 'a', 'l', 'i', 'k', 3, 'e', 't', 'h', 0},
		},
		{
			"empty",
			"",
			[]byte{0},
		},
		{
			"with leading dot",
			".eth",
			[]byte{3, 'e', 't', 'h', 0},
		},
		{
			"with trailing dot",
			"eth.",
			[]byte{3, 'e', 't', 'h', 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ens.PacketToBytes(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("PacketToBytes(%q) length = %d, want %d", tt.input, len(result), len(tt.expected))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("PacketToBytes(%q)[%d] = %d, want %d", tt.input, i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestToCoinType(t *testing.T) {
	tests := []struct {
		name     string
		chainId  int
		expected uint64
		hasError bool
	}{
		{
			"ethereum mainnet",
			1,
			60,
			false,
		},
		{
			"optimism",
			10,
			2147483658,
			false,
		},
		{
			"polygon",
			137,
			2147483785,
			false,
		},
		{
			"arbitrum",
			42161,
			2147525809,
			false,
		},
		{
			"negative chain id",
			-1,
			0,
			true,
		},
		{
			"too large chain id",
			0x80000000,
			0,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ens.ToCoinType(tt.chainId)
			if tt.hasError && err == nil {
				t.Errorf("ToCoinType(%d) expected error", tt.chainId)
			}
			if !tt.hasError && err != nil {
				t.Errorf("ToCoinType(%d) unexpected error: %v", tt.chainId, err)
			}
			if !tt.hasError && result != tt.expected {
				t.Errorf("ToCoinType(%d) = %d, want %d", tt.chainId, result, tt.expected)
			}
		})
	}
}

func TestNamehashWithEncodedLabel(t *testing.T) {
	// Test that namehash works with encoded labels
	label := "eth"
	labelHash := ens.Labelhash(label)
	encoded := ens.EncodeLabelhash(labelHash)

	// Namehash of encoded label should equal namehash of original
	result1 := ens.Namehash(label)
	result2 := ens.Namehash(encoded)

	if result1 != result2 {
		t.Errorf("Namehash(%q) = %s, Namehash(%q) = %s, expected equal", label, result1, encoded, result2)
	}
}

// Benchmark tests

func BenchmarkLabelhash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ens.Labelhash("vitalik")
	}
}

func BenchmarkNamehash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ens.Namehash("vitalik.eth")
	}
}

func BenchmarkNormalize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ens.Normalize("Vitalik.ETH")
	}
}
